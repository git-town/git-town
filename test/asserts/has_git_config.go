package asserts

import (
	"os"
	"testing"

	"github.com/shoenig/test/must"
)

func HasGitConfiguration(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	must.NoError(t, err)
	for e := range entries {
		if entries[e].Name() == ".gitconfig" {
			return
		}
	}
	t.Fatalf(".gitconfig not found in %q", dir)
}
