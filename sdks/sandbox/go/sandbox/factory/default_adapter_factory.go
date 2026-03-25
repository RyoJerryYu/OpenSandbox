package factory

import (
	"context"
	"errors"
	"net/http"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
	lifecycleapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/lifecycle"
)

type DefaultAdapterFactory struct{}

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

func (f *DefaultAdapterFactory) CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error) {
	return &ExecdStack{}, nil
}

func (f *DefaultAdapterFactory) CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error) {
	return &EgressStack{}, nil
}
