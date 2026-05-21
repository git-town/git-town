package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v23/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestReadCanonicalRTAVersionLine(t *testing.T) {
	t.Parallel()

	t.Run("returns first declaration", func(t *testing.T) {
		tempDir := t.TempDir()
		mainMakefilePath := filepath.Join(tempDir, "Makefile")
		makefileContents := `
# header
RTA_VERSION = 1.2.3  # run-that-app version to use
OTHER = ok
RTA_VERSION = should not matter

`[1:]
		asserts.NoError(os.WriteFile(mainMakefilePath, []byte(makefileContents), 0o600))
		line, hasLine := readCanonicalRTAVersionLine(mainMakefilePath).Get()
		must.True(t, hasLine)
		must.EqOp(t, "RTA_VERSION = 1.2.3  # run-that-app version to use", line)
	})
}

func TestReplaceRTAVersionAssignment(t *testing.T) {
	t.Parallel()

	t.Run("replace declaration", func(t *testing.T) {
		give := `
HEADER = ok
RTA_VERSION = 9.9.9
TAIL=x
`[1:]
		have, modified := replaceRTAVersionAssignment(give, "RTA_VERSION = 1.2.3")
		want := `
HEADER = ok
RTA_VERSION = 1.2.3
TAIL=x
`[1:]
		must.True(t, modified)
		must.EqOp(t, want, have)
	})

	t.Run("unchanged", func(t *testing.T) {
		give := `
HEADER = ok
RTA_VERSION = 1.2.3
TAIL=x
`[1:]
		have, modified := replaceRTAVersionAssignment(give, "RTA_VERSION = 1.2.3")
		want := `
HEADER = ok
RTA_VERSION = 1.2.3
TAIL=x
`[1:]
		must.False(t, modified)
		must.EqOp(t, want, have)
	})
}
