package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestNewRemote(t *testing.T) {
	t.Parallel()
	tests := map[string]gitdomain.Remote{
		"origin":   gitdomain.RemoteOrigin,
		"upstream": gitdomain.RemoteUpstream,
		"":         gitdomain.RemoteNone,
		"foo":      gitdomain.RemoteOther,
	}
	for give, want := range tests {
		have := gitdomain.NewRemote(give)
		must.Eq(t, want, have)
	}
}
