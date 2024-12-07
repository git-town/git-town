package git_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/git"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseLsFilesUnmergedOutput(t *testing.T) {
	t.Parallel()
	t.Run("conflicting changes", func(t *testing.T) {
		t.Parallel()
		give := `
	100755 c887ff2255bb9e9440f9456bcf8d310bc8d718d4 2	file
	100755 ece1e56bf2125e5b114644258872f04bc375ba69 3	file`[1:]
		have, err := git.ParseLsFilesUnmergedOutput(give)
		want := []git.FileConflictQuickInfo{
			{
				BaseChange: None[git.BlobInfo](),
				CurrentBranchChange: git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "c887ff2255bb9e9440f9456bcf8d310bc8d718d4",
				},
				IncomingChange: git.BlobInfo{
					FilePath:   "file",
					Permission: "100755",
					SHA:        "ece1e56bf2125e5b114644258872f04bc375ba69",
				},
			},
		}
		must.NoError(t, err)
		must.Eq(t, want, have)
	})
}
