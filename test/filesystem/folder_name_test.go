package filesystem_test

import (
	"testing"

	"github.com/git-town/git-town/v13/test/filesystem"
	"github.com/shoenig/test/must"
)

func TestFolderName(t *testing.T) {
	t.Parallel()
	tests := map[string]string{
		"foo":                                 "foo",
		`globally set to "true", local unset`: "globally_set_to_true_local_unset",
	}
	for give, want := range tests {
		have := filesystem.FolderName(give)
		must.EqOp(t, want, have)
	}
}
