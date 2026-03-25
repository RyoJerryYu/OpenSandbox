package sandbox

import (
	"context"
	"errors"
	"testing"
	"time"

	sandboxerrors "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/errors"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestWaitUntilReadyPollsLifecycleThenHealth(t *testing.T) {
	lifecycle := &sandboxLifecycleFake{
		getSandboxResponses: []*models.SandboxInfo{
			{Status: models.SandboxStatus{State: "Pending"}},
			{Status: models.SandboxStatus{State: "Running"}},
		},
	}
	health := &sandboxHealthFake{results: []bool{false, true}}
	sbx := &Sandbox{
		ID:               "sbx-ready",
		connectionConfig: configWithDefaults(),
		sandboxes:        lifecycle,
		Health:           health,
	}

	err := sbx.WaitUntilReady(context.Background(), &WaitUntilReadyOptions{
		Timeout:           2 * time.Second,
		PollInterval:      5 * time.Millisecond,
		SkipStatePoll:     false,
		CustomHealthCheck: nil,
	})
	if err != nil {
		t.Fatalf("wait until ready: %v", err)
	}
	if lifecycle.getSandboxCallCount < 2 {
		t.Fatalf("expected lifecycle polling, got %d calls", lifecycle.getSandboxCallCount)
	}
	if health.pingCallCount < 2 {
		t.Fatalf("expected health polling, got %d calls", health.pingCallCount)
	}
}

func TestWaitUntilReadyUsesCustomHealthCheckWhenProvided(t *testing.T) {
	lifecycle := &sandboxLifecycleFake{
		getSandboxResponses: []*models.SandboxInfo{
			{Status: models.SandboxStatus{State: "Running"}},
		},
	}
	health := &sandboxHealthFake{}
	customCalls := 0
	sbx := &Sandbox{
		ID:               "sbx-ready",
		connectionConfig: configWithDefaults(),
		sandboxes:        lifecycle,
		Health:           health,
	}

	err := sbx.WaitUntilReady(context.Background(), &WaitUntilReadyOptions{
		Timeout:      2 * time.Second,
		PollInterval: 5 * time.Millisecond,
		CustomHealthCheck: func(ctx context.Context, sandbox *Sandbox) (bool, error) {
			customCalls++
			return true, nil
		},
	})
	if err != nil {
		t.Fatalf("wait until ready: %v", err)
	}
	if customCalls != 1 {
		t.Fatalf("expected custom health check call, got %d", customCalls)
	}
	if health.pingCallCount != 0 {
		t.Fatalf("expected default ping to be skipped, got %d calls", health.pingCallCount)
	}
}

func TestWaitUntilReadyReturnsTimeoutError(t *testing.T) {
	sbx := &Sandbox{
		ID:               "sbx-ready",
		connectionConfig: configWithDefaults(),
		sandboxes: &sandboxLifecycleFake{
			getSandboxResponses: []*models.SandboxInfo{
				{Status: models.SandboxStatus{State: "Running"}},
			},
		},
		Health: &sandboxHealthFake{results: []bool{false, false, false}},
	}

	err := sbx.WaitUntilReady(context.Background(), &WaitUntilReadyOptions{
		Timeout:      20 * time.Millisecond,
		PollInterval: 5 * time.Millisecond,
	})
	var timeoutErr *sandboxerrors.ReadyTimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected timeout error, got %v", err)
	}
}
