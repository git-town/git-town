package test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// Runner runs shell commands in a customizable environment.
// Possible customizations:
// - override certain shell commands permanently or temporary
type Runner struct {
	baseDir string
}

func NewRunner(baseDir string) *Runner {
	return &Runner{baseDir: baseDir}
}

// AddTempShellOverride adds a temporary mock of a shell command
// with the given name and content.
func (r *Runner) AddTempShellOverride(name, content string) error {
	err := r.createTempShellOverrideDir()
	if err != nil {
		return errors.Wrap(err, "cannot create temp shell overrides directory")
	}
	filePath := path.Join(r.tempShellOverrideDirPath(), name)
	err = ioutil.WriteFile(filePath, []byte(content), 0744)
	if err != nil {
		return errors.Wrapf(err, "cannot create temp shell override file %s", filePath)
	}
	return nil
}

func (r *Runner) RemoveTempShellOverrides() {
	os.RemoveAll(r.tempShellOverrideDirPath())
}

func (r *Runner) createTempShellOverrideDir() error {
	return os.MkdirAll(r.tempShellOverrideDirPath(), 0777)
}

func (r *Runner) tempShellOverrideDirPath() string {
	return path.Join(r.baseDir, "temp_shell_overrides")
}

// Run runs the given command with the given argv-like arguments in the current directory
// and stores the output and error for later analysis.
func (r *Runner) Run(name string, arguments ...string) (string, error) {

	// create an environment with the temp shell overrides directory added to the PATH
	customEnv := os.Environ()
	for i, entry := range customEnv {
		if strings.HasPrefix(entry, "PATH=") {
			parts := strings.SplitN(entry, "=", 2)
			parts[1] = r.tempShellOverrideDirPath() + ":" + parts[1]
			customEnv[i] = strings.Join(parts, "=")
			break
		}
	}

	// run the command inside the custom environment
	cmd := exec.Command(name, arguments...)
	cmd.Env = customEnv
	rawOutput, err := cmd.CombinedOutput()
	return string(rawOutput), err
}

// RunString runs the given command (that can contain arguments) in the current directory
// and stores the output and error for later analysis.
//
// Currently this splits the string by space,
// this only works for simple commands without quotes.
func (r *Runner) RunString(command string) (string, error) {
	parts := strings.Fields(command)
	command, args := parts[0], parts[1:]
	return r.Run(command, args...)
}

// RunMany runs all given commands in current directory.
// Failed commands cause abortion of the function with the received error.
func (r *Runner) RunMany(commands [][]string) error {
	for _, commandList := range commands {
		command, args := commandList[0], commandList[1:]
		_, err := r.Run(command, args...)
		if err != nil {
			return err
		}
	}
	return nil
}
