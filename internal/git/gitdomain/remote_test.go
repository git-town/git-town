package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestNewRemote(t *testing.T) {
	t.Parallel()
	tests := map[stringss.TrimmedString]Option[gitdomain.Remote]{
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
