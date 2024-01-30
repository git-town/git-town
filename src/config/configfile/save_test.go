package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/configfile"
	"github.com/shoenig/test/must"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("RenderTOML", func(t *testing.T) {
		t.Parallel()
		give := configdomain.DefaultConfig()
		have := configfile.RenderTOML(&give)
		want := `
push-hook = true
push-new-branches = false
ship-delete-tracking-branch = false
sync-before-ship = false
sync-upstream = true

[branches]
  main = "main"
  perennials = ["one", "two"]

[sync-strategy]
  feature-branches = "merge"
  perennial-branches = "rebase"
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()
		give := configdomain.DefaultConfig()
		err := configfile.Save(&give)
		defer os.Remove(configfile.FileName)
		must.NoError(t, err)
		bytes, err := os.ReadFile(configfile.FileName)
		must.NoError(t, err)
		have := string(bytes)
		want := `
push-hook = true
push-new-branches = false
ship-delete-tracking-branch = false
sync-before-ship = false
sync-upstream = true

[branches]
  main = "main"
  perennials = ["one", "two"]

[sync-strategy]
  feature-branches = "merge"
  perennial-branches = "rebase"
`[1:]
		must.EqOp(t, want, have)
	})
}
