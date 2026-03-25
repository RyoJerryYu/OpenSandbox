package unit

import (
	"testing"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
)

func TestFactoryExposesLifecycleExecdAndEgressStacks(t *testing.T) {
	var _ factory.AdapterFactory = (*factory.DefaultAdapterFactory)(nil)
}
