package run_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v7/src/run"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run(".InDir()", func(t *testing.T) {
		t.Parallel()
		dir, err := os.MkdirTemp("", "")
		assert.NoError(t, err)
		dirPath := filepath.Join(dir, "mydir")
		err = os.Mkdir(dirPath, 0o700)
		assert.NoError(t, err)
		err = os.WriteFile(filepath.Join(dirPath, "one"), []byte{}, 0o500)
		assert.NoError(t, err)
		res, err := run.InDir(dirPath, "ls", "-1")
		assert.NoError(t, err)
		assert.Equal(t, "one", res.OutputSanitized())
	})
}
