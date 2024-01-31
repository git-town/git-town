package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/shoenig/test/must"
)

func TestConfigStorage(t *testing.T) {
	t.Parallel()
	t.Run("Short", func(t *testing.T) {
		t.Parallel()
		tests := map[dialog.ConfigStorageOption]string{
			dialog.ConfigStorageOptionFile: "file",
			dialog.ConfigStorageOptionGit:  "git",
		}
		for give, want := range tests {
			have := give.Short()
			must.EqOp(t, want, have)
		}
	})
}
