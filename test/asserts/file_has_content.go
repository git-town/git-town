package asserts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

func FileHasContent(t *testing.T, dir, filename, expectedContent string) {
	t.Helper()
	fileContent, err := os.ReadFile(filepath.Join(dir, filename))
	must.NoError(t, err)
	must.EqOp(t, expectedContent, string(fileContent))
}
