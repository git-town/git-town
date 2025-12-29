package mockproposals_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestPersistence(t *testing.T) {
	t.Parallel()

	t.Run("Load", func(t *testing.T) {
		t.Parallel()

		t.Run("file exists", func(t *testing.T) {
			t.Parallel()
			workspaceDir := t.TempDir()
			proposalsFile := mockproposals.FilePath(workspaceDir)
			content := `
[
  {
    "Body": "test body",
    "Number": 123,
    "Source": "feature-branch",
    "Target": "main",
    "Title": "Test Proposal",
    "URL": "https://example.com/pr/123"
  }
]`[1:]
			asserts.NoError(os.WriteFile(proposalsFile, []byte(content), 0o600))
			have := mockproposals.Load(workspaceDir)
			want := mockproposals.MockProposals{
				{
					Body:   gitdomain.NewProposalBodyOpt("test body"),
					Number: 123,
					Source: "feature-branch",
					Target: "main",
					Title:  "Test Proposal",
					URL:    "https://example.com/pr/123",
				},
			}
			must.Eq(t, want, have)
		})

		t.Run("file does not exist", func(t *testing.T) {
			t.Parallel()
			workspaceDir := t.TempDir()
			have := mockproposals.Load(workspaceDir)
			must.Len(t, 0, have)
		})

		t.Run("file without proposals", func(t *testing.T) {
			t.Parallel()
			workspaceDir := t.TempDir()
			proposalsFile := filepath.Join(workspaceDir, "proposals.json")
			asserts.NoError(os.WriteFile(proposalsFile, []byte("[]"), 0o600))
			have := mockproposals.Load(workspaceDir)
			must.Len(t, 0, have)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()

		t.Run("save and load", func(t *testing.T) {
			t.Parallel()
			workspaceDir := t.TempDir()
			give := mockproposals.MockProposals{
				{
					Body:   gitdomain.NewProposalBodyOpt("test body"),
					Source: "feature-branch",
					Number: 123,
					Target: "main",
					Title:  "Test Proposal",
					URL:    "https://example.com/pr/123",
				},
			}
			mockproposals.Save(workspaceDir, give)
			proposalsFile := filepath.Join(workspaceDir, "proposals.json")
			have := asserts.NoError1(os.ReadFile(proposalsFile))
			want := `
[
  {
    "Active": false,
    "Body": "test body",
    "MergeWithAPI": false,
    "Number": 123,
    "Source": "feature-branch",
    "Target": "main",
    "Title": "Test Proposal",
    "URL": "https://example.com/pr/123"
  }
]`[1:]
			must.Eq(t, want, string(have))
		})

		t.Run("save empty proposals", func(t *testing.T) {
			t.Parallel()
			workspaceDir := t.TempDir()
			give := mockproposals.MockProposals{}
			mockproposals.Save(workspaceDir, give)
			proposalsFile := filepath.Join(workspaceDir, "proposals.json")
			have := asserts.NoError1(os.ReadFile(proposalsFile))
			must.EqOp(t, "[]", string(have))
		})

		t.Run("overwrite existing file", func(t *testing.T) {
			t.Parallel()
			workspaceDir := t.TempDir()
			proposalsFile := filepath.Join(workspaceDir, "proposals.json")
			err := os.WriteFile(proposalsFile, []byte(`[{"Number": 999}]`), 0o600)
			must.NoError(t, err)

			newProposals := mockproposals.MockProposals{
				{
					Active:       false,
					Body:         None[gitdomain.ProposalBody](),
					MergeWithAPI: false,
					Number:       456,
					Source:       gitdomain.NewLocalBranchName("new-branch"),
					Target:       gitdomain.NewLocalBranchName("main"),
					Title:        gitdomain.ProposalTitle("New Proposal"),
					URL:          "https://example.com/pr/456",
				},
			}

			mockproposals.Save(workspaceDir, newProposals)

			content, err := os.ReadFile(proposalsFile)
			must.NoError(t, err)
			must.StrContains(t, string(content), `"Number": 456`)
			must.StrContains(t, string(content), `"Source": "new-branch"`)
			must.StrNotContains(t, string(content), `"Number": 999`)
		})
	})

	t.Run("Load and Save roundtrip", func(t *testing.T) {
		t.Parallel()
		workspaceDir := t.TempDir()
		originalProposals := mockproposals.MockProposals{
			{
				Active:       true,
				Body:         gitdomain.NewProposalBodyOpt("body 1"),
				MergeWithAPI: true,
				Number:       1,
				Source:       gitdomain.NewLocalBranchName("branch1"),
				Target:       gitdomain.NewLocalBranchName("main"),
				Title:        gitdomain.ProposalTitle("Title 1"),
				URL:          "https://example.com/pr/1",
			},
			{
				Active:       false,
				Body:         None[gitdomain.ProposalBody](),
				MergeWithAPI: false,
				Number:       2,
				Source:       gitdomain.NewLocalBranchName("branch2"),
				Target:       gitdomain.NewLocalBranchName("main"),
				Title:        gitdomain.ProposalTitle("Title 2"),
				URL:          "https://example.com/pr/2",
			},
		}

		mockproposals.Save(workspaceDir, originalProposals)
		loadedProposals := mockproposals.Load(workspaceDir)

		must.Len(t, len(originalProposals), loadedProposals)
		for i, original := range originalProposals {
			must.Eq(t, original.Active, loadedProposals[i].Active)
			must.Eq(t, original.Body, loadedProposals[i].Body)
			must.Eq(t, original.MergeWithAPI, loadedProposals[i].MergeWithAPI)
			must.Eq(t, original.Number, loadedProposals[i].Number)
			must.Eq(t, original.Source, loadedProposals[i].Source)
			must.Eq(t, original.Target, loadedProposals[i].Target)
			must.Eq(t, original.Title, loadedProposals[i].Title)
			must.Eq(t, original.URL, loadedProposals[i].URL)
		}
	})
}
