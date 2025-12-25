package git_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestParseGitStatusZ(t *testing.T) {
	t.Parallel()

	t.Run("copied file", func(t *testing.T) {
		t.Parallel()
		give := "C  vendor/golang.org/x/mod/PATENTS\x00vendor/golang.org/x/exp/PATENTS\x00"
		have, err := git.ParseGitStatusZ(give)
		must.NoError(t, err)
		want := []git.FileStatus{
			{
				OriginalPath: Some("vendor/golang.org/x/exp/PATENTS"),
				Path:         "vendor/golang.org/x/mod/PATENTS",
				ShortStatus:  "C ",
			},
		}
		must.Eq(t, want, have)
	})

	t.Run("one file", func(t *testing.T) {
		t.Parallel()
		give := " M internal/git/parse_git_status_z.go\x00"
		have, err := git.ParseGitStatusZ(give)
		must.NoError(t, err)
		want := []git.FileStatus{
			{
				OriginalPath: None[string](),
				Path:         "internal/git/parse_git_status_z.go",
				ShortStatus:  " M",
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
				OriginalPath: Some("internal/git/parse_git_status.go"),
				Path:         "internal/git/parse_git_status_z.go",
				ShortStatus:  "R ",
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
				OriginalPath: None[string](),
				Path:         "internal/git/parse_git_status_z.go",
				ShortStatus:  " M",
			},
			{
				OriginalPath: None[string](),
				Path:         "internal/git/parse_git_status_z_test.go",
				ShortStatus:  " M",
			},
		}
		must.Eq(t, want, have)
	})
}
