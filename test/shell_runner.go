package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Originate/git-town/src/command"
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

	// globalDir contains the path that contains the global Git configuration.
	globalDir string

	// tempShellOverrideDirDir contains the directory path that stores the mock shell command implementations.
	// This variable is populated when shell overrides are being set.
	// An empty string indicates that no shell overrides have been set.
	tempShellOverridesDir string
}

// NewShellRunner provides a new ShellRunner instance that executes in the given directory.
func NewShellRunner(dir string, globalDir string) ShellRunner {
	return ShellRunner{dir: dir, globalDir: globalDir}
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

// Run runs the given command with the given arguments
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (runner *ShellRunner) Run(name string, arguments ...string) (output string, err error) {
	// create an environment with the temp shell overrides directory added to the PATH
	customEnv := os.Environ()

	// set HOME to the given global directory so that Git puts the global configuration there.
	for i := range customEnv {
		if strings.HasPrefix(customEnv[i], "HOME=") {
			customEnv[i] = fmt.Sprintf("HOME=%s", runner.globalDir)
		}
	}

	// enable shell overrides
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
	outcome := command.RunDirEnv(runner.dir, customEnv, name, arguments...)
	if Debug {
		fmt.Println(path.Base(runner.dir), ">", name, strings.Join(arguments, " "))
		fmt.Println(outcome.Output())
	}
	return outcome.Output(), outcome.Err()
}

// RunString runs the given command (including possible arguments)
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
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
// Commands are provided as a list of argv-style strings.
// Shell overrides apply for the first command only.
// Failed commands abort immediately with the encountered error.
func (runner *ShellRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		command, args := argv[0], argv[1:]
		output, err := runner.Run(command, args...)
		if err != nil {
			return errors.Wrapf(err, "error running command %q: %s", argv, output)
		}
	}
	return nil
}
