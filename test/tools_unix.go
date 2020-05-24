// +build !windows

package test

import "io/ioutil"

// This package contains platform-specific testing tool implementations for Unix-like platforms.

// ScriptName provides the name of the given script file on the Windows.
func ScriptName(command string) string {
	return command
}

// CreateLsTool creates a tool in the given folder that lists all files in its current folder.
func CreateLsTool(toolPath string) error {
	return ioutil.WriteFile(toolPath, []byte("#!/usr/bin/env bash\n\nls\n"), 0744)
}
