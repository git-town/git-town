package git_test

import (
	"testing"

	"github.com/git-town/git-town/v20/internal/git"
	. "github.com/git-town/git-town/v20/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseGitStatusZ(t *testing.T) {
	t.Parallel()

	t.Run("one file", func(t *testing.T) {
		t.Parallel()
		give := " M internal/git/parse_git_status_z.go\x00"
		have, err := git.ParseGitStatusZ(give)
		must.NoError(t, err)
		want := []git.FileStatus{
			{
				ShortStatus: " M",
				Path:        "internal/git/parse_git_status_z.go",
				RenamedFrom: None[string](),
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("renamed file", func(t *testing.T) {
		t.Parallel()
		give := "R  internal/git/parse_git_status_z.go\x00internal/git/parse_git_status.go\x00"
		have, err := git.ParseGitStatusZ(give)
		must.NoError(t, err)
		want := []git.FileStatus{
			{
				ShortStatus: "R ",
				Path:        "internal/git/parse_git_status_z.go",
				RenamedFrom: Some("internal/git/parse_git_status.go"),
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("two files", func(t *testing.T) {
		t.Parallel()
		give := " M internal/git/parse_git_status_z.go\x00 M internal/git/parse_git_status_z_test.go\x00"
		have, err := git.ParseGitStatusZ(give)
		must.NoError(t, err)
		want := []git.FileStatus{
			{
				ShortStatus: " M",
				Path:        "internal/git/parse_git_status_z.go",
				RenamedFrom: None[string](),
			},
			{
				ShortStatus: " M",
				Path:        "internal/git/parse_git_status_z_test.go",
				RenamedFrom: None[string](),
			},
		}
		must.Eq(t, want, have)
	})
}
