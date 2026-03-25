package adapters

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

const tailCursorHeader = "EXECD-COMMANDS-TAIL-CURSOR"

type CommandsClient interface {
	InterruptCommandWithResponse(ctx context.Context, params *execdapi.InterruptCommandParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.InterruptCommandResponse, error)
	GetCommandStatusWithResponse(ctx context.Context, id string, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetCommandStatusResponse, error)
	GetBackgroundCommandLogsWithResponse(ctx context.Context, id string, params *execdapi.GetBackgroundCommandLogsParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetBackgroundCommandLogsResponse, error)
}

type CommandsAdapter struct {
	client CommandsClient
}

func NewCommandsAdapter(client CommandsClient) *CommandsAdapter {
	return &CommandsAdapter{client: client}
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
