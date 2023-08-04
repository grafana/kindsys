package themasys_test

import (
	"testing"

	"github.com/grafana/kindsys/pkg/themasys"
)

func TestFramework(t *testing.T) {
	// please don't panic, that's all I ask
	_ = themasys.CUEFramework(nil)
}
