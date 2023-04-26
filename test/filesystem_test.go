//nolint:testpackage
package test

import (
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/test/asserts"
	"github.com/stretchr/testify/assert"
)

func TestCopyDirectory(t *testing.T) {
	t.Parallel()
	t.Run("normal directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := t.TempDir()
		srcDir := filepath.Join(tmpDir, "src")
		dstDir := filepath.Join(tmpDir, "dst")
		createFile(t, srcDir, "one.txt")
		createFile(t, srcDir, "f1/a.txt")
		createFile(t, srcDir, "f2/b.txt")
		err := CopyDirectory(srcDir, dstDir)
		assert.NoError(t, err)
		asserts.FileExists(t, dstDir, "one.txt")
		asserts.FileExists(t, dstDir, "f1/a.txt")
		asserts.FileExists(t, dstDir, "f2/b.txt")
	})

	t.Run("Git repository", func(t *testing.T) {
		t.Parallel()
		origin := CreateRunner(t)
		createFile(t, origin.WorkingDir(), "one.txt")
		dstDir := filepath.Join(t.TempDir(), "dest")
		err := CopyDirectory(origin.WorkingDir(), dstDir)
		assert.NoError(t, err)
		asserts.FileExists(t, dstDir, "one.txt")
		asserts.FileHasContent(t, dstDir, ".git/HEAD", "ref: refs/heads/initial\n")
	})
}
