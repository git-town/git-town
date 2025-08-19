package envconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/shoenig/test/must"
)

func TestEnvVars(t *testing.T) {
	t.Parallel()

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		t.Run("lookup by name", func(t *testing.T) {
			t.Parallel()
			envVars := envconfig.NewEnvVars([]string{"SYNC_TAGS=yes"})
			have := envVars.Get("SYNC_TAGS")
			must.EqOp(t, "yes", have)
		})
	})
}
