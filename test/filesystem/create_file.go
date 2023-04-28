package filesystem

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createFile creates a file with the given filename in the given directory.
func CreateFile(t *testing.T, dir, filename string) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.MkdirAll(filepath.Dir(filePath), 0o744)
	assert.NoError(t, err)
	err = os.WriteFile(filePath, []byte(filename+" content"), 0o500)
	assert.NoError(t, err)
}
