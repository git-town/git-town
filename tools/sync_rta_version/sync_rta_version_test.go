package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/git-town/git-town/v23/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestReadCanonicalRTAVersionLine_readsFirstMatchingLine(t *testing.T) {
	t.Parallel()
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
}

func TestReplaceRTAVersionAssignment_replacesLfAndPreservesTrailingNewline(t *testing.T) {
	t.Parallel()
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
}

func TestReplaceRTAVersionAssignment_rewritesCrlf(t *testing.T) {
	t.Parallel()
	canonicalLine := `RTA_VERSION = 9.9.9`
	before := "RTA_VERSION = old\r\nFOO=1\r\n"
	after, modified := replaceRTAVersionAssignment(before, canonicalLine)
	want := canonicalLine + "\r\nFOO=1\r\n"
	if !modified || after != want {
		t.Fatalf("unexpected result modified=%v after=%q", modified, after)
	}
}

func TestReplaceRTAVersionAssignment_reportsNotModifiedWhenUnchanged(t *testing.T) {
	t.Parallel()
	canonicalLine := `RTA_VERSION = 1`
	original := canonicalLine + "\nNEXT=ok\n"
	after, modified := replaceRTAVersionAssignment(original, canonicalLine)
	if modified || after != original {
		t.Fatalf("want unchanged got modified=%v after=%q", modified, after)
	}
}

func TestReplaceRTAVersionAssignment_replacesIndentedAssignment(t *testing.T) {
	t.Parallel()
	canonicalLine := `RTA_VERSION = 2`
	before := "use tab\n\tRTA_VERSION = 1\nmore\n"
	after, modified := replaceRTAVersionAssignment(before, canonicalLine)
	want := "use tab\n" + canonicalLine + "\nmore\n"
	if !modified || after != want {
		t.Fatalf("unexpected modified=%v after=%q want %q", modified, after, want)
	}
}
