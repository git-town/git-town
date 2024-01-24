package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestCodeHostingPlatforms(t *testing.T) {
	t.Parallel()

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		give := configdomain.CodeHostingPlatforms{
			configdomain.CodeHostingPlatformGitHub,
			configdomain.CodeHostingPlatformGitLab,
		}
		have := give.Strings()
		want := []string{"github", "gitlab"}
		must.Eq(t, want, have)
	})
}
