package browserdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/browser/browserdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestUseBrowser(t *testing.T) {
	t.Parallel()
	tests := map[Option[browserdomain.BrowserExecutable]]bool{
		None[browserdomain.BrowserExecutable]():         true,
		Some(browserdomain.BrowserExecutable("")):       false,
		Some(browserdomain.BrowserExecutable("(none)")): false,
		Some(browserdomain.BrowserExecutable("chrome")): true,
	}
	for give, want := range tests {
		have := browserdomain.BrowserEnabled(give)
		must.EqOp(t, want, have)
	}
}
