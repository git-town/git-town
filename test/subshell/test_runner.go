package subshell

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/git-town/git-town/v9/test/envvars"
	"github.com/kballard/go-shellquote"
)

// TestRunner runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
//   - overide environment variables
//   - Temporarily override certain shell commands with mock implementations.
//     Temporary mocks are only valid for the next command being run.
type TestRunner struct {
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
func (r *TestRunner) createBinDir() {
	if r.usesBinDir {
		// binDir already created --> nothing to do here
		return
	}
	asserts.NoError(os.Mkdir(r.BinDir, 0o700))
	r.usesBinDir = true
}

// createMockBinary creates an executable with the given name and content in ms.binDir.
func (r *TestRunner) createMockBinary(name string, content string) {
	r.createBinDir()
	//nolint:gosec // intentionally creating an executable here
	asserts.NoError(os.WriteFile(filepath.Join(r.BinDir, name), []byte(content), 0x744))
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (r *TestRunner) MockBrokenCommand(name string) {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(r.BinDir, name))
	r.createMockBinary("which", content)
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	r.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (r *TestRunner) MockCommand(name string) {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(r.BinDir, name))
	r.createMockBinary("which", content)
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	r.createMockBinary(name, content)
}

// MockCommitMessage sets up this runner with an editor that enters the given commit message.
func (r *TestRunner) MockCommitMessage(message string) {
	r.gitEditor = "git_editor"
	r.createMockBinary(r.gitEditor, fmt.Sprintf("#!/usr/bin/env bash\n\necho %q > $1", message))
}

// MockGit pretends that this repo has Git in the given version installed.
func (r *TestRunner) MockGit(version string) {
	if runtime.GOOS == "windows" {
		// create Windows binary
		content := fmt.Sprintf("echo git version %s\n", version)
		r.createMockBinary("git.cmd", content)
		return
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
		log.Fatalf("cannot locate the git executable: %v", err)
	}
	content := fmt.Sprintf(mockGit, version, gitPath)
	r.createMockBinary("git", content)
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (r *TestRunner) MockNoCommandsInstalled() {
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	r.createMockBinary("which", content)
}

// MustQuery provides the output of the given command with the given arguments.
// Overrides will be used and removed when done.
func (r *TestRunner) MustQuery(name string, arguments ...string) string {
	return r.MustQueryWith(&Options{}, name, arguments...)
}

func (r *TestRunner) MustQueryStringCode(fullCmd string) (string, int) {
	return r.MustQueryStringCodeWith(fullCmd, &Options{})
}

func (r *TestRunner) MustQueryStringCodeWith(fullCmd string, opts *Options) (string, int) {
	parts, err := shellquote.Split(fullCmd)
	asserts.NoError(err)
	cmd, args := parts[0], parts[1:]
	output, exitCode, err := r.QueryWithCode(opts, cmd, args...)
	asserts.NoError(err)
	return output, exitCode
}

// MustQueryWith provides the output of the given command and didn't encounter any form of error.
func (r *TestRunner) MustQueryWith(opts *Options, cmd string, args ...string) string {
	output, err := r.QueryWith(opts, cmd, args...)
	asserts.NoError(err)
	return output
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (r *TestRunner) MustRun(name string, arguments ...string) {
	err := r.Run(name, arguments...)
	if err != nil {
		panic(fmt.Sprintf("Error executing \"%s %v\": %v", name, arguments, err))
	}
}

func (r *TestRunner) MustRunMany(commands [][]string) {
	asserts.NoError(r.RunMany(commands))
}

// Query provides the output of the given command.
// Overrides will be used and removed when done.
func (r *TestRunner) Query(name string, arguments ...string) (string, error) {
	return r.QueryWith(&Options{}, name, arguments...)
}

// QueryString runs the given command (including possible arguments).
// Overrides will be used and removed when done.
func (r *TestRunner) QueryString(fullCmd string) (string, error) {
	return r.QueryStringWith(fullCmd, &Options{})
}

// QueryStringWith runs the given command (including possible arguments) using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Overrides will be used and removed when done.
func (r *TestRunner) QueryStringWith(fullCmd string, opts *Options) (string, error) {
	parts, err := shellquote.Split(fullCmd)
	asserts.NoError(err)
	cmd, args := parts[0], parts[1:]
	return r.QueryWith(opts, cmd, args...)
}

// QueryWith provides the output of the given command and ensures it exited with code 0.
func (r *TestRunner) QueryWith(opts *Options, cmd string, args ...string) (string, error) {
	output, exitCode, err := r.QueryWithCode(opts, cmd, args...)
	if exitCode != 0 {
		err = fmt.Errorf("process \"%s %s\" failed with code %d, output:\n%s", cmd, strings.Join(args, " "), exitCode, output)
	}
	return output, err
}

// QueryWith runs the given command with the given options in this ShellRunner's directory.
func (r *TestRunner) QueryWithCode(opts *Options, cmd string, args ...string) (string, int, error) {
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
	asserts.NoError(err)
	asserts.NoError(subProcess.Start())
	for _, userInput := range opts.Input {
		// Here we simply wait for some time until the subProcess needs the input.
		// Capturing the output and scanning for the actual content needed
		// would introduce substantial amounts of multi-threaded complexity
		// for not enough gains.
		// https://github.com/git-town/go-execplus could help make this more robust.
		time.Sleep(500 * time.Millisecond)
		_, err := input.Write([]byte(userInput))
		if err != nil {
			log.Fatalf("can't write %q to subprocess '%s %s': %v", userInput, cmd, strings.Join(args, " "), err)
		}
	}
	err = subProcess.Wait()
	var exitCode int
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
			err = nil
		} else {
			err = subshell.ErrorDetails(cmd, args, err, output.Bytes())
		}
	}
	if r.Debug {
		fmt.Println(filepath.Base(r.WorkingDir), ">", cmd, strings.Join(args, " "))
		os.Stdout.Write(output.Bytes())
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
	if opts.IgnoreOutput {
		return "", exitCode, err
	}
	return strings.TrimSpace(output.String()), exitCode, err
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (r *TestRunner) Run(name string, arguments ...string) error {
	_, err := r.QueryWith(&Options{IgnoreOutput: true}, name, arguments...)
	return err
}

func (r *TestRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := r.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// SetTestOrigin adds the given environment variable to subsequent runs of commands.
func (r *TestRunner) SetTestOrigin(content string) {
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

	// when set, captures the output and returns it
	IgnoreOutput bool
}
