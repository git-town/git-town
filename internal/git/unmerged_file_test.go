package git_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestUnmergedFile(t *testing.T) {
	t.Parallel()

	t.Run("ParseLsFilesUnmergedOutput", func(t *testing.T) {
		t.Parallel()
		give := `
100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file
`[1:]
		have, err := git.ParseLsFilesUnmergedOutput(give)
		must.NoError(t, err)
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
		must.Eq(t, want, have)
	})

	t.Run("ParseLsTreeOutput", func(t *testing.T) {
		t.Parallel()
		t.Run("happy path", func(t *testing.T) {
			t.Parallel()
			give := `100755 blob ece1e56bf2125e5b114644258872f04bc375ba69	file`
			have, err := git.ParseLsTreeOutput(give)
			must.NoError(t, err)
			want := git.BlobInfo{
				FilePath:   "file",
				Permission: "100755",
				SHA:        gitdomain.NewSHA("ece1e56bf2125e5b114644258872f04bc375ba69"),
			}
			must.Eq(t, want, have)
		})
	})
}
