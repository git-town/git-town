package dialog

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestSwitchBranch(t *testing.T) {
	t.Parallel()
	t.Run("addEntryAndChildren", func(t *testing.T) {
		t.Parallel()
		t.Run("add an entry with children to an empty list", func(t *testing.T) {
			entries := ModalSelectEntries{}
			mainBranch := domain.NewLocalBranchName("main")
			feature1 := domain.NewLocalBranchName("feature-1")
			feature2 := domain.NewLocalBranchName("feature-2")
			lineage := config.Lineage{
				feature1: mainBranch,
				feature2: feature1,
			}
			have, err := addEntryAndChildren(entries, feature1, 0, lineage)
			must.NoError(t, err)
			want := ModalSelectEntries{ModalSelectEntry{
				Text:  "feature-1",
				Value: "feature-1",
			}, ModalSelectEntry{
				Text:  "  feature-2",
				Value: "feature-2",
			}}
			must.Eq(t, want, have)
		})
		t.Run("add an entry to an existing list", func(t *testing.T) {
			t.Parallel()
			entries := ModalSelectEntries{ModalSelectEntry{
				Text:  "existing",
				Value: "existing",
			}}
			mainBranch := domain.NewLocalBranchName("main")
			feature1 := domain.NewLocalBranchName("feature-1")
			feature2 := domain.NewLocalBranchName("feature-2")
			lineage := config.Lineage{
				feature1: mainBranch,
				feature2: feature1,
			}
			have, err := addEntryAndChildren(entries, feature1, 1, lineage)
			must.NoError(t, err)
			want := ModalSelectEntries{
				ModalSelectEntry{
					Text:  "existing",
					Value: "existing",
				},
				ModalSelectEntry{
					Text:  "  feature-1",
					Value: "feature-1",
				}, ModalSelectEntry{
					Text:  "    feature-2",
					Value: "feature-2",
				}}
			must.Eq(t, want, have)
		})
	})
	t.Run("createEntries", func(t *testing.T) {
		t.Parallel()
		mainBranch := domain.NewLocalBranchName("main")
		feature1 := domain.NewLocalBranchName("feature-1")
		feature2 := domain.NewLocalBranchName("feature-2")
		lineage := config.Lineage{
			feature1: mainBranch,
			feature2: feature1,
		}
		have, err := createEntries(lineage, feature2)
		must.NoError(t, err)
		want := ModalSelectEntries{
			ModalSelectEntry{
				Text:  "main",
				Value: "main",
			},
			ModalSelectEntry{
				Text:  "  feature-1",
				Value: "feature-1",
			}, ModalSelectEntry{
				Text:  "    feature-2",
				Value: "feature-2",
			}}
		must.Eq(t, want, have)
	})
}
