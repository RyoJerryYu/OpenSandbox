package sandbox

import (
	"context"
	"time"

	sandboxerrors "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/errors"
)

// WaitUntilReadyOptions customizes lifecycle-state and health polling.
type WaitUntilReadyOptions struct {
	Timeout           time.Duration
	PollInterval      time.Duration
	SkipStatePoll     bool
	CustomHealthCheck func(ctx context.Context, sandbox *Sandbox) (bool, error)
}

// WaitUntilReady blocks until the sandbox reports Running and passes the configured health check.
func (s *Sandbox) WaitUntilReady(ctx context.Context, opts *WaitUntilReadyOptions) error {
	if s == nil {
		return nil
	}

	timeout := 60 * time.Second
	pollInterval := 500 * time.Millisecond
	skipStatePoll := false
	var customHealthCheck func(context.Context, *Sandbox) (bool, error)

	if opts != nil {
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		if opts.PollInterval > 0 {
			pollInterval = opts.PollInterval
		}
		skipStatePoll = opts.SkipStatePoll
		customHealthCheck = opts.CustomHealthCheck
	}

	deadline := time.Now().Add(timeout)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if time.Now().After(deadline) {
			return &sandboxerrors.ReadyTimeoutError{Timeout: timeout}
		}

		if !skipStatePoll {
			info, err := s.GetInfo(ctx)
			if err != nil {
				goto wait
			}
			if info == nil || info.Status.State != "Running" {
				goto wait
			}
		}

		if customHealthCheck != nil {
			ok, err := customHealthCheck(ctx, s)
			if err == nil && ok {
				return nil
			}
		} else {
			ok := s.IsHealthy(ctx)
			if ok {
				return nil
			}
		}

	wait:
		timer := time.NewTimer(pollInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
