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

	t.Run("Load and Save roundtrip", func(t *testing.T) {
		t.Parallel()
		workspaceDir := t.TempDir()
		give := mockproposals.MockProposals{
			{
				Body:   gitdomain.NewProposalBodyOpt("body 1"),
				Number: 1,
				Source: "branch1",
				Target: "main",
				Title:  "Title 1",
				URL:    "https://example.com/pr/1",
			},
			{
				Body:   None[gitdomain.ProposalBody](),
				Number: 2,
				Source: "branch2",
				Target: "main",
				Title:  "Title 2",
				URL:    "https://example.com/pr/2",
			},
		}
		mockproposals.Save(workspaceDir, give)
		have := mockproposals.Load(workspaceDir)
		must.Eq(t, give, have)
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
			asserts.NoError(os.WriteFile(proposalsFile, []byte(`[{"Number": 999}]`), 0o600))
			newProposals := mockproposals.MockProposals{
				{
					Body:   None[gitdomain.ProposalBody](),
					Number: 456,
					Source: "new-branch",
					Target: "main",
					Title:  "Test Proposal",
					URL:    "https://example.com/pr/456",
				},
			}
			mockproposals.Save(workspaceDir, newProposals)
			have := asserts.NoError1(os.ReadFile(proposalsFile))
			want := `
[
  {
    "Active": false,
    "Body": null,
    "MergeWithAPI": false,
    "Number": 456,
    "Source": "new-branch",
    "Target": "main",
    "Title": "Test Proposal",
    "URL": "https://example.com/pr/456"
  }
]`[1:]
			must.Eq(t, want, string(have))
		})
	})
}
