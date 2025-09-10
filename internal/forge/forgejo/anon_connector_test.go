package forgejo_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/forgejo"
	"github.com/shoenig/test/must"
)

func TestAnonConnector(t *testing.T) {
	t.Run("RepositoryURL", func(t *testing.T) {
		t.Parallel()
		connector := forgejo.AnonConnector{
			Data: forgedomain.Data{
				Hostname:     "codeberg.org",
				Organization: "org",
				Repository:   "repo",
			},
		}
		have := connector.RepositoryURL()
		must.EqOp(t, "https://codeberg.org/org/repo", have)
	})
}
