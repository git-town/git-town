package runtime_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v8/test/asserts"
	"github.com/git-town/git-town/v8/test/runtime"
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
		err := runtime.CopyDirectory(srcDir, dstDir)
		assert.NoError(t, err)
		asserts.FileExists(t, dstDir, "one.txt")
		asserts.FileExists(t, dstDir, "f1/a.txt")
		asserts.FileExists(t, dstDir, "f2/b.txt")
	})

	t.Run("Git repository", func(t *testing.T) {
		t.Parallel()
		origin := runtime.Create(t)
		createFile(t, origin.Dir(), "one.txt")
		dstDir := filepath.Join(t.TempDir(), "dest")
		err := runtime.CopyDirectory(origin.Dir(), dstDir)
		assert.NoError(t, err)
		asserts.FileExists(t, dstDir, "one.txt")
		asserts.FileHasContent(t, dstDir, ".git/HEAD", "ref: refs/heads/initial\n")
	})
}

// createFile creates a file with the given filename in the given directory.
func createFile(t *testing.T, dir, filename string) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.MkdirAll(filepath.Dir(filePath), 0o744)
	assert.NoError(t, err)
	err = os.WriteFile(filePath, []byte(filename+" content"), 0o500)
	assert.NoError(t, err)
}
