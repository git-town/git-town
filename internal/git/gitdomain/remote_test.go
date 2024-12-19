package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestNewRemote(t *testing.T) {
	t.Parallel()
	tests := map[string]Option[gitdomain.Remote]{
		"origin":   Some(gitdomain.RemoteOrigin),
		"upstream": Some(gitdomain.RemoteUpstream),
		"foo":      gitdomain.NewRemote("foo"),
		"":         None[gitdomain.Remote](),
	}
	for give, want := range tests {
		have := gitdomain.NewRemote(give)
		must.Eq(t, want, have)
	}
}
