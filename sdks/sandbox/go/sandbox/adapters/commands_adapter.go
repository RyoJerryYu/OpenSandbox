package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/sse"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

const tailCursorHeader = "EXECD-COMMANDS-TAIL-CURSOR"

type CommandsClient interface {
	InterruptCommandWithResponse(ctx context.Context, params *execdapi.InterruptCommandParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.InterruptCommandResponse, error)
	GetCommandStatusWithResponse(ctx context.Context, id string, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetCommandStatusResponse, error)
	GetBackgroundCommandLogsWithResponse(ctx context.Context, id string, params *execdapi.GetBackgroundCommandLogsParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetBackgroundCommandLogsResponse, error)
}

type CommandsAdapter struct {
	client           CommandsClient
	baseURL          string
	connectionConfig *config.ConnectionConfig
}

func NewCommandsAdapter(client CommandsClient) *CommandsAdapter {
	return &CommandsAdapter{client: client}
}

func NewStreamingCommandsAdapter(client CommandsClient, baseURL string, connectionConfig *config.ConnectionConfig) *CommandsAdapter {
	return &CommandsAdapter{
		client:           client,
		baseURL:          baseURL,
		connectionConfig: connectionConfig,
	}
}

func (a *CommandsAdapter) Run(ctx context.Context, command string, opts *models.RunCommandOptions) (*models.CommandExecution, error) {
	stream, err := a.RunStream(ctx, command, opts)
	if err != nil {
		return nil, err
	}

	execution := &models.CommandExecution{}
	for event := range stream.Events {
		dispatchCommandEvent(execution, event)
	}
	if err := <-stream.Done; err != nil {
		return nil, err
	}
	return execution, nil
}

func (a *CommandsAdapter) RunStream(ctx context.Context, command string, opts *models.RunCommandOptions) (*models.CommandStream, error) {
	if strings.TrimSpace(command) == "" {
		return nil, fmt.Errorf("command cannot be empty")
	}
	if a.baseURL == "" {
		return nil, fmt.Errorf("missing execd streaming configuration")
	}

	body, err := json.Marshal(toRunCommandRequest(command, opts))
	if err != nil {
		return nil, err
	}
	req, err := a.newStreamingRequest(ctx, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := a.sseHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	rawStream := sse.NewStream(ctx, resp)
	events := make(chan models.ServerStreamEvent)
	done := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(done)

		for event := range rawStream.Events {
			var parsed models.ServerStreamEvent
			if err := json.Unmarshal(event.Data, &parsed); err != nil {
				continue
			}

			select {
			case <-ctx.Done():
				done <- ctx.Err()
				return
			case events <- parsed:
			}
		}

		done <- <-rawStream.Done
	}()

	return &models.CommandStream{
		Events: events,
		Done:   done,
	}, nil
}

func (a *CommandsAdapter) Interrupt(ctx context.Context, sessionID string) error {
	resp, err := a.client.InterruptCommandWithResponse(ctx, &execdapi.InterruptCommandParams{Id: sessionID})
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("interrupt command failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (a *CommandsAdapter) GetCommandStatus(ctx context.Context, commandID string) (*models.CommandStatus, error) {
	resp, err := a.client.GetCommandStatusWithResponse(ctx, commandID)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("get command status failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromExecdCommandStatus(resp.JSON200), nil
}

func (a *CommandsAdapter) GetBackgroundCommandLogs(ctx context.Context, commandID string, cursor *int64) (*models.CommandLogs, error) {
	params := &execdapi.GetBackgroundCommandLogsParams{}
	if cursor != nil {
		params.Cursor = cursor
	}
	resp, err := a.client.GetBackgroundCommandLogsWithResponse(ctx, commandID, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, NormalizeHTTPError(fmt.Errorf("get background command logs failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}

	logs := &models.CommandLogs{Content: string(resp.Body)}
	if resp.HTTPResponse != nil {
		if raw := getHeaderInsensitive(resp.HTTPResponse.Header, tailCursorHeader); raw != "" {
			if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil {
				logs.Cursor = &parsed
			}
		}
	}
	return logs, nil
}

func getHeaderInsensitive(headers http.Header, target string) string {
	for key, values := range headers {
		if strings.EqualFold(key, target) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func (a *CommandsAdapter) newStreamingRequest(ctx context.Context, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/command", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	applyFilesystemConnectionHeaders(req, a.connectionConfig)
	return req, nil
}

func (a *CommandsAdapter) sseHTTPClient() *http.Client {
	if a.connectionConfig != nil {
		if a.connectionConfig.SSEHTTPClient != nil {
			return a.connectionConfig.SSEHTTPClient
		}
		if a.connectionConfig.HTTPClient != nil {
			return a.connectionConfig.HTTPClient
		}
	}
	return http.DefaultClient
}

func toRunCommandRequest(command string, opts *models.RunCommandOptions) execdapi.RunCommandJSONRequestBody {
	req := execdapi.RunCommandJSONRequestBody{
		Command: command,
	}
	if opts == nil {
		return req
	}
	if opts.Background {
		req.Background = &opts.Background
	}
	if opts.WorkingDirectory != "" {
		req.Cwd = &opts.WorkingDirectory
	}
	if opts.Timeout != nil {
		timeoutMillis := opts.Timeout.Milliseconds()
		req.Timeout = &timeoutMillis
	}
	if opts.UID != nil {
		req.Uid = opts.UID
	}
	if opts.GID != nil {
		req.Gid = opts.GID
	}
	if len(opts.Envs) > 0 {
		envs := make(map[string]string, len(opts.Envs))
		for k, v := range opts.Envs {
			envs[k] = v
		}
		req.Envs = &envs
	}
	return req
}

func dispatchCommandEvent(execution *models.CommandExecution, event models.ServerStreamEvent) {
	switch event.Type {
	case "init":
		execution.ID = event.Text
	case "stdout":
		execution.Logs.Stdout = append(execution.Logs.Stdout, models.OutputMessage{
			Text:      event.Text,
			Timestamp: event.Timestamp,
			IsError:   false,
		})
	case "stderr":
		execution.Logs.Stderr = append(execution.Logs.Stderr, models.OutputMessage{
			Text:      event.Text,
			Timestamp: event.Timestamp,
			IsError:   true,
		})
	case "result":
		text := ""
		if v, ok := event.Results["text/plain"]; ok {
			text, _ = v.(string)
		}
		if text == "" {
			if v, ok := event.Results["text"]; ok {
				text, _ = v.(string)
			}
		}
		execution.Result = append(execution.Result, models.ExecutionResult{
			Text:      text,
			Timestamp: event.Timestamp,
			Raw:       event.Results,
		})
	case "error":
		if event.Error != nil {
			execution.Error = &models.ExecutionError{
				Name:      event.Error.Name,
				Value:     event.Error.Value,
				Timestamp: event.Timestamp,
				Traceback: append([]string(nil), event.Error.Traceback...),
			}
		}
	case "execution_count":
		execution.ExecutionCount = event.ExecutionCount
	case "execution_complete":
		execution.Complete = &models.ExecutionComplete{
			Timestamp:       event.Timestamp,
			ExecutionTimeMs: derefInt64(event.ExecutionTime),
		}
	}
}

func derefInt64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}
