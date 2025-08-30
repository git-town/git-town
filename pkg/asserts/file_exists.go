package asserts

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

func FileExists(t *testing.T, dir, filename string) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	info, err := os.Stat(filePath)
	must.NoError(t, err)
	must.False(t, info.IsDir())
}
