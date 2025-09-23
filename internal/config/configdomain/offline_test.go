package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestOffline(t *testing.T) {
	t.Parallel()
	tests := []bool{true, false}
	for _, give := range tests {
		offline := configdomain.Offline(give)
		must.EqOp(t, give, offline.IsOffline())
		must.EqOp(t, !give, offline.IsOnline())
	}
}
