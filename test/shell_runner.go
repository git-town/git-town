package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Originate/git-town/src/command"
	"github.com/kballard/go-shellquote"
)

// ShellRunner runs shell commands in the given directory, using a customizable environment.
// Possible customizations:
// - Temporarily override certain shell commands with mock implementations.
//   Temporary mocks are only valid for the next command being run.
type ShellRunner struct {

	// workingDir contains the directory path in which this runner runs.
	workingDir string

	// homeDir contains the path that contains the global Git configuration.
	homeDir string

	// tempShellOverrideDirDir contains the directory path that stores the mock shell command implementations.
	// This variable is populated when shell overrides are being set.
	// An empty string indicates that no shell overrides have been set.
	tempShellOverridesDir string
}

// NewShellRunner provides a new ShellRunner instance that executes in the given directory.
func NewShellRunner(workingDir string, homeDir string) ShellRunner {
	return ShellRunner{workingDir: workingDir, homeDir: homeDir}
}

// AddTempShellOverride temporarily mocks the shell command with the given name
// with the given Bash script.
func (runner *ShellRunner) AddTempShellOverride(name, content string) error {
	if !runner.hasTempShellOverrides() {
		err := runner.createTempShellOverridesDir()
		if err != nil {
			return fmt.Errorf("cannot create temp shell overrides dir: %w", err)
		}
	}
	return ioutil.WriteFile(runner.tempShellOverrideFilePath(name), []byte(content), 0744)
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

// RemoveTempShellOverrides removes all custom shell overrides.
func (runner *ShellRunner) RemoveTempShellOverrides() {
	os.RemoveAll(runner.tempShellOverridesDir)
	runner.tempShellOverridesDir = ""
}

// Run runs the given command with the given arguments
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (runner *ShellRunner) Run(name string, arguments ...string) (*command.Result, error) {
	return runner.RunWith(command.Options{}, name, arguments...)
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Shell overrides apply for the first command only.
// Failed commands abort immediately with the encountered error.
func (runner *ShellRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		command, args := argv[0], argv[1:]
		outcome, err := runner.Run(command, args...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w\n%v", argv, err, outcome)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments)
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (runner *ShellRunner) RunString(fullCmd string) (*command.Result, error) {
	return runner.RunStringWith(fullCmd, command.Options{})
}

// RunStringWith runs the given command (including possible arguments)
// in this ShellRunner's directory using the given options.
// Shell overrides will be used and removed when done.
func (runner *ShellRunner) RunStringWith(fullCmd string, opts command.Options) (result *command.Result, err error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return result, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return runner.RunWith(opts, cmd, args...)
}

// RunWith runs the given command with the given options in this ShellRunner's directory.
func (runner *ShellRunner) RunWith(opts command.Options, cmd string, args ...string) (result *command.Result, err error) {
	// create an environment with the temp shell overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	for i := range opts.Env {
		if strings.HasPrefix(opts.Env[i], "HOME=") {
			opts.Env[i] = fmt.Sprintf("HOME=%s", runner.homeDir)
		}
	}
	// enable shell overrides
	if runner.hasTempShellOverrides() {
		for i := range opts.Env {
			if strings.HasPrefix(opts.Env[i], "PATH=") {
				parts := strings.SplitN(opts.Env[i], "=", 2)
				parts[1] = runner.tempShellOverridesDir + ":" + parts[1]
				opts.Env[i] = strings.Join(parts, "=")
				break
			}
		}
		defer runner.RemoveTempShellOverrides()
	}
	// set the working dir
	opts.Dir = filepath.Join(runner.workingDir, opts.Dir)

	// run the command inside the custom environment
	result, err = command.RunWith(opts, cmd, args...)
	if Debug {
		fmt.Println(filepath.Base(runner.workingDir), ">", cmd, strings.Join(args, " "))
		fmt.Println(result.Output())
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
	return result, err
}

// tempShellOverrideFilePath provides the full file path where to store a temp shell command with the given name.
func (runner *ShellRunner) tempShellOverrideFilePath(shellOverrideFilename string) string {
	return filepath.Join(runner.tempShellOverridesDir, shellOverrideFilename)
}
