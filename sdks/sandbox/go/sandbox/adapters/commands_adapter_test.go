package adapters

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
)

func TestCommandsAdapterInterruptTreats2xxAsSuccess(t *testing.T) {
	adapter := NewCommandsAdapter(&fakeCommandsClient{
		interruptResponse: &execdapi.InterruptCommandResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
	})

	if err := adapter.Interrupt(context.Background(), "session-1"); err != nil {
		t.Fatalf("interrupt: %v", err)
	}
}

func TestCommandsAdapterGetCommandStatusConvertsFields(t *testing.T) {
	startedAt := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)
	finishedAt := startedAt.Add(2 * time.Second)
	adapter := NewCommandsAdapter(&fakeCommandsClient{
		statusResponse: &execdapi.GetCommandStatusResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &execdapi.CommandStatusResponse{
				Id:         stringPtr("cmd-1"),
				Content:    stringPtr("echo hi"),
				Running:    boolPtr(false),
				ExitCode:   int32Ptr(0),
				Error:      stringPtr(""),
				StartedAt:  &startedAt,
				FinishedAt: &finishedAt,
			},
		},
	})

	status, err := adapter.GetCommandStatus(context.Background(), "cmd-1")
	if err != nil {
		t.Fatalf("get command status: %v", err)
	}
	if status.ID != "cmd-1" {
		t.Fatalf("unexpected status id: %s", status.ID)
	}
	if status.Command != "echo hi" {
		t.Fatalf("unexpected command: %s", status.Command)
	}
	if status.Running {
		t.Fatal("expected finished command")
	}
	if status.ExitCode == nil || *status.ExitCode != 0 {
		t.Fatalf("unexpected exit code: %+v", status.ExitCode)
	}
	if status.StartedAt == nil || !status.StartedAt.Equal(startedAt) {
		t.Fatalf("unexpected started at: %+v", status.StartedAt)
	}
	if status.FinishedAt == nil || !status.FinishedAt.Equal(finishedAt) {
		t.Fatalf("unexpected finished at: %+v", status.FinishedAt)
	}
}

func TestCommandsAdapterBackgroundLogCursorParsing(t *testing.T) {
	adapter := NewCommandsAdapter(&fakeCommandsClient{
		logsResponse: &execdapi.GetBackgroundCommandLogsResponse{
			HTTPResponse: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"EXECD-COMMANDS-TAIL-CURSOR": []string{"17"},
				},
			},
			Body: []byte("stdout line\nstderr line\n"),
		},
	})

	logs, err := adapter.GetBackgroundCommandLogs(context.Background(), "cmd-1", int64Ptr(5))
	if err != nil {
		t.Fatalf("get background logs: %v", err)
	}
	if logs.Content != "stdout line\nstderr line\n" {
		t.Fatalf("unexpected logs content: %q", logs.Content)
	}
	if logs.Cursor == nil || *logs.Cursor != 17 {
		t.Fatalf("unexpected cursor: %+v", logs.Cursor)
	}
}

func TestCommandsAdapterRunStreamExposesOrderedEvents(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripCaptureFunc(func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Accept") != "text/event-stream" {
				t.Fatalf("unexpected accept header: %q", req.Header.Get("Accept"))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"text/event-stream"},
				},
				Body: io.NopCloser(bytes.NewBufferString(
					"data: {\"type\":\"init\",\"text\":\"cmd-1\",\"timestamp\":1}\n\n" +
						"data: {\"type\":\"stdout\",\"text\":\"hello\",\"timestamp\":2}\n\n",
				)),
			}, nil
		}),
	}
	adapter := NewStreamingCommandsAdapter(nil, "http://sandbox.example:44772", &config.ConnectionConfig{
		SSEHTTPClient: httpClient,
	})

	stream, err := adapter.RunStream(context.Background(), "echo hello", nil)
	if err != nil {
		t.Fatalf("run stream: %v", err)
	}

	var eventTypes []string
	for event := range stream.Events {
		eventTypes = append(eventTypes, event.Type)
	}
	if err := <-stream.Done; err != nil {
		t.Fatalf("stream done: %v", err)
	}
	if len(eventTypes) != 2 || eventTypes[0] != "init" || eventTypes[1] != "stdout" {
		t.Fatalf("unexpected event order: %+v", eventTypes)
	}
}

func TestCommandsAdapterRunAggregatesStdoutResultAndComplete(t *testing.T) {
	httpClient := &http.Client{
		Transport: roundTripCaptureFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"text/event-stream"},
				},
				Body: io.NopCloser(bytes.NewBufferString(
					"data: {\"type\":\"init\",\"text\":\"cmd-1\",\"timestamp\":1}\n\n" +
						"data: {\"type\":\"stdout\",\"text\":\"hello\",\"timestamp\":2}\n\n" +
						"data: {\"type\":\"result\",\"results\":{\"text/plain\":\"done\"},\"timestamp\":3}\n\n" +
						"data: {\"type\":\"execution_complete\",\"timestamp\":4,\"execution_time\":5}\n\n",
				)),
			}, nil
		}),
	}
	adapter := NewStreamingCommandsAdapter(nil, "http://sandbox.example:44772", &config.ConnectionConfig{
		SSEHTTPClient: httpClient,
	})

	execution, err := adapter.Run(context.Background(), "echo hello", nil)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if execution.ID != "cmd-1" {
		t.Fatalf("unexpected execution id: %s", execution.ID)
	}
	if len(execution.Logs.Stdout) != 1 || execution.Logs.Stdout[0].Text != "hello" {
		t.Fatalf("unexpected stdout logs: %+v", execution.Logs.Stdout)
	}
	if len(execution.Result) != 1 || execution.Result[0].Text != "done" {
		t.Fatalf("unexpected results: %+v", execution.Result)
	}
	if execution.Complete == nil || execution.Complete.ExecutionTimeMs != 5 {
		t.Fatalf("unexpected completion: %+v", execution.Complete)
	}
}

type fakeCommandsClient struct {
	interruptResponse *execdapi.InterruptCommandResponse
	statusResponse    *execdapi.GetCommandStatusResponse
	logsResponse      *execdapi.GetBackgroundCommandLogsResponse
}

func (f *fakeCommandsClient) InterruptCommandWithResponse(ctx context.Context, params *execdapi.InterruptCommandParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.InterruptCommandResponse, error) {
	return f.interruptResponse, nil
}

func (f *fakeCommandsClient) GetCommandStatusWithResponse(ctx context.Context, id string, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetCommandStatusResponse, error) {
	return f.statusResponse, nil
}

func (f *fakeCommandsClient) GetBackgroundCommandLogsWithResponse(ctx context.Context, id string, params *execdapi.GetBackgroundCommandLogsParams, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetBackgroundCommandLogsResponse, error) {
	return f.logsResponse, nil
}

func int32Ptr(v int32) *int32 {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}
