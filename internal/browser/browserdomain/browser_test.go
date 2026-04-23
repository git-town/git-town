package browserdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestUseBrowser(t *testing.T) {
	t.Parallel()
	tests := map[Option[browserdomain.Browser]]bool{
		None[browserdomain.Browser]():         true,
		Some(browserdomain.Browser("")):       false,
		Some(browserdomain.Browser("(none)")): false,
		Some(browserdomain.Browser("chrome")): true,
	}
	for give, want := range tests {
		have := browserdomain.UseBrowser(give)
		must.EqOp(t, want, have)
	}
}
