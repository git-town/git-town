package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.ValidatedConfigData{
			MainBranch: "main",
		}
		must.False(t, config.IsMainBranch("feature"))
		must.True(t, config.IsMainBranch("main"))
		must.False(t, config.IsMainBranch("peren1"))
		must.False(t, config.IsMainBranch("peren2"))
	})
}
