//nolint:testpackage
package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyDirectory(t *testing.T) {
	t.Run("normal directory", func(t *testing.T) {
		t.Parallel()
		tmpDir := CreateTempDir(t)
		srcDir := filepath.Join(tmpDir, "src")
		dstDir := filepath.Join(tmpDir, "dst")
		createFile(t, srcDir, "one.txt")
		createFile(t, srcDir, "f1/a.txt")
		createFile(t, srcDir, "f2/b.txt")
		err := CopyDirectory(srcDir, dstDir)
		assert.NoError(t, err)
		assertFileExists(t, dstDir, "one.txt")
		assertFileExists(t, dstDir, "f1/a.txt")
		assertFileExists(t, dstDir, "f2/b.txt")
	})

	t.Run("Git repository", func(t *testing.T) {
		t.Parallel()
		origin := CreateRepo(t)
		createFile(t, origin.WorkingDir(), "one.txt")
		dstDir := filepath.Join(CreateTempDir(t), "dest")
		err := CopyDirectory(origin.WorkingDir(), dstDir)
		assert.NoError(t, err)
		assertFileExists(t, dstDir, "one.txt")
		assertFileExistsWithContent(t, dstDir, ".git/HEAD", "ref: refs/heads/initial\n")
	})
}
