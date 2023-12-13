package common_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/hosting/common"
	"github.com/shoenig/test/must"
)

func TestHostnameWithStandardPort(t *testing.T) {
	t.Parallel()

	t.Run("no port in hostname", func(t *testing.T) {
		t.Parallel()
		config := common.Config{
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
		config := common.Config{
			Hostname:     "git.example.com:4022",
			Organization: "org",
			Repository:   "repo",
		}
		have := config.HostnameWithStandardPort()
		want := "git.example.com"
		must.EqOp(t, want, have)
	})
}
