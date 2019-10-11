package helpers

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// CreateFile creates a file with the given name and content inside the given directory.
func CreateFile(t *testing.T, dir, filename, content string) {
	assert.Nilf(t, os.MkdirAll(dir, 0744), "cannot create directory %q", dir)
	ioutil.WriteFile(path.Join(dir, filename), []byte(content), 0744)
}
