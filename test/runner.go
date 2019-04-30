package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

// Runner runs shell commands in a customizable environment.
// Possible customizations:
// - override certain shell commands permanently or temporary
type Runner struct {

	// tempShellOverrideDirDir contains the path of the directory in which the temp shell overrides exist.
	// This variable is only populated when temp shell overrides are set.
	tempShellOverridesDir string
}

// RunResult represents the outcomes of a command that was run.
type RunResult struct {
	Output string
	Err    error
}

// AddTempShellOverride adds a temporary mock of a shell command
// with the given name and file content.
func (r *Runner) AddTempShellOverride(name, content string) error {
	if !r.hasTempShellOverride() {
		err := r.createTempShellOverridesDir()
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(r.tempShellOverrideFilePath(name), []byte(content), 0744)
}

// tempShellOverrideFilePath returns the path where to store the given
// temp shell override on disk.
func (r *Runner) tempShellOverrideFilePath(shellOverrideFilename string) string {
	return path.Join(r.tempShellOverridesDir, shellOverrideFilename)
}

func (r *Runner) RemoveTempShellOverrides() {
	os.RemoveAll(r.tempShellOverridesDir)
	r.tempShellOverridesDir = ""
}

// createTempShellOverridesDir creates the folder that will contain the temp shell overrides.
// It is safe to call this method multiple times.
func (r *Runner) createTempShellOverridesDir() error {
	var err error
	r.tempShellOverridesDir, err = ioutil.TempDir("", "")
	return err
}

// hasTempShellOverrideDir returns whether a folder for the temp shell overrides was already created.
func (r *Runner) hasTempShellOverride() bool {
	return r.tempShellOverridesDir != ""
}

// Run runs the given command with the given argv-like arguments in the current directory
// and stores the output and error for later analysis.
func (r *Runner) Run(name string, arguments ...string) RunResult {

	// create an environment with the temp shell overrides directory added to the PATH
	customEnv := os.Environ()
	if r.hasTempShellOverride() {
		for i, entry := range customEnv {
			if strings.HasPrefix(entry, "PATH=") {
				parts := strings.SplitN(entry, "=", 2)
				parts[1] = r.tempShellOverridesDir + ":" + parts[1]
				customEnv[i] = strings.Join(parts, "=")
				break
			}
		}
	}

	// run the command inside the custom environment
	cmd := exec.Command(name, arguments...)
	cmd.Env = customEnv
	rawOutput, err := cmd.CombinedOutput()

	r.RemoveTempShellOverrides()

	return RunResult{string(rawOutput), err}
}

func (r *Runner) hasTempShellOverrides() bool {
	return len(r.tempShellOverridesDir) > 0
}

// RunString runs the given command (that can contain arguments) in the current directory
// and stores the output and error for later analysis.
//
// Currently this splits the string by space,
// this only works for simple commands without quotes.
func (r *Runner) RunString(command string) RunResult {
	parts := strings.Fields(command)
	command, args := parts[0], parts[1:]
	return r.Run(command, args...)
}

// RunMany runs all given commands in current directory.
// Failed commands cause abortion of the function with the received error.
func (r *Runner) RunMany(commands [][]string) error {
	for _, commandList := range commands {
		command, args := commandList[0], commandList[1:]
		result := r.Run(command, args...)
		if result.Err != nil {
			return result.Err
		}
	}
	return nil
}
