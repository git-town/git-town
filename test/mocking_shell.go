package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/git-town/git-town/src/run"
	"github.com/kballard/go-shellquote"
)

// MockingShell runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
// - overide environment variables
// - Temporarily override certain shell commands with mock implementations.
//   Temporary mocks are only valid for the next command being run.
type MockingShell struct {
	workingDir     string // the directory in which this runner runs
	homeDir        string // the directory that contains the global Git configuration
	binDir         string // the directory that stores the mock shell command implementations, ignored if empty
	testOrigin     string // optional content of the GIT_TOWN_REMOTE environment variable
	hasMockCommand bool   // indicates whether the current test has mocked a command
}

// NewMockingShell provides a new MockingShell instance that executes in the given directory.
func NewMockingShell(workingDir string, homeDir string, binDir string) *MockingShell {
	return &MockingShell{workingDir: workingDir, homeDir: homeDir, binDir: binDir}
}

// WorkingDir provides the directory this MockingShell operates in.
func (ms *MockingShell) WorkingDir() string {
	return ms.workingDir
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (ms *MockingShell) MockBrokenCommand(name string) error {
	// create "bin" dir
	err := os.Mkdir(ms.binDir, 0744)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir: %w", err)
	}
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(ms.binDir, name))
	err = ioutil.WriteFile(filepath.Join(ms.binDir, "which"), []byte(content), 0500)
	if err != nil {
		return fmt.Errorf("cannot write custom which command: %w", err)
	}
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	err = ioutil.WriteFile(filepath.Join(ms.binDir, name), []byte(content), 0500)
	if err != nil {
		return fmt.Errorf("cannot write custom command: %w", err)
	}
	ms.hasMockCommand = true
	return nil
}

// MockCommand adds a mock for the command with the given name.
func (ms *MockingShell) MockCommand(name string) error {
	// create "bin" dir
	err := os.Mkdir(ms.binDir, 0744)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir: %w", err)
	}
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(ms.binDir, name))
	err = ioutil.WriteFile(filepath.Join(ms.binDir, "which"), []byte(content), 0500)
	if err != nil {
		return fmt.Errorf("cannot write custom which command: %w", err)
	}
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	err = ioutil.WriteFile(filepath.Join(ms.binDir, name), []byte(content), 0500)
	if err != nil {
		return fmt.Errorf("cannot write custom command: %w", err)
	}
	ms.hasMockCommand = true
	return nil
}

// MockGit pretends that this repo has Git in the given version installed.
func (ms *MockingShell) MockGit(version string) error {
	// create "bin" dir
	err := os.Mkdir(ms.binDir, 0744)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir %q: %w", ms.binDir, err)
	}
	// write custom Git command
	if runtime.GOOS == "windows" {
		content := fmt.Sprintf("echo git version %s\n", version)
		err = ioutil.WriteFile(filepath.Join(ms.binDir, "git.cmd"), []byte(content), 0500)
	} else {
		content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" = \"version\" ]; then\n  echo git version %s\nfi\n", version)
		err = ioutil.WriteFile(filepath.Join(ms.binDir, "git"), []byte(content), 0500)
	}
	if err != nil {
		return fmt.Errorf("cannot create custom Git binary: %w", err)
	}
	ms.hasMockCommand = true
	return nil
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (ms *MockingShell) MockNoCommandsInstalled() error {
	// create "bin" dir
	err := os.Mkdir(ms.binDir, 0744)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir: %w", err)
	}
	// write custom "which" command
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	err = ioutil.WriteFile(filepath.Join(ms.binDir, "which"), []byte(content), 0500)
	if err != nil {
		return fmt.Errorf("cannot write custom which command: %w", err)
	}
	ms.hasMockCommand = true
	return nil
}

// Run runs the given command with the given arguments
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) Run(name string, arguments ...string) (*run.Result, error) {
	return ms.RunWith(run.Options{}, name, arguments...)
}

// RunMany runs all given commands in current directory.
// Commands are provided as a list of argv-style strings.
// Shell overrides apply for the first command only.
// Failed commands abort immediately with the encountered error.
func (ms *MockingShell) RunMany(commands [][]string) error {
	for _, argv := range commands {
		command, args := argv[0], argv[1:]
		_, err := ms.Run(command, args...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments)
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) RunString(fullCmd string) (*run.Result, error) {
	return ms.RunStringWith(fullCmd, run.Options{})
}

// RunStringWith runs the given command (including possible arguments)
// in this ShellRunner's directory using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) RunStringWith(fullCmd string, opts run.Options) (result *run.Result, err error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return result, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return ms.RunWith(opts, cmd, args...)
}

// RunWith runs the given command with the given options in this ShellRunner's directory.
func (ms *MockingShell) RunWith(opts run.Options, cmd string, args ...string) (result *run.Result, err error) {
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
	// add the custom origin
	if ms.testOrigin != "" {
		opts.Env = append(opts.Env, fmt.Sprintf("GIT_TOWN_REMOTE=%s", ms.testOrigin))
	}
	// add the custom bin dir to the PATH
	if ms.hasMockCommand {
		for i := range opts.Env {
			if strings.HasPrefix(opts.Env[i], "PATH=") {
				parts := strings.SplitN(opts.Env[i], "=", 2)
				parts[1] = ms.binDir + string(os.PathListSeparator) + parts[1]
				opts.Env[i] = strings.Join(parts, "=")
				break
			}
		}
	}
	// set the working dir
	opts.Dir = filepath.Join(ms.workingDir, opts.Dir)
	// run the command inside the custom environment
	result, err = run.WithOptions(opts, cmd, args...)
	if Debug {
		fmt.Println(filepath.Base(ms.workingDir), ">", cmd, strings.Join(args, " "))
		fmt.Println(result.Output())
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
	return result, err
}

// SetTestOrigin adds the given environment variable to subsequent runs of commands.
func (ms *MockingShell) SetTestOrigin(content string) {
	ms.testOrigin = content
}
