package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

// ShellRunner runs shell commands in the given directory, using a customizable environment.
// Possible customizations:
// - Temporarily override certain shell commands with mock implementations.
//   Temporary mocks are only valid for the next command being run.
type ShellRunner struct {

	// dir contains the directory path in which this runner runs.
	dir string

	// tempShellOverrideDirDir contains the directory path that stores the mock shell command implementations.
	// This variable is populated when shell overrides are being set.
	// An empty string indicates that no shell overrides have been set.
	tempShellOverridesDir string
}

// NewShellRunner provides a new ShellRunner instance that executes in the given directory.
func NewShellRunner(dir string) ShellRunner {
	return ShellRunner{dir: dir}
}

// AddTempShellOverride temporarily mocks the shell command with the given name
// with the given Bash script.
func (runner *ShellRunner) AddTempShellOverride(name, content string) error {
	if !runner.hasTempShellOverrides() {
		err := runner.createTempShellOverridesDir()
		if err != nil {
			return errors.Wrap(err, "cannot create temp shell overrides dir")
		}
	}
	return ioutil.WriteFile(runner.tempShellOverrideFilePath(name), []byte(content), 0744)
}

// tempShellOverrideFilePath provides the full file path where to store a temp shell command with the given name.
func (runner *ShellRunner) tempShellOverrideFilePath(shellOverrideFilename string) string {
	return path.Join(runner.tempShellOverridesDir, shellOverrideFilename)
}

// RemoveTempShellOverrides removes all custom shell overrides.
func (runner *ShellRunner) RemoveTempShellOverrides() {
	os.RemoveAll(runner.tempShellOverridesDir)
	runner.tempShellOverridesDir = ""
}

// createTempShellOverridesDir creates the folder that will contain the temp shell overrides.
// It is safe to call this method multiple times.
func (runner *ShellRunner) createTempShellOverridesDir() error {
	var err error
	runner.tempShellOverridesDir, err = ioutil.TempDir("", "")
	return err
}

// hasTempShellOverrides indicates whether there are temp shell overrides for the next command.
func (runner *ShellRunner) hasTempShellOverrides() bool {
	return runner.tempShellOverridesDir != ""
}

// Run runs the given command with the given argv-like arguments
// in this ShellRunner's directory.
// Shell overrides will be used and removed.
func (runner *ShellRunner) Run(name string, arguments ...string) (output string, err error) {
	// create an environment with the temp shell overrides directory added to the PATH
	customEnv := os.Environ()
	if runner.hasTempShellOverrides() {
		for i := range customEnv {
			if strings.HasPrefix(customEnv[i], "PATH=") {
				parts := strings.SplitN(customEnv[i], "=", 2)
				parts[1] = runner.tempShellOverridesDir + ":" + parts[1]
				customEnv[i] = strings.Join(parts, "=")
				break
			}
		}
		defer runner.RemoveTempShellOverrides()
	}

	// run the command inside the custom environment
	cmd := exec.Command(name, arguments...)
	cmd.Dir = runner.dir
	cmd.Env = customEnv
	rawOutput, err := cmd.CombinedOutput()
	return string(rawOutput), err
}

// RunString runs the given command (including possible arguments)
// in this ShellRunner's directory.
// Shell overrides will be used and removed.
//
// The current implementation splits the string by space
// and therefore only works for simple commands without quoted arguments.
func (runner *ShellRunner) RunString(command string) (output string, err error) {
	parts, err := shellquote.Split(command)
	if err != nil {
		return "", errors.Wrapf(err, "cannot split command: %q", command)
	}
	command, args := parts[0], parts[1:]
	return runner.Run(command, args...)
}

// RunMany runs all given commands in current directory.
// Failed commands cause abortion of the function with the received error.
func (runner *ShellRunner) RunMany(commands [][]string) error {
	for _, commandList := range commands {
		command, args := commandList[0], commandList[1:]
		_, err := runner.Run(command, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
