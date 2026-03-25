package sandbox

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestSandboxRenewUsesNowPlusTimeout(t *testing.T) {
	fake := &sandboxLifecycleFake{}
	sbx := &Sandbox{
		ID:               "sbx-1",
		connectionConfig: config.DefaultConnectionConfig(),
		sandboxes:        fake,
	}

	before := time.Now()
	if _, err := sbx.Renew(context.Background(), 90*time.Second); err != nil {
		t.Fatalf("renew sandbox: %v", err)
	}

	minExpected := before.Add(88 * time.Second)
	maxExpected := before.Add(92 * time.Second)
	if fake.renewExpiresAt.Before(minExpected) || fake.renewExpiresAt.After(maxExpected) {
		t.Fatalf("unexpected renew time: %s", fake.renewExpiresAt)
	}
}

func TestSandboxGetEndpointUsesConnectionProxyFlag(t *testing.T) {
	fake := &sandboxLifecycleFake{
		endpointResponse: &models.SandboxEndpoint{Endpoint: "domain/proxy/44772"},
	}
	sbx := &Sandbox{
		ID: "sbx-1",
		connectionConfig: &config.ConnectionConfig{
			Domain:         config.DefaultDomain,
			Protocol:       config.DefaultProtocol,
			UseServerProxy: true,
		},
		sandboxes: fake,
	}

	if _, err := sbx.GetEndpoint(context.Background(), 44772); err != nil {
		t.Fatalf("get endpoint: %v", err)
	}
	if !fake.useServerProxy {
		t.Fatal("expected useServerProxy to be forwarded")
	}
}

func TestSandboxCloseCallsCloseFunc(t *testing.T) {
	closed := false
	sbx := &Sandbox{
		closeFn: func() error {
			closed = true
			return nil
		},
	}

	if err := sbx.Close(); err != nil {
		t.Fatalf("close sandbox: %v", err)
	}
	if !closed {
		t.Fatal("expected close function to be called")
	}
}

func TestCreateResolvesExecdAndEgressEndpoints(t *testing.T) {
	lifecycle := &sandboxLifecycleFake{
		createResponse: &models.CreateSandboxResponse{ID: "sbx-1"},
		endpointResponseByPort: map[int]*models.SandboxEndpoint{
			defaultExecdPort:  {Endpoint: "execd.test:44772"},
			defaultEgressPort: {Endpoint: "egress.test:18080"},
		},
	}
	factorySpy := &sandboxFactoryFake{
		lifecycle: &factory.LifecycleStack{Sandboxes: lifecycle},
		execd:     &factory.ExecdStack{},
		egress:    &factory.EgressStack{},
	}

	timeout := 5 * time.Minute
	sbx, err := Create(context.Background(), SandboxCreateOptions{
		ConnectionConfig: config.DefaultConnectionConfig(),
		AdapterFactory:   factorySpy,
		Timeout:          &timeout,
		Image:            models.ImageSpec{URI: "ubuntu"},
	})
	if err != nil {
		t.Fatalf("create sandbox: %v", err)
	}

	if sbx.ID != "sbx-1" {
		t.Fatalf("unexpected sandbox id: %s", sbx.ID)
	}
	if factorySpy.execdBaseURL != "http://execd.test:44772" {
		t.Fatalf("unexpected execd base url: %s", factorySpy.execdBaseURL)
	}
	if factorySpy.egressBaseURL != "http://egress.test:18080" {
		t.Fatalf("unexpected egress base url: %s", factorySpy.egressBaseURL)
	}
}

func TestCreateCleansUpRemoteSandboxOnInitializationFailure(t *testing.T) {
	lifecycle := &sandboxLifecycleFake{
		createResponse: &models.CreateSandboxResponse{ID: "sbx-2"},
		endpointResponseByPort: map[int]*models.SandboxEndpoint{
			defaultExecdPort: {Endpoint: "execd.test:44772"},
		},
	}
	expected := errors.New("execd stack failed")
	factorySpy := &sandboxFactoryFake{
		lifecycle: &factory.LifecycleStack{Sandboxes: lifecycle},
		execdErr:  expected,
	}

	timeout := 5 * time.Minute
	_, err := Create(context.Background(), SandboxCreateOptions{
		ConnectionConfig: config.DefaultConnectionConfig(),
		AdapterFactory:   factorySpy,
		Timeout:          &timeout,
		Image:            models.ImageSpec{URI: "ubuntu"},
	})
	if !errors.Is(err, expected) {
		t.Fatalf("expected execd stack error, got %v", err)
	}
	if lifecycle.deletedSandboxID != "sbx-2" {
		t.Fatalf("expected remote cleanup for sbx-2, got %s", lifecycle.deletedSandboxID)
	}
}

