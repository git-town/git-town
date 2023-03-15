package test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/kballard/go-shellquote"
)

// MockingRunner runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
//   - overide environment variables
//   - Temporarily override certain shell commands with mock implementations.
//     Temporary mocks are only valid for the next command being run.
type MockingRunner struct {
	// the directory that contains mock executables, ignored if empty
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

	// the directory in which this runner executes shell commands
	workingDir string
}

// NewMockingRunner provides a new MockingRunner instance that executes in the given directory.
func NewMockingRunner(workingDir string, homeDir string, binDir string) MockingRunner {
	return MockingRunner{
		workingDir: workingDir,
		homeDir:    homeDir,
		binDir:     binDir,
	}
}

// createBinDir creates the directory that contains mock executables.
// This method is idempotent.
func (r *MockingRunner) createBinDir() error {
	if r.usesBinDir {
		// binDir already created --> nothing to do here
		return nil
	}
	err := os.Mkdir(r.binDir, 0o700)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir: %w", err)
	}
	r.usesBinDir = true
	return nil
}

// createMockBinary creates an executable with the given name and content in ms.binDir.
func (r *MockingRunner) createMockBinary(name string, content string) error {
	if err := r.createBinDir(); err != nil {
		return err
	}
	err := os.WriteFile(filepath.Join(r.binDir, name), []byte(content), 0o500)
	if err != nil {
		return fmt.Errorf("cannot write custom %q command: %w", name, err)
	}
	return nil
}

// WorkingDir provides the directory this MockingRunner operates in.
func (r *MockingRunner) WorkingDir() string {
	return r.workingDir
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (r *MockingRunner) MockBrokenCommand(name string) error {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(r.binDir, name))
	err := r.createMockBinary("which", content)
	if err != nil {
		return err
	}
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	return r.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (r *MockingRunner) MockCommand(name string) error {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(r.binDir, name))
	if err := r.createMockBinary("which", content); err != nil {
		return fmt.Errorf("cannot write custom which command: %w", err)
	}
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	return r.createMockBinary(name, content)
}

// MockGit pretends that this repo has Git in the given version installed.
func (r *MockingRunner) MockGit(version string) error {
	if runtime.GOOS == "windows" {
		// create Windows binary
		content := fmt.Sprintf("echo git version %s\n", version)
		return r.createMockBinary("git.cmd", content)
	}
	// create Unix binary
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" = \"version\" ]; then\n  echo git version %s\nfi\n", version)
	return r.createMockBinary("git", content)
}

// MockCommitMessage sets up this runner with an editor that enters the given commit message.
func (r *MockingRunner) MockCommitMessage(message string) error {
	r.gitEditor = "git_editor"
	return r.createMockBinary(r.gitEditor, fmt.Sprintf("#!/usr/bin/env bash\n\necho %q > $1", message))
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (r *MockingRunner) MockNoCommandsInstalled() error {
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	return r.createMockBinary("which", content)
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (r *MockingRunner) Run(name string, arguments ...string) (*subshell.Result, error) {
	return r.RunWith(&subshell.Options{}, name, arguments...)
}

// RunMany runs all given commands.
// Commands are provided as a list of argv-style strings.
// Overrides apply for the first command only.
// Failed commands abort immediately with the encountered error.
func (r *MockingRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		command, args := argv[0], argv[1:]
		_, err := r.Run(command, args...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// RunString runs the given command (including possible arguments).
// Overrides will be used and removed when done.
func (r *MockingRunner) RunString(fullCmd string) (*subshell.Result, error) {
	return r.RunStringWith(fullCmd, &subshell.Options{})
}

// RunStringWith runs the given command (including possible arguments) using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Overrides will be used and removed when done.
func (r *MockingRunner) RunStringWith(fullCmd string, opts *subshell.Options) (*subshell.Result, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return nil, fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.RunWith(opts, cmd, args...)
}

// RunWith runs the given command with the given options in this ShellRunner's directory.
func (r *MockingRunner) RunWith(opts *subshell.Options, cmd string, args ...string) (*subshell.Result, error) {
	// create an environment with the temp Overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	opts.Env = ReplaceEnvVar(opts.Env, "HOME", r.homeDir)
	// add the custom origin
	if r.testOrigin != "" {
		opts.Env = ReplaceEnvVar(opts.Env, "GIT_TOWN_REMOTE", r.testOrigin)
	}
	// add the custom bin dir to the PATH
	if r.usesBinDir {
		opts.Env = PrependEnvPath(opts.Env, r.binDir)
	}
	// add the custom GIT_EDITOR
	if r.gitEditor != "" {
		opts.Env = ReplaceEnvVar(opts.Env, "GIT_EDITOR", filepath.Join(r.binDir, "git_editor"))
	}
	// set the working dir
	opts.Dir = filepath.Join(r.workingDir, opts.Dir)
	// run the command inside the custom environment
	result, err := subshell.WithOptions(opts, cmd, args...)
	if r.Debug {
		fmt.Println(filepath.Base(r.workingDir), ">", cmd, strings.Join(args, " "))
		fmt.Println(result.Output)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
	return result, err
}

// SetTestOrigin adds the given environment variable to subsequent runs of commands.
func (r *MockingRunner) SetTestOrigin(content string) {
	r.testOrigin = content
}
