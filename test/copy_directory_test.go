package test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Originate/exit"
	"github.com/stretchr/testify/assert"
)

func TestCopyDirectory(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	srcDir := path.Join(tmpDir, "src")
	dstDir := path.Join(tmpDir, "dst")

	// create a few files and folders
	createFile(srcDir, "one.txt")
	createFile(srcDir, "f1/a.txt")
	createFile(srcDir, "f2/b.txt")

	// copy them
	err = CopyDirectory(srcDir, dstDir)
	assert.Nil(t, err)

	// verify that the destination exists
	assertFileExists(t, dstDir, "one.txt")
	assertFileExists(t, dstDir, "f1/a.txt")
	assertFileExists(t, dstDir, "f2/b.txt")
}

func TestCopyGitRepo(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	srcDir := path.Join(tmpDir, "src")
	dstDir := path.Join(tmpDir, "dst")

	// create a few files and folders
	InitGitRepository(srcDir, false)
	createFile(srcDir, "one.txt")

	// copy them
	err = CopyDirectory(srcDir, dstDir)
	assert.Nil(t, err)

	// verify that the destination exists
	assertFileExists(t, dstDir, "one.txt")
	assertFileExistsWithContent(t, dstDir, ".git/HEAD", "ref: refs/heads/master\n")
}

func createFile(dir, filename string) {
	filePath := path.Join(dir, filename)
	err := os.MkdirAll(path.Dir(filePath), 0744)
	exit.If(err)
	err = ioutil.WriteFile(filePath, []byte(filename+" content"), 0644)
	exit.If(err)
}
