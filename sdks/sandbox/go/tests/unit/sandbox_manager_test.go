package unit

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

func TestManagerRenewUsesNowPlusTimeout(t *testing.T) {
	fake := &fakeSandboxes{}
	manager, err := sandbox.NewSandboxManager(sandbox.SandboxManagerOptions{
		AdapterFactory: &fakeFactory{lifecycle: &factory.LifecycleStack{Sandboxes: fake}},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	before := time.Now()
	if err := manager.RenewSandbox(context.Background(), "sbx-1", 2*time.Minute); err != nil {
		t.Fatalf("renew sandbox: %v", err)
	}

	if fake.renewSandboxID != "sbx-1" {
		t.Fatalf("unexpected sandbox id: %s", fake.renewSandboxID)
	}

	minExpected := before.Add(2*time.Minute - 2*time.Second)
	maxExpected := before.Add(2*time.Minute + 2*time.Second)
	if fake.renewExpiresAt.Before(minExpected) || fake.renewExpiresAt.After(maxExpected) {
		t.Fatalf("unexpected renew time: %s", fake.renewExpiresAt)
	}
}

func TestManagerCloseReleasesLocalResourcesOnly(t *testing.T) {
	manager, err := sandbox.NewSandboxManager(sandbox.SandboxManagerOptions{
		AdapterFactory: &fakeFactory{lifecycle: &factory.LifecycleStack{Sandboxes: &fakeSandboxes{}}},
		CloseFunc: func() error {
			return nil
		},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	if err := manager.Close(); err != nil {
		t.Fatalf("close manager: %v", err)
	}
}

func TestNewSandboxManagerFailsWhenLifecycleStackCannotBeBuilt(t *testing.T) {
	expected := errors.New("missing lifecycle client")
	_, err := sandbox.NewSandboxManager(sandbox.SandboxManagerOptions{
		AdapterFactory: &failingFactory{err: expected},
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected lifecycle factory error, got %v", err)
	}
}

func TestNewSandboxManagerUsesDefaultFactoryWithGeneratedLifecycleClient(t *testing.T) {
	var gotAPIKey string
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			gotAPIKey = r.Header.Get("OPEN-SANDBOX-API-KEY")
			body := []byte(`{"items":[],"pagination":{"page":1,"pageSize":10,"totalItems":0,"totalPages":0,"hasNextPage":false}}`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
		Timeout: 5 * time.Second,
	}

	manager, err := sandbox.NewSandboxManager(sandbox.SandboxManagerOptions{
		ConnectionConfig: &config.ConnectionConfig{
			Domain:     "example.test",
			Protocol:   "http",
			APIKey:     "api-key-1",
			HTTPClient: client,
		},
	})
	if err != nil {
		t.Fatalf("new manager: %v", err)
	}

	if _, err := manager.ListSandboxInfos(context.Background(), models.SandboxFilter{}); err != nil {
		t.Fatalf("list sandbox infos: %v", err)
	}

	if gotAPIKey != "api-key-1" {
		t.Fatalf("unexpected api key header: %q", gotAPIKey)
	}
}

type fakeFactory struct {
	lifecycle *factory.LifecycleStack
}

func (f *fakeFactory) CreateLifecycleStack(opts factory.CreateLifecycleStackOptions) (*factory.LifecycleStack, error) {
	return f.lifecycle, nil
}

func (f *fakeFactory) CreateExecdStack(opts factory.CreateExecdStackOptions) (*factory.ExecdStack, error) {
	return &factory.ExecdStack{}, nil
}

func (f *fakeFactory) CreateEgressStack(opts factory.CreateEgressStackOptions) (*factory.EgressStack, error) {
	return &factory.EgressStack{}, nil
}

type failingFactory struct {
	err error
}

func (f *failingFactory) CreateLifecycleStack(opts factory.CreateLifecycleStackOptions) (*factory.LifecycleStack, error) {
	return nil, f.err
}

func (f *failingFactory) CreateExecdStack(opts factory.CreateExecdStackOptions) (*factory.ExecdStack, error) {
	return nil, f.err
}

func (f *failingFactory) CreateEgressStack(opts factory.CreateEgressStackOptions) (*factory.EgressStack, error) {
	return nil, f.err
}

type fakeSandboxes struct {
	renewSandboxID string
	renewExpiresAt time.Time
}

func (f *fakeSandboxes) CreateSandbox(context.Context, models.CreateSandboxRequest) (*models.CreateSandboxResponse, error) {
	return nil, nil
}

func (f *fakeSandboxes) GetSandbox(context.Context, string) (*models.SandboxInfo, error) {
	return nil, nil
}

func (f *fakeSandboxes) ListSandboxes(context.Context, models.SandboxFilter) (*models.ListSandboxesResponse, error) {
	return nil, nil
}

func (f *fakeSandboxes) DeleteSandbox(context.Context, string) error {
	return nil
}

func (f *fakeSandboxes) PauseSandbox(context.Context, string) error {
	return nil
}

func (f *fakeSandboxes) ResumeSandbox(context.Context, string) error {
	return nil
}

func (f *fakeSandboxes) RenewSandboxExpiration(ctx context.Context, sandboxID string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error) {
	f.renewSandboxID = sandboxID
	f.renewExpiresAt = expiresAt
	return &models.RenewSandboxExpirationResponse{ExpiresAt: &expiresAt}, nil
}

func (f *fakeSandboxes) GetSandboxEndpoint(context.Context, string, int, bool) (*models.SandboxEndpoint, error) {
	return nil, nil
}

var _ services.Sandboxes = (*fakeSandboxes)(nil)
