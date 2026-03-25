package factory

import (
	"errors"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
)

type DefaultAdapterFactory struct{}

func (f *DefaultAdapterFactory) CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error) {
	if opts.LifecycleClient == nil {
		return nil, errors.New("missing lifecycle client")
	}
	return &LifecycleStack{
		Sandboxes: adapters.NewSandboxesAdapter(opts.LifecycleClient),
	}, nil
}

func (f *DefaultAdapterFactory) CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error) {
	return &ExecdStack{}, nil
}

func (f *DefaultAdapterFactory) CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error) {
	return &EgressStack{}, nil
}
