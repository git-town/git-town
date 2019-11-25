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
	result, err := RunInDir("", cmd, args...)
	if err != nil {
		fmt.Printf("\n\nError running '%s %s': %s", cmd, strings.Join(args, " "), err)
		os.Exit(1)
	}
	return result
}

// Run executes the command given in argv notation.
func Run(cmd string, args ...string) (*Result, error) {
	return RunInDir("", cmd, args...)
}

// RunInDir executes the given command in the given directory.
func RunInDir(dir string, cmd string, args ...string) (*Result, error) {
	return RunDirEnv(dir, os.Environ(), cmd, args...)
}

// RunDirEnv executes the given command in the given directory, using the given environment variables.
func RunDirEnv(dir string, env []string, cmd string, args ...string) (*Result, error) {
	logRun(cmd, args...)
	subProcess := exec.Command(cmd, args...) // #nosec
	if dir != "" {
		subProcess.Dir = dir
	}
	subProcess.Env = env
	output, err := subProcess.CombinedOutput()
	result := &Result{
		cmd:    cmd,
		args:   args,
		output: string(output),
	}
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok {
			result.exitCode = exitError.ExitCode()
			return result, nil
		}
	}
	return result, err
}
