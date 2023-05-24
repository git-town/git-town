package helpers_test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/git-town/git-town/v9/test/filesystem"
	"github.com/git-town/git-town/v9/test/helpers"
	"github.com/git-town/git-town/v9/test/testruntime"
	"github.com/stretchr/testify/assert"
)

func TestCopyDirectory(t *testing.T) {
	t.Parallel()
	t.Run("normal directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		srcDir := filepath.Join(tmpDir, "src")
		dstDir := filepath.Join(tmpDir, "dst")
		filesystem.CreateFile(t, srcDir, "one.txt")
		filesystem.CreateFile(t, srcDir, "f1/a.txt")
		filesystem.CreateFile(t, srcDir, "f2/b.txt")
		err := helpers.CopyDirectory(srcDir, dstDir)
		assert.NoError(t, err)
		asserts.FileExists(t, dstDir, "one.txt")
		asserts.FileExists(t, dstDir, "f1/a.txt")
		asserts.FileExists(t, dstDir, "f2/b.txt")
	})

	t.Run("Git repository", func(t *testing.T) {
		t.Parallel()
		origin := testruntime.Create(t)
		filesystem.CreateFile(t, origin.WorkingDir, "one.txt")
		dstDir := filepath.Join(t.TempDir(), "dest")
		err := helpers.CopyDirectory(origin.WorkingDir, dstDir)
		assert.NoError(t, err)
		asserts.FileExists(t, dstDir, "one.txt")
		asserts.FileHasContent(t, dstDir, ".git/HEAD", "ref: refs/heads/initial\n")
	})
}
