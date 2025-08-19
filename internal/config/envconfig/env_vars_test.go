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
		t.Run("lookup by alternative name", func(t *testing.T) {
			t.Parallel()
			envVars := envconfig.NewEnvVars([]string{"SYNC_TAGS=yes"})
			have := envVars.Get("OTHER_NAME", "OTHER_NAME_2", "SYNC_TAGS")
			must.EqOp(t, "yes", have)
		})
		t.Run("entry doesn't exist", func(t *testing.T) {
			t.Parallel()
			envVars := envconfig.NewEnvVars([]string{"SYNC_TAGS=yes"})
			have := envVars.Get("OTHER_NAME")
			must.EqOp(t, "", have)
		})
	})
}
