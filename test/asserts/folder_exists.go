package asserts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FolderExists(t *testing.T, dir string) {
	t.Helper()
	_, err := os.Stat(dir)
	assert.Falsef(t, os.IsNotExist(err), "directory %q not found", dir)
}
