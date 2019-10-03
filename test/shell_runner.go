package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// ShellRunner runs shell commands in a particular directory, using a customizable environment.
// Possible customizations:
// - override certain shell commands permanently or temporary
type ShellRunner struct {

	// dir contains the directory in which this runner runs.
	dir string

	// tempShellOverrideDirDir contains the path of the directory in which the temp shell overrides exist.
	// This variable is only populated when temp shell overrides are set.
	tempShellOverridesDir string
}

// NewShellRunner provides a new ShellRunner instance that runs in the given directory.
func NewShellRunner(dir string) ShellRunner {
	return ShellRunner{dir: dir}
}

// AddTempShellOverride adds a temporary mock of a shell command
// with the given name and file content.
func (runner *ShellRunner) AddTempShellOverride(name, content string) error {
	if !runner.hasTempShellOverride() {
		err := runner.createTempShellOverridesDir()
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(runner.tempShellOverrideFilePath(name), []byte(content), 0744)
}

// tempShellOverrideFilePath returns the path where to store the given
// temp shell override on disk.
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

// hasTempShellOverrideDir indicates whether a folder for the temp shell overrides was already created.
func (runner *ShellRunner) hasTempShellOverride() bool {
	return runner.tempShellOverridesDir != ""
}

// Run runs the given command with the given argv-like arguments in the current directory
// and stores the output and error for later analysis.
func (runner *ShellRunner) Run(name string, arguments ...string) (output string, err error) {

	// create an environment with the temp shell overrides directory added to the PATH
	customEnv := os.Environ()
	if runner.hasTempShellOverride() {
		for i, entry := range customEnv {
			if strings.HasPrefix(entry, "PATH=") {
				parts := strings.SplitN(entry, "=", 2)
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

// RunString runs the given command (that can contain arguments) in the current directory
// and stores the output and error for later analysis.
//
// Currently this splits the string by space,
// this only works for simple commands without quotes.
func (runner *ShellRunner) RunString(command string) (output string, err error) {
	parts := strings.Fields(command)
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
