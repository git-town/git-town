package subshell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/git-town/git-town/v8/src/subshell"
	"github.com/git-town/git-town/v8/test/envvars"
	"github.com/kballard/go-shellquote"
)

// Mocking runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
//   - overide environment variables
//   - Temporarily override certain shell commands with mock implementations.
//     Temporary mocks are only valid for the next command being run.
type Mocking struct {
	// the directory that contains mock executables, ignored if empty
	BinDir string

	// whether to log the output of subshell commands
	Debug bool `exhaustruct:"optional"`

	// name of the binary to use as the custom editor during "git commit"
	gitEditor string `exhaustruct:"optional"`

	// the directory that contains the global Git configuration
	HomeDir string

	// optional content of the GIT_TOWN_REMOTE environment variable
	testOrigin string `exhaustruct:"optional"`

	// indicates whether the current test has created the binDir
	usesBinDir bool `exhaustruct:"optional"`

	// the directory in which this runner executes shell commands
	WorkingDir string
}

// createBinDir creates the directory that contains mock executables.
// This method is idempotent.
func (r *Mocking) createBinDir() error {
	if r.usesBinDir {
		// binDir already created --> nothing to do here
		return nil
	}
	err := os.Mkdir(r.BinDir, 0o700)
	if err != nil {
		return fmt.Errorf("cannot create mock bin dir: %w", err)
	}
	r.usesBinDir = true
	return nil
}

// createMockBinary creates an executable with the given name and content in ms.binDir.
func (r *Mocking) createMockBinary(name string, content string) error {
	if err := r.createBinDir(); err != nil {
		return err
	}
	err := os.WriteFile(filepath.Join(r.BinDir, name), []byte(content), 0o500)
	if err != nil {
		return fmt.Errorf("cannot write custom %q command: %w", name, err)
	}
	return nil
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (r *Mocking) MockBrokenCommand(name string) error {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(r.BinDir, name))
	err := r.createMockBinary("which", content)
	if err != nil {
		return err
	}
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	return r.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (r *Mocking) MockCommand(name string) error {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(r.BinDir, name))
	if err := r.createMockBinary("which", content); err != nil {
		return fmt.Errorf("cannot write custom which command: %w", err)
	}
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	return r.createMockBinary(name, content)
}

// MockGit pretends that this repo has Git in the given version installed.
func (r *Mocking) MockGit(version string) error {
	if runtime.GOOS == "windows" {
		// create Windows binary
		content := fmt.Sprintf("echo git version %s\n", version)
		return r.createMockBinary("git.cmd", content)
	}
	// create Unix binary
	mockGit := `#!/usr/bin/env bash
if [ "$1" = "version" ]; then
  echo "git version %s"
else
	%s "$@"
fi`
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("cannot locate the git executable: %w", err)
	}
	content := fmt.Sprintf(mockGit, version, gitPath)
	return r.createMockBinary("git", content)
}

// MockCommitMessage sets up this runner with an editor that enters the given commit message.
func (r *Mocking) MockCommitMessage(message string) error {
	r.gitEditor = "git_editor"
	return r.createMockBinary(r.gitEditor, fmt.Sprintf("#!/usr/bin/env bash\n\necho %q > $1", message))
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (r *Mocking) MockNoCommandsInstalled() error {
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	return r.createMockBinary("which", content)
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (r *Mocking) Run(name string, arguments ...string) (string, error) {
	return r.RunWith(&Options{}, name, arguments...)
}

// RunMany runs all given commands.
// Commands are provided as a list of argv-style strings.
// Overrides apply for the first command only.
// Failed commands abort immediately with the encountered error.
func (r *Mocking) RunMany(commands [][]string) error {
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
func (r *Mocking) RunString(fullCmd string) (string, error) {
	return r.RunStringWith(fullCmd, &Options{})
}

// RunStringWith runs the given command (including possible arguments) using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Overrides will be used and removed when done.
func (r *Mocking) RunStringWith(fullCmd string, opts *Options) (string, error) {
	parts, err := shellquote.Split(fullCmd)
	if err != nil {
		return "", fmt.Errorf("cannot split command %q: %w", fullCmd, err)
	}
	cmd, args := parts[0], parts[1:]
	return r.RunWith(opts, cmd, args...)
}

// RunWith runs the given command with the given options in this ShellRunner's directory.
func (r *Mocking) RunWith(opts *Options, cmd string, args ...string) (string, error) {
	// create an environment with the temp Overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	opts.Env = envvars.Replace(opts.Env, "HOME", r.HomeDir)
	// add the custom origin
	if r.testOrigin != "" {
		opts.Env = envvars.Replace(opts.Env, "GIT_TOWN_REMOTE", r.testOrigin)
	}
	// add the custom bin dir to the PATH
	if r.usesBinDir {
		opts.Env = envvars.PrependPath(opts.Env, r.BinDir)
	}
	// add the custom GIT_EDITOR
	if r.gitEditor != "" {
		opts.Env = envvars.Replace(opts.Env, "GIT_EDITOR", filepath.Join(r.BinDir, "git_editor"))
	}
	// set the working dir
	opts.Dir = filepath.Join(r.WorkingDir, opts.Dir)
	// run the command inside the custom environment
	subProcess := exec.Command(cmd, args...) // #nosec
	if opts.Dir != "" {
		subProcess.Dir = opts.Dir
	}
	if opts.Env != nil {
		subProcess.Env = opts.Env
	}
	var output bytes.Buffer
	subProcess.Stdout = &output
	subProcess.Stderr = &output
	input, err := subProcess.StdinPipe()
	if err != nil {
		return "", err
	}
	err = subProcess.Start()
	if err != nil {
		return "", fmt.Errorf("can't start subprocess '%s %s': %w", cmd, strings.Join(args, " "), err)
	}
	for _, userInput := range opts.Input {
		// Here we simply wait for some time until the subProcess needs the input.
		// Capturing the output and scanning for the actual content needed
		// would introduce substantial amounts of multi-threaded complexity
		// for not enough gains.
		// https://github.com/git-town/go-execplus could help make this more robust.
		time.Sleep(500 * time.Millisecond)
		_, err := input.Write([]byte(userInput))
		if err != nil {
			return "", fmt.Errorf("can't write %q to subprocess '%s %s': %w", userInput, cmd, strings.Join(args, " "), err)
		}
	}
	err = subProcess.Wait()
	if err != nil {
		err = subshell.ErrorDetails(cmd, args, err, output.Bytes())
	}
	exitCode := subProcess.ProcessState.ExitCode()
	if r.Debug {
		fmt.Println(filepath.Base(r.WorkingDir), ">", cmd, strings.Join(args, " "))
		fmt.Println(output.String())
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}

	if exitCode != 0 {
		err = fmt.Errorf("process \"%s %s\" failed with code %d, output:\n%s", cmd, strings.Join(args, " "), exitCode, output.String())
	}
	return strings.TrimSpace(output.String()), err
}

// SetTestOrigin adds the given environment variable to subsequent runs of commands.
func (r *Mocking) SetTestOrigin(content string) {
	r.testOrigin = content
}

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {
	// Dir contains the directory in which to execute the command.
	// If empty, runs in the current directory.
	Dir string

	// Env allows to override the environment variables to use in the subshell, in the format provided by os.Environ()
	// If empty, uses the environment variables of this process.
	Env []string

	// Input contains the user input to enter into the running command.
	// It is written to the subprocess one element at a time, with a delay defined by command.InputDelay in between.
	Input []string // input into the subprocess
}

// WorkingDir provides the folder that this shell operates in.
func (r *Mocking) Dir() string {
	return r.WorkingDir
}
