package dialog_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestBuilder(t *testing.T) {
	t.Parallel()

	t.Run("AddEntryAndChildren", func(t *testing.T) {
		t.Parallel()
		t.Run("add an entry with children to an empty list", func(t *testing.T) {
			main := domain.NewLocalBranchName("main")
			feature1 := domain.NewLocalBranchName("feature-1")
			feature2 := domain.NewLocalBranchName("feature-2")
			featureA := domain.NewLocalBranchName("feature-A")
			lineage := configdomain.Lineage{
				feature1: main,
				feature2: feature1,
				featureA: main,
			}
			builder := dialog.NewBuilder(lineage)
			// add entries to the empty builder
			err := builder.AddEntryAndChildren(feature1, 0)
			must.NoError(t, err)
			want := dialog.ModalSelectEntries{dialog.ModalSelectEntry{
				Text:  "feature-1",
				Value: "feature-1",
			}, dialog.ModalSelectEntry{
				Text:  "  feature-2",
				Value: "feature-2",
			}}
			must.Eq(t, want, builder.Entries)
			// add more entries to the already populated builder
			err = builder.AddEntryAndChildren(featureA, 0)
			must.NoError(t, err)
			want = dialog.ModalSelectEntries{dialog.ModalSelectEntry{
				Text:  "feature-1",
				Value: "feature-1",
			}, dialog.ModalSelectEntry{
				Text:  "  feature-2",
				Value: "feature-2",
			}, dialog.ModalSelectEntry{
				Text:  "feature-A",
				Value: "feature-A",
			}}
			must.Eq(t, want, builder.Entries)
		})
	})

	t.Run("CreateEntries", func(t *testing.T) {
		t.Parallel()
		mainBranch := domain.NewLocalBranchName("main")
		feature1 := domain.NewLocalBranchName("feature-1")
		feature2 := domain.NewLocalBranchName("feature-2")
		lineage := configdomain.Lineage{
			feature1: mainBranch,
			feature2: feature1,
		}
		builder := dialog.NewBuilder(lineage)
		roots := domain.LocalBranchNames{mainBranch}
		err := builder.CreateEntries(roots, feature1)
		must.NoError(t, err)
		want := dialog.ModalSelectEntries{
			dialog.ModalSelectEntry{
				Text:  "main",
				Value: "main",
			},
			dialog.ModalSelectEntry{
				Text:  "  feature-1",
				Value: "feature-1",
			},
			dialog.ModalSelectEntry{
				Text:  "    feature-2",
				Value: "feature-2",
			},
		}
		must.Eq(t, want, builder.Entries)
	})
}
