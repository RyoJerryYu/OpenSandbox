package factory

type DefaultAdapterFactory struct{}

func (f *DefaultAdapterFactory) CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error) {
	return &LifecycleStack{}, nil
}

func (f *DefaultAdapterFactory) CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error) {
	return &ExecdStack{}, nil
}

func (f *DefaultAdapterFactory) CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error) {
	return &EgressStack{}, nil
}
