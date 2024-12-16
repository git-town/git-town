package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/git-town/git-town/v16/test/git"
	"github.com/shoenig/test/must"
)

func TestNewRemote(t *testing.T) {
	t.Parallel()
	tests := map[string]Option[gitdomain.Remote]{
		"origin":   Some(git.RemoteOrigin),
		"upstream": Some(gitdomain.RemoteUpstream),
		"foo":      gitdomain.NewRemote("foo"),
		"":         None[gitdomain.Remote](),
	}
	for give, want := range tests {
		have := gitdomain.NewRemote(give)
		must.Eq(t, want, have)
	}
}
