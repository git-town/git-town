package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// MustRun executes an essential subshell command given in argv notation.
// Essential subshell commands are essential for the functioning of Git Town.
// If they fail, Git Town ends right there.
func MustRun(cmd string, args ...string) *Result {
	result := RunInDir("", cmd, args...)
	if result.Err() != nil {
		fmt.Printf("\n\nError running '%s %s': %s", cmd, strings.Join(args, " "), result.Err())
		os.Exit(1)
	}
	return result
}

// MustRunInDir executes an essential subshell command given in argv notation.
// Essential subshell commands are essential for the functioning of Git Town.
// If they fail, Git Town ends right there.
func MustRunInDir(dir string, cmd string, args ...string) *Result {
	result := RunInDir(dir, cmd, args...)
	if result.Err() != nil {
		fmt.Printf("\n\nError running '%s %s' in %s: %s", cmd, strings.Join(args, " "), dir, result.Err())
		os.Exit(1)
	}
	return result
}

// Run executes the command given in argv notation.
func Run(cmd string, args ...string) *Result {
	return RunInDir("", cmd, args...)
}

// RunInDir executes the given command in the given directory.
func RunInDir(dir string, cmd string, args ...string) *Result {
	return RunDirEnv(dir, os.Environ(), cmd, args...)
}

// RunDirEnv executes the given command in the given directory, using the given environment variables.
func RunDirEnv(dir string, env []string, cmd string, args ...string) *Result {
	logRun(cmd, args...)
	subProcess := exec.Command(cmd, args...) // #nosec
	if dir != "" {
		subProcess.Dir = dir
	}
	subProcess.Env = env
	output, err := subProcess.CombinedOutput()
	return &Result{
		command: cmd,
		args:    args,
		err:     err,
		output:  string(output),
	}
}
