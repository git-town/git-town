//go:build windows

package ostools

import (
	"os"

	"github.com/git-town/git-town/v22/pkg/asserts"
)

// This package contains OS-specific testing tool implementations for Windows.

// CreateLsTool creates a tool in the given folder that lists all files in its current folder.
func CreateLsTool(toolPath string) {
	asserts.NoError(os.WriteFile(ScriptName(toolPath), []byte("@dir /B"), 0o744)) //nolint:gosec
}

// ScriptName provides the name of the given script file on the Windows.
func ScriptName(command string) string {
	return command + ".cmd"
}
