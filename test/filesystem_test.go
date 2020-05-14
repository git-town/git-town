package test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyDirectory(t *testing.T) {
	tmpDir := createTempDir(t)
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	createFile(t, srcDir, "one.txt")
	createFile(t, srcDir, "f1/a.txt")
	createFile(t, srcDir, "f2/b.txt")
	err := CopyDirectory(srcDir, dstDir)
	assert.Nil(t, err)
	assertFileExists(t, dstDir, "one.txt")
	assertFileExists(t, dstDir, "f1/a.txt")
	assertFileExists(t, dstDir, "f2/b.txt")
}

func TestCopyDirectory_GitRepo(t *testing.T) {
	tmpDir := createTempDir(t)
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")
	_, err := InitGitRepository(srcDir, tmpDir, "")
	assert.Nil(t, err)
	createFile(t, srcDir, "one.txt")
	err = CopyDirectory(srcDir, dstDir)
	assert.Nil(t, err)
	assertFileExists(t, dstDir, "one.txt")
	assertFileExistsWithContent(t, dstDir, ".git/HEAD", "ref: refs/heads/master\n")
}
