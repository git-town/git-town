//go:build windows
// +build windows

package ostools

import (
	"os"

	"github.com/git-town/git-town/v11/test/asserts"
)

// This package contains platform-specific testing tool implementations for the Windows platform.

// CallScriptArgs provides the command and arguments to call the given script on Windows.
func CallScriptArgs(toolPath string) (cmd string, args []string) {
	return "cmd.exe", []string{"/C", toolPath}
}

// CreateInputTool creates a tool that reads two inputs from STDIN and prints them back to the user.
func CreateInputTool(toolPath string) {
	asserts.NoError(os.WriteFile(ScriptName(toolPath), []byte(`
set /p i1=""
set /p i2=""
echo You entered %i1% and %i2%
`), 0o744)) //nolint:gosec
}

// CreateLsTool creates a tool in the given folder that lists all files in its current folder.
func CreateLsTool(toolPath string) {
	asserts.NoError(os.WriteFile(ScriptName(toolPath), []byte("@dir /B"), 0o744)) //nolint:gosec
}

// ScriptName provides the name of the given script file on the Windows.
func ScriptName(command string) string {
	return command + ".cmd"
}
