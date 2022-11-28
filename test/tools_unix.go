//go:build !windows
// +build !windows

package test

import "os"

// This package contains platform-specific testing tool implementations for Unix-like platforms.

// CallScriptArgs provides the command and arguments to call the given script on Windows.
//
//nolint:nonamedreturns
func CallScriptArgs(toolPath string) (cmd string, args []string) {
	return toolPath, []string{}
}

// CreateInputTool creates a tool that reads two inputs from STDIN and prints them back to the user.
func CreateInputTool(toolPath string) error {
	return os.WriteFile(toolPath, []byte(`#!/usr/bin/env bash
read i1
read i2
echo You entered $i1 and $i2
`), 0o500)
}

// CreateLsTool creates a tool in the given folder that lists all files in its current folder.
func CreateLsTool(toolPath string) error {
	return os.WriteFile(toolPath, []byte("#!/usr/bin/env bash\n\nls\n"), 0o500)
}

// ScriptName provides the name of the given script file on the Windows.
func ScriptName(command string) string {
	return command
}
