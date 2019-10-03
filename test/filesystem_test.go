package test

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyDirectory(t *testing.T) {
	tmpDir := createTempDir(t)
	srcDir := path.Join(tmpDir, "src")
	dstDir := path.Join(tmpDir, "dst")
	createFile(t, srcDir, "one.txt")
	createFile(t, srcDir, "f1/a.txt")
	createFile(t, srcDir, "f2/b.txt")

	err := CopyDirectory(srcDir, dstDir)

	assert.Nil(t, err)
	assertFileExists(t, dstDir, "one.txt")
	assertFileExists(t, dstDir, "f1/a.txt")
	assertFileExists(t, dstDir, "f2/b.txt")
}

func TestCopyDirectoryGitRepo(t *testing.T) {
	tmpDir := createTempDir(t)
	srcDir := path.Join(tmpDir, "src")
	dstDir := path.Join(tmpDir, "dst")
	_, err := InitGitRepository(srcDir, false)
	assert.Nil(t, err)
	createFile(t, srcDir, "one.txt")

	err = CopyDirectory(srcDir, dstDir)

	assert.Nil(t, err)
	assertFileExists(t, dstDir, "one.txt")
	assertFileExistsWithContent(t, dstDir, ".git/HEAD", "ref: refs/heads/master\n")
}
