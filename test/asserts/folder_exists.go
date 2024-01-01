package asserts

import (
	"os"
	"testing"

	"github.com/shoenig/test/must"
)

func FolderExists(t *testing.T, dir string) {
	t.Helper()
	_, err := os.Stat(dir)
	must.False(t, os.IsNotExist(err))
}
