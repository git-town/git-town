package asserts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FileHasContent(t *testing.T, dir, filename, expectedContent string) {
	t.Helper()
	fileContent, err := os.ReadFile(filepath.Join(dir, filename))
	assert.NoError(t, err)
	assert.Equal(t, expectedContent, string(fileContent))
}
