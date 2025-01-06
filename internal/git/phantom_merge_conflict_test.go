package git_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/git"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseLsFilesUnmergedLine(t *testing.T) {
	t.Parallel()

	t.Run("base stage", func(t *testing.T) {
		t.Parallel()
		give := `100755 9f8f8acb41baba910c147c21eb61c55cf6d0447b 1	file`
		haveBlobInfo, haveStage, havePath, err := git.ParseLsFilesUnmergedLine(give)
		must.NoError(t, err)
		wantBlobInfo := git.BlobInfo{
			FilePath:   "file",
			Permission: "100755",
			SHA:        "9f8f8acb41baba910c147c21eb61c55cf6d0447b",
		}
		must.Eq(t, wantBlobInfo, haveBlobInfo)
		wantStage := git.UnmergedStageBase
		must.Eq(t, wantStage, haveStage)
		wantPath := "file"
		must.Eq(t, wantPath, havePath)
	})

	t.Run("current branch stage", func(t *testing.T) {
		t.Parallel()
		give := `100755 9f8f8acb41baba910c147c21eb61c55cf6d0447b 2	file`
		haveBlobInfo, haveStage, havePath, err := git.ParseLsFilesUnmergedLine(give)
		must.NoError(t, err)
		wantBlobInfo := git.BlobInfo{
			FilePath:   "file",
			Permission: "100755",
			SHA:        "9f8f8acb41baba910c147c21eb61c55cf6d0447b",
		}
		must.Eq(t, wantBlobInfo, haveBlobInfo)
		wantStage := git.UnmergedStageCurrentBranch
		must.Eq(t, wantStage, haveStage)
		wantPath := "file"
		must.Eq(t, wantPath, havePath)
	})

	t.Run("incoming stage", func(t *testing.T) {
		t.Parallel()
		give := `100755 9f8f8acb41baba910c147c21eb61c55cf6d0447b 3	file`
		haveBlobInfo, haveStage, havePath, err := git.ParseLsFilesUnmergedLine(give)
		must.NoError(t, err)
		wantBlobInfo := git.BlobInfo{
			FilePath:   "file",
			Permission: "100755",
			SHA:        "9f8f8acb41baba910c147c21eb61c55cf6d0447b",
		}
		must.Eq(t, wantBlobInfo, haveBlobInfo)
		wantStage := git.UnmergedStageIncoming
		must.Eq(t, wantStage, haveStage)
		wantPath := "file"
		must.Eq(t, wantPath, havePath)
	})
}

func TestParseLsFilesUnmergedOutput(t *testing.T) {
	t.Parallel()

	t.Run("conflicting changes", func(t *testing.T) {
		t.Parallel()
		give := `
			100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
			100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file`
		have, err := git.ParseLsFilesUnmergedOutput(give)
		want := []git.FileConflictQuickInfo{
			{
				BaseChange: None[git.BlobInfo](),
				CurrentBranchChange: Some(git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "c887ff2255bb9e9440f9456bcf8d310bc8d718d4",
				}),
				IncomingChange: Some(git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "ece1e56bf2125e5b114644258872f04bc375ba69",
				}),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("file deleted on current branch", func(t *testing.T) {
		t.Parallel()
		give := `
			100755 9f8f8acb41baba910c147c21eb61c55cf6d0447b 1	file
			100755 554e589880fc9e46b8b313499d325337187b1ee1 3	file`
		have, err := git.ParseLsFilesUnmergedOutput(give)
		want := []git.FileConflictQuickInfo{
			{
				BaseChange: Some(git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "9f8f8acb41baba910c147c21eb61c55cf6d0447b",
				}),
				CurrentBranchChange: None[git.BlobInfo](),
				IncomingChange: Some(git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "554e589880fc9e46b8b313499d325337187b1ee1",
				}),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})

	t.Run("file deleted on incoming branch", func(t *testing.T) {
		t.Parallel()
		give := `
			100755 9f8f8acb41baba910c147c21eb61c55cf6d0447b 1	file
			100755 554e589880fc9e46b8b313499d325337187b1ee1 2	file`
		have, err := git.ParseLsFilesUnmergedOutput(give)
		want := []git.FileConflictQuickInfo{
			{
				BaseChange: Some(git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "9f8f8acb41baba910c147c21eb61c55cf6d0447b",
				}),
				CurrentBranchChange: Some(git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "554e589880fc9e46b8b313499d325337187b1ee1",
				}),
				IncomingChange: None[git.BlobInfo](),
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})
}
