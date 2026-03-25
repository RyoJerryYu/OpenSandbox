package unit

import (
	"testing"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
)

func TestFactoryExposesLifecycleExecdAndEgressStacks(t *testing.T) {
	var _ factory.AdapterFactory = (*factory.DefaultAdapterFactory)(nil)
}

func TestDefaultFactoryFailsWithoutInjectedLifecycleClient(t *testing.T) {
	_, err := (&factory.DefaultAdapterFactory{}).CreateLifecycleStack(factory.CreateLifecycleStackOptions{})
	if err == nil {
		t.Fatal("expected missing lifecycle client error")
	}
}
