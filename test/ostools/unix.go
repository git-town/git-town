//go:build !windows
// +build !windows

package ostools

import (
	"os"

	"github.com/git-town/git-town/v13/test/asserts"
)

// This package contains platform-specific testing tool implementations for Unix-like platforms.

// CreateLsTool creates a tool in the given folder that lists all files in its current folder.
func CreateLsTool(toolPath string) {
	//nolint:gosec // intentionally creating an executable here
	asserts.NoError(os.WriteFile(toolPath, []byte("#!/usr/bin/env bash\n\nls\n"), 0o744))
}

// ScriptName provides the name of the given script file on the Windows.
func ScriptName(command string) string {
	return command
}
