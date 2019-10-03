package test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createFile creates a file with the given filename in the given directory.
func createFile(t *testing.T, dir, filename string) {
	filePath := path.Join(dir, filename)
	err := os.MkdirAll(path.Dir(filePath), 0744)
	assert.Nil(t, err)
	err = ioutil.WriteFile(filePath, []byte(filename+" content"), 0644)
	assert.Nil(t, err)
}

// createTestDir provides the path to an empty directory in the system's temp directory
func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "")
	assert.Nil(t, err, "cannot create TempDir")
	return dir
}
