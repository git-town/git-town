package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestConfigStorage(t *testing.T) {
	t.Parallel()

	t.Run("Short", func(t *testing.T) {
		t.Parallel()
		must.EqOp(t, "file", dialog.ConfigStorageOptionFile.Short())
		must.EqOp(t, "git", dialog.ConfigStorageOptionGit.Short())
	})
}