func TestResumeReturnsFreshSandboxInstance(t *testing.T) {
	lifecycle := &sandboxLifecycleFake{
		endpointResponseByPort: map[int]*models.SandboxEndpoint{
			defaultExecdPort:  {Endpoint: "execd.test:44772"},
			defaultEgressPort: {Endpoint: "egress.test:18080"},
		},
	}
	factorySpy := &sandboxFactoryFake{
		lifecycle: &factory.LifecycleStack{Sandboxes: lifecycle},
		execd:     &factory.ExecdStack{},
		egress:    &factory.EgressStack{},
	}

	original := &Sandbox{
		ID:               "sbx-3",
		connectionConfig: config.DefaultConnectionConfig(),
		sandboxes:        lifecycle,
	}

	resumed, err := original.Resume(context.Background(), ResumeOptions{
		AdapterFactory: factorySpy,
	})
	if err != nil {
		t.Fatalf("resume sandbox: %v", err)
	}
	if resumed == original {
		t.Fatal("expected fresh sandbox instance")
	}
	if lifecycle.resumedSandboxID != "sbx-3" {
		t.Fatalf("expected resume call for sbx-3, got %s", lifecycle.resumedSandboxID)
	}
}

type sandboxLifecycleFake struct {
	createResponse      *models.CreateSandboxResponse
	endpointResponse    *models.SandboxEndpoint
	endpointResponseByPort map[int]*models.SandboxEndpoint
	renewExpiresAt     time.Time
	useServerProxy     bool
	deletedSandboxID   string
	resumedSandboxID   string
}

func (f *sandboxLifecycleFake) CreateSandbox(context.Context, models.CreateSandboxRequest) (*models.CreateSandboxResponse, error) {
	return f.createResponse, nil
}

func (f *sandboxLifecycleFake) GetSandbox(context.Context, string) (*models.SandboxInfo, error) {
	return nil, nil
}

func (f *sandboxLifecycleFake) ListSandboxes(context.Context, models.SandboxFilter) (*models.ListSandboxesResponse, error) {
	return nil, nil
}

func (f *sandboxLifecycleFake) DeleteSandbox(_ context.Context, sandboxID string) error {
	f.deletedSandboxID = sandboxID
	return nil
}

func (f *sandboxLifecycleFake) PauseSandbox(context.Context, string) error {
	return nil
}

func (f *sandboxLifecycleFake) ResumeSandbox(_ context.Context, sandboxID string) error {
	f.resumedSandboxID = sandboxID
	return nil
}

func (f *sandboxLifecycleFake) RenewSandboxExpiration(_ context.Context, _ string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error) {
	f.renewExpiresAt = expiresAt
	return &models.RenewSandboxExpirationResponse{ExpiresAt: &expiresAt}, nil
}

func (f *sandboxLifecycleFake) GetSandboxEndpoint(_ context.Context, _ string, port int, useServerProxy bool) (*models.SandboxEndpoint, error) {
	f.useServerProxy = useServerProxy
	if f.endpointResponseByPort != nil {
		return f.endpointResponseByPort[port], nil
	}
	return f.endpointResponse, nil
}

type sandboxFactoryFake struct {
	lifecycle     *factory.LifecycleStack
	execd         *factory.ExecdStack
	egress        *factory.EgressStack
	execdErr      error
	egressErr     error
	execdBaseURL  string
	egressBaseURL string
}

func (f *sandboxFactoryFake) CreateLifecycleStack(opts factory.CreateLifecycleStackOptions) (*factory.LifecycleStack, error) {
	return f.lifecycle, nil
}

func (f *sandboxFactoryFake) CreateExecdStack(opts factory.CreateExecdStackOptions) (*factory.ExecdStack, error) {
	f.execdBaseURL = opts.ExecdBaseURL
	if f.execdErr != nil {
		return nil, f.execdErr
	}
	return f.execd, nil
}

func (f *sandboxFactoryFake) CreateEgressStack(opts factory.CreateEgressStackOptions) (*factory.EgressStack, error) {
	f.egressBaseURL = opts.EgressBaseURL
	if f.egressErr != nil {
		return nil, f.egressErr
	}
	return f.egress, nil
}
