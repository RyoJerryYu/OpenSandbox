package adapters

import (
	"context"
	"net/http"
	"testing"
	"time"

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
