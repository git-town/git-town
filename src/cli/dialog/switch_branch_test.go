package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()
	t.Run("addEntryAndChildren", func(t *testing.T) {
		t.Parallel()
		entries := dialog.ModalSelectEntries{dialog.ModalSelectEntry{
			Text:  "One",
			Value: "1",
		}, dialog.ModalSelectEntry{
			Text:  "Two",
			Value: "2",
		}}
		branch := domain.NewLocalBranchName("new-branch")
		mainBranch := domain.NewLocalBranchName("main")
		feature1 := domain.NewLocalBranchName("feature-1")
		lineage := config.Lineage{
			feature1: mainBranch,
		}
		have, err := dialog.AddEntryAndChildren(entries, branch, 2, lineage)
		must.NoError(t, err)
		want := dialog.ModalSelectEntries{dialog.ModalSelectEntry{
			Text:  "One",
			Value: "1",
		}, dialog.ModalSelectEntry{
			Text:  "Two",
			Value: "2",
		}, dialog.ModalSelectEntry{
			Text:  "    new-branch",
			Value: "new-branch",
		}}
		must.Eq(t, want, have)
	})
}
