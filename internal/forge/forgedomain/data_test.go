package forgedomain_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/shoenig/test/must"
)

func TestHostnameWithStandardPort(t *testing.T) {
	t.Parallel()

	t.Run("no port in hostname", func(t *testing.T) {
		t.Parallel()
		config := forgedomain.Data{
			Hostname:     "git.example.com",
			Organization: "org",
			Repository:   "repo",
		}
		have := config.HostnameWithStandardPort()
		want := "git.example.com"
		must.EqOp(t, want, have)
	})

	t.Run("port in hostname", func(t *testing.T) {
		t.Parallel()
		config := forgedomain.Data{
			Hostname:     "git.example.com:4022",
			Organization: "org",
			Repository:   "repo",
		}
		have := config.HostnameWithStandardPort()
		want := "git.example.com"
		must.EqOp(t, want, have)
	})
}
