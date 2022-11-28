//go:build windows
// +build windows

package test

import "os"

// This package contains platform-specific testing tool implementations for the Windows platform.

// CallScriptArgs provides the command and arguments to call the given script on Windows.
func CallScriptArgs(toolPath string) (cmd string, args []string) {
	return "cmd.exe", []string{"/C", toolPath}
}

// CreateInputTool creates a tool that reads two inputs from STDIN and prints them back to the user.
func CreateInputTool(toolPath string) error {
	return os.WriteFile(ScriptName(toolPath), []byte(`
set /p i1=""
set /p i2=""
echo You entered %i1% and %i2%
`), 0o744)
}

// CreateLsTool creates a tool in the given folder that lists all files in its current folder.
func CreateLsTool(toolPath string) error {
	return os.WriteFile(ScriptName(toolPath), []byte("@dir /B"), 0o744)
}

// ScriptName provides the name of the given script file on the Windows.
func ScriptName(command string) string {
	return command + ".cmd"
}
