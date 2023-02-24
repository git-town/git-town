package test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/git-town/git-town/v7/src/envvars"
	"github.com/git-town/git-town/v7/src/run"
	"github.com/kballard/go-shellquote"
)

// MockingShell runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
//   - overide environment variables
//   - Temporarily override certain shell commands with mock implementations.
//     Temporary mocks are only valid for the next command being run.
type MockingShell struct {
	// the directory that stores the mock shell command implementations, ignored if empty
	binDir string

	// whether to log the output of subshell commands
	Debug bool `exhaustruct:"optional"`

	// name of the binary to use as the custom editor during "git commit"
	gitEditor string `exhaustruct:"optional"`

	// the directory that contains the global Git configuration
	homeDir string

	// optional content of the GIT_TOWN_REMOTE environment variable
	testOrigin string `exhaustruct:"optional"`

	// indicates whether the current test has created the binDir
	usesBinDir bool `exhaustruct:"optional"`

	// the directory in which this runner runs
	workingDir string
}

// NewMockingShell provides a new MockingShell instance that executes in the given directory.
func NewMockingShell(workingDir string, homeDir string, binDir string) MockingShell {
	return MockingShell{
		workingDir: workingDir,
		homeDir:    homeDir,
		binDir:     binDir,
	}
}

// createBinDir creates the directory that contains mock shell command implementations.
// This method is idempotent.
func (ms *MockingShell) createBinDir() error {
	if ms.usesBinDir {
		// binDir already created --> nothing to do here
		return nil
	}
	err := os.Mkdir(ms.binDir, 0o700)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir: %w", err)
	}
	ms.usesBinDir = true
	return nil
}

// createMockBinary creates an executable with the given name and content in ms.binDir.
func (ms *MockingShell) createMockBinary(name string, content string) error {
	if err := ms.createBinDir(); err != nil {
		return err
	}
	err := os.WriteFile(filepath.Join(ms.binDir, name), []byte(content), 0o500)
	if err != nil {
		return fmt.Errorf("cannot write custom %q command: %w", name, err)
	}
	return nil
}

// WorkingDir provides the directory this MockingShell operates in.
func (ms *MockingShell) WorkingDir() string {
	return ms.workingDir
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (ms *MockingShell) MockBrokenCommand(name string) error {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(ms.binDir, name))
	err := ms.createMockBinary("which", content)
	if err != nil {
		return err
	}
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	return ms.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (ms *MockingShell) MockCommand(name string) error {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(ms.binDir, name))
	if err := ms.createMockBinary("which", content); err != nil {
		return fmt.Errorf("cannot write custom which command: %w", err)
	}
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	return ms.createMockBinary(name, content)
}

// MockGit pretends that this repo has Git in the given version installed.
func (ms *MockingShell) MockGit(version string) error {
	if runtime.GOOS == "windows" {
		// create Windows binary
		content := fmt.Sprintf("echo git version %s\n", version)
		return ms.createMockBinary("git.cmd", content)
	}
	// create Unix binary
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" = \"version\" ]; then\n  echo git version %s\nfi\n", version)
	return ms.createMockBinary("git", content)
}

// MockCommitMessage sets up this shell with an editor that enters the given commit message.
func (ms *MockingShell) MockCommitMessage(message string) error {
	ms.gitEditor = "git_editor"
	return ms.createMockBinary(ms.gitEditor, fmt.Sprintf("#!/usr/bin/env bash\n\necho %q > $1", message))
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (ms *MockingShell) MockNoCommandsInstalled() error {
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	return ms.createMockBinary("which", content)
}

// Run runs the given command with the given arguments
// in this ShellRunner's directory.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) Run(name string, arguments ...string) (*run.Result, error) {
	return ms.RunWith(&run.Options{}, name, arguments...)
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
	return ms.RunStringWith(fullCmd, &run.Options{})
}

// RunStringWith runs the given command (including possible arguments)
// in this ShellRunner's directory using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Shell overrides will be used and removed when done.
func (ms *MockingShell) RunStringWith(fullCmd string, opts *run.Options) (*run.Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return ms.RunWith(opts, cmd, args...)
}

// RunWith runs the given command with the given options in this ShellRunner's directory.
func (ms *MockingShell) RunWith(opts *run.Options, cmd string, args ...string) (*run.Result, error) {
	// create an environment with the temp shell overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	opts.Env = envvars.Replace(opts.Env, "HOME", ms.homeDir)
	// add the custom origin
	if ms.testOrigin != "" {
		opts.Env = envvars.Replace(opts.Env, "GIT_TOWN_REMOTE", ms.testOrigin)
	}
	// add the custom bin dir to the PATH
	if ms.usesBinDir {
		opts.Env = envvars.PrependPath(opts.Env, ms.binDir)
	}
	// add the custom GIT_EDITOR
	if ms.gitEditor != "" {
		opts.Env = envvars.Replace(opts.Env, "GIT_EDITOR", filepath.Join(ms.binDir, "git_editor"))
	}
	// set the working dir
	opts.Dir = filepath.Join(ms.workingDir, opts.Dir)
	// run the command inside the custom environment
	result, err := run.WithOptions(opts, cmd, args...)
	if ms.Debug {
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
