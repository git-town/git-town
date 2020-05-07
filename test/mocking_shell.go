package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/git-town/git-town/src/command"
	"github.com/kballard/go-shellquote"
)

// MockingShell runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
// - overide environment variables
// - Temporarily override certain shell commands with mock implementations.
//   Temporary mocks are only valid for the next command being run.
type MockingShell struct {
	workingDir            string // the directory in which this runner runs
	homeDir               string // the directory that contains the global Git configuration
	tempShellOverridesDir string // the directory that stores the mock shell command implementations, ignored if empty
}

// NewMockingShell provides a new MockingShell instance that executes in the given directory.
func NewMockingShell(workingDir string, homeDir string) *MockingShell {
	return &MockingShell{workingDir: workingDir, homeDir: homeDir}
}

// AddTempShellOverride temporarily mocks the shell command with the given name
// with the given Bash script.
func (ms *MockingShell) AddTempShellOverride(name, content string) error {
	if !ms.hasTempShellOverrides() {
		err := ms.createTempShellOverridesDir()
		if err != nil {
			return fmt.Errorf("cannot create temp shell overrides dir: %w", err)
		}
	}
	return ioutil.WriteFile(ms.tempShellOverrideFilePath(name), []byte(content), 0744)
}

// createTempShellOverridesDir creates the folder that will contain the temp shell overrides.
// It is safe to call this method multiple times.
func (ms *MockingShell) createTempShellOverridesDir() error {
	var err error
	ms.tempShellOverridesDir, err = ioutil.TempDir("", "")
	return err
}

// hasTempShellOverrides indicates whether there are temp shell overrides for the next command.
func (ms *MockingShell) hasTempShellOverrides() bool {
	return ms.tempShellOverridesDir != ""
}

// MustRun runs the given command and returns the result. Panics on error.
func (ms *MockingShell) MustRun(cmd string, args ...string) *command.Result {
	res, err := ms.RunWith(command.Options{Essential: true}, cmd, args...)
	if err != nil {
		panic(err)
	}
	return res
}

// removeTempShellOverrides removes all custom shell overrides.
func (ms *MockingShell) removeTempShellOverrides() {
	os.RemoveAll(ms.tempShellOverridesDir)
	ms.tempShellOverridesDir = ""
}

// Run runs the given command with the given arguments
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) Run(name string, arguments ...string) (*command.Result, error) {
	return ms.RunWith(command.Options{}, name, arguments...)
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Shell overrides apply for the first command only.
// Failed commands abort immediately with the encountered error.
func (ms *MockingShell) RunMany(commands [][]string) error {
	for _, argv := range commands {
		command, args := argv[0], argv[1:]
		outcome, err := ms.Run(command, args...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w\n%v", argv, err, outcome)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments)
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) RunString(fullCmd string) (*command.Result, error) {
	return ms.RunStringWith(fullCmd, command.Options{})
}

// RunStringWith runs the given command (including possible arguments)
// in this ShellRunner's directory using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) RunStringWith(fullCmd string, opts command.Options) (result *command.Result, err error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return result, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return ms.RunWith(opts, cmd, args...)
}

// RunWith runs the given command with the given options in this ShellRunner's directory.
func (ms *MockingShell) RunWith(opts command.Options, cmd string, args ...string) (result *command.Result, err error) {
	// create an environment with the temp shell overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	for i := range opts.Env {
		if strings.HasPrefix(opts.Env[i], "HOME=") {
			opts.Env[i] = fmt.Sprintf("HOME=%s", ms.homeDir)
		}
	}
	// enable shell overrides
	if ms.hasTempShellOverrides() {
		for i := range opts.Env {
			if strings.HasPrefix(opts.Env[i], "PATH=") {
				parts := strings.SplitN(opts.Env[i], "=", 2)
				parts[1] = ms.tempShellOverridesDir + ":" + parts[1]
				opts.Env[i] = strings.Join(parts, "=")
				break
			}
		}
		defer ms.removeTempShellOverrides()
	}
	// set the working dir
	opts.Dir = filepath.Join(ms.workingDir, opts.Dir)
	// run the command inside the custom environment
	result, err = command.RunWith(opts, cmd, args...)
	if Debug {
		fmt.Println(filepath.Base(ms.workingDir), ">", cmd, strings.Join(args, " "))
		fmt.Println(result.Output())
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
	return result, err
}

// tempShellOverrideFilePath provides the full file path where to store a temp shell command with the given name.
func (ms *MockingShell) tempShellOverrideFilePath(shellOverrideFilename string) string {
	return filepath.Join(ms.tempShellOverridesDir, shellOverrideFilename)
}
