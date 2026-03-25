package sse

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	sandboxerrors "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/errors"
)

func TestParserDecodesJSONEvents(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(bytes.NewBufferString(
			": keep-alive\n" +
				"event: message\n" +
				"data: {\"type\":\"init\",\"text\":\"cmd-1\",\"timestamp\":1}\n\n" +
				"{\"type\":\"stdout\",\"text\":\"hello\",\"timestamp\":2}\n" +
				"\n" +
				"data: {\"type\":\"execution_complete\",\"timestamp\":3}\n\n",
		)),
	}

	stream := NewStream(context.Background(), resp)

	var events []Event
	for event := range stream.Events {
		events = append(events, event)
	}

	if err := <-stream.Done; err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if string(events[0].Data) != "{\"type\":\"init\",\"text\":\"cmd-1\",\"timestamp\":1}" {
		t.Fatalf("unexpected first event: %s", string(events[0].Data))
	}
	if string(events[1].Data) != "{\"type\":\"stdout\",\"text\":\"hello\",\"timestamp\":2}" {
		t.Fatalf("unexpected second event: %s", string(events[1].Data))
	}
}

func TestParserReturnsErrorOnNonOKResponse(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadGateway,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
			"X-Request-Id": []string{"req-123"},
		},
		Body: io.NopCloser(bytes.NewBufferString(`{"code":"bad_gateway","message":"upstream failed"}`)),
	}

	stream := NewStream(context.Background(), resp)

	if event, ok := <-stream.Events; ok {
		t.Fatalf("expected no events, got %+v", event)
	}

	err := <-stream.Done
	var sandboxErr *sandboxerrors.SandboxError
	if !errors.As(err, &sandboxErr) {
		t.Fatalf("expected sandbox error, got %v", err)
	}
	if sandboxErr.Code != "bad_gateway" {
		t.Fatalf("unexpected code: %s", sandboxErr.Code)
	}
	if sandboxErr.RequestID != "req-123" {
		t.Fatalf("unexpected request id: %s", sandboxErr.RequestID)
	}
}

func TestParserStopsOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	reader := &blockingReader{
		first: []byte("data: {\"type\":\"stdout\",\"text\":\"hello\",\"timestamp\":1}\n\n"),
		wait:  make(chan struct{}),
	}
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"text/event-stream"},
		},
		Body: io.NopCloser(reader),
	}

	stream := NewStream(ctx, resp)

	select {
	case <-stream.Events:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for first event")
	}

	cancel()
	close(reader.wait)

	err := <-stream.Done
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

type blockingReader struct {
	first []byte
	sent  bool
	wait  chan struct{}
}

func (r *blockingReader) Read(p []byte) (int, error) {
	if !r.sent {
		r.sent = true
		n := copy(p, r.first)
		return n, nil
	}
	<-r.wait
	return 0, io.EOF
}
