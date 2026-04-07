package sse

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	sandboxerrors "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/errors"
)

func NewStream(ctx context.Context, resp *http.Response) *Stream {
	events := make(chan Event)
	done := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(done)
		if resp == nil {
			done <- fmt.Errorf("nil response")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			done <- parseStreamError(resp)
			return
		}

		done <- parseEvents(ctx, resp.Body, events)
	}()

	return &Stream{
		Events: events,
		Done:   done,
	}
}

func parseEvents(ctx context.Context, body io.Reader, events chan<- Event) error {
	scanner := bufio.NewScanner(body)
	buffer := make([]byte, 0, 64*1024)
	scanner.Buffer(buffer, 1024*1024)

	var dataLines []string
	var eventName string

	flush := func() error {
		if len(dataLines) == 0 {
			eventName = ""
			return nil
		}

		payload := strings.TrimSpace(strings.Join(dataLines, "\n"))
		dataLines = nil
		ev := Event{
			Event: eventName,
			Data:  []byte(payload),
		}
		eventName = ""

		if payload == "" || !json.Valid(ev.Data) {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case events <- ev:
			return nil
		}
	}

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := strings.TrimRight(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" {
			if err := flush(); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		if strings.HasPrefix(line, "event:") {
			eventName = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}
		if strings.HasPrefix(line, "id:") || strings.HasPrefix(line, "retry:") {
			continue
		}
		if strings.HasPrefix(line, "data:") {
			dataLines = append(dataLines, strings.TrimSpace(strings.TrimPrefix(line, "data:")))
			continue
		}
		dataLines = append(dataLines, strings.TrimSpace(line))
	}

	if err := scanner.Err(); err != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}

	if err := flush(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func parseStreamError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	errResp := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{}
	if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
		return &sandboxerrors.SandboxError{
			Code:       errResp.Code,
			Message:    errResp.Message,
			RequestID:  resp.Header.Get("X-Request-Id"),
			StatusCode: resp.StatusCode,
		}
	}

	message := strings.TrimSpace(string(bytes.TrimSpace(body)))
	if message == "" {
		message = fmt.Sprintf("stream request failed with status %d", resp.StatusCode)
	}
	return &sandboxerrors.SandboxError{
		Message:    message,
		RequestID:  resp.Header.Get("X-Request-Id"),
		StatusCode: resp.StatusCode,
	}
}
