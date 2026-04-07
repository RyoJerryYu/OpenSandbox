package factory

import (
	"context"
	"errors"
	"net/http"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	egressapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/egress"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	lifecycleapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/lifecycle"
)

// DefaultAdapterFactory builds SDK adapters backed by the generated OpenAPI clients.
type DefaultAdapterFactory struct{}

// CreateLifecycleStack builds lifecycle adapters for sandbox management APIs.
func (f *DefaultAdapterFactory) CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error) {
	lifecycleClient := opts.LifecycleClient
	if lifecycleClient == nil {
		if opts.ConnectionConfig == nil || opts.LifecycleBaseURL == "" {
			return nil, errors.New("missing lifecycle client")
		}
		client, err := lifecycleapi.NewClientWithResponses(
			opts.LifecycleBaseURL,
			lifecycleapi.WithHTTPClient(opts.ConnectionConfig.HTTPClient),
			lifecycleapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
				if opts.ConnectionConfig.APIKey != "" {
					req.Header.Set("OPEN-SANDBOX-API-KEY", opts.ConnectionConfig.APIKey)
				}
				for k, v := range opts.ConnectionConfig.Headers {
					req.Header.Set(k, v)
				}
				return nil
			}),
		)
		if err != nil {
			return nil, err
		}
		lifecycleClient = adapters.NewGeneratedLifecycleClient(client)
	}
	return &LifecycleStack{
		Sandboxes: adapters.NewSandboxesAdapter(lifecycleClient),
	}, nil
}

// CreateExecdStack builds execd adapters for commands, files, health, and metrics.
func (f *DefaultAdapterFactory) CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error) {
	if opts.ConnectionConfig == nil || opts.ExecdBaseURL == "" {
		return nil, errors.New("missing execd client configuration")
	}

	client, err := execdapi.NewClientWithResponses(
		opts.ExecdBaseURL,
		execdapi.WithHTTPClient(opts.ConnectionConfig.HTTPClient),
		execdapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			return applyConnectionHeaders(req, opts.ConnectionConfig)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &ExecdStack{
		Commands: adapters.NewStreamingCommandsAdapter(client, opts.ExecdBaseURL, opts.ConnectionConfig),
		Files:    adapters.NewFilesystemAdapter(client, opts.ExecdBaseURL, opts.ConnectionConfig),
		Health:   adapters.NewHealthAdapter(client),
		Metrics:  adapters.NewMetricsAdapter(client),
	}, nil
}

// CreateEgressStack builds adapters for the sandbox egress sidecar.
func (f *DefaultAdapterFactory) CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error) {
	if opts.ConnectionConfig == nil || opts.EgressBaseURL == "" {
		return nil, errors.New("missing egress client configuration")
	}

	client, err := egressapi.NewClientWithResponses(
		opts.EgressBaseURL,
		egressapi.WithHTTPClient(opts.ConnectionConfig.HTTPClient),
		egressapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			return applyConnectionHeaders(req, opts.ConnectionConfig)
		}),
	)
	if err != nil {
		return nil, err
	}

	return &EgressStack{
		Egress: adapters.NewEgressAdapter(client),
	}, nil
}

func applyConnectionHeaders(req *http.Request, connectionConfig *config.ConnectionConfig) error {
	if connectionConfig == nil {
		return nil
	}
	if connectionConfig.APIKey != "" {
		req.Header.Set("OPEN-SANDBOX-API-KEY", connectionConfig.APIKey)
	}
	for k, v := range connectionConfig.Headers {
		req.Header.Set(k, v)
	}
	return nil
}
