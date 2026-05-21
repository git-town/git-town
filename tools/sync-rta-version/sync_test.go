package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadCanonicalRTAVersionLine_readsFirstMatchingLine(t *testing.T) {
	t.Parallel()
	canonicalWant := `RTA_VERSION = 1.2.3  # run-that-app version to use`
	tempDir := t.TempDir()
	mainMakefilePath := filepath.Join(tempDir, "Makefile")
	makefileContents := strings.Join([]string{
		"# header",
		canonicalWant,
		"OTHER = ok",
		"RTA_VERSION = should not matter",
	}, "\n") + "\n"
	writeErr := os.WriteFile(mainMakefilePath, []byte(makefileContents), 0o600)
	if writeErr != nil {
		t.Fatalf("write makefile: %v", writeErr)
	}
	line, parseErr := readCanonicalRTAVersionLine(mainMakefilePath)
	if parseErr != nil {
		t.Fatalf("read: %v", parseErr)
	}
	if line != canonicalWant {
		t.Fatalf("got %q want %q", line, canonicalWant)
	}
}

func TestReplaceRTAVersionAssignment_replacesLfAndPreservesTrailingNewline(t *testing.T) {
	t.Parallel()
	canonicalLine := `RTA_VERSION = 9.9.9  # run-that-app version to use`
	before := "HEADER = ok\nRTA_VERSION = 0.nope  # run-that-app version to use\nTAIL=x\n"
	after, modified := replaceRTAVersionAssignment(before, canonicalLine)
	want := "HEADER = ok\n" + canonicalLine + "\nTAIL=x\n"
	if !modified || after != want {
		t.Fatalf("unexpected result modified=%v after=%q", modified, after)
	}
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
