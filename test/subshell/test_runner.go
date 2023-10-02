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
func (tr *TestRunner) createBinDir() {
	if tr.usesBinDir {
		// binDir already created --> nothing to do here
		return
	}
	asserts.NoError(os.Mkdir(tr.BinDir, 0o700))
	tr.usesBinDir = true
}

// createMockBinary creates an executable with the given name and content in ms.binDir.
func (tr *TestRunner) createMockBinary(name string, content string) {
	tr.createBinDir()
	//nolint:gosec // intentionally creating an executable here
	asserts.NoError(os.WriteFile(filepath.Join(tr.BinDir, name), []byte(content), 0x744))
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (tr *TestRunner) MockBrokenCommand(name string) {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(tr.BinDir, name))
	tr.createMockBinary("which", content)
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	tr.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (tr *TestRunner) MockCommand(name string) {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(tr.BinDir, name))
	tr.createMockBinary("which", content)
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	tr.createMockBinary(name, content)
}

// MockCommitMessage sets up this runner with an editor that enters the given commit message.
func (tr *TestRunner) MockCommitMessage(message string) {
	tr.gitEditor = "git_editor"
	tr.createMockBinary(tr.gitEditor, fmt.Sprintf("#!/usr/bin/env bash\n\necho %q > $1", message))
}

// MockGit pretends that this repo has Git in the given version installed.
func (tr *TestRunner) MockGit(version string) {
	if runtime.GOOS == "windows" {
		// create Windows binary
		content := fmt.Sprintf("echo git version %s\n", version)
		tr.createMockBinary("git.cmd", content)
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
	tr.createMockBinary("git", content)
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (tr *TestRunner) MockNoCommandsInstalled() {
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	tr.createMockBinary("which", content)
}

// MustQuery provides the output of the given command with the given arguments.
// Overrides will be used and removed when done.
func (tr *TestRunner) MustQuery(name string, arguments ...string) string {
	return tr.MustQueryWith(&Options{}, name, arguments...)
}

func (tr *TestRunner) MustQueryStringCode(fullCmd string) (string, int) {
	return tr.MustQueryStringCodeWith(fullCmd, &Options{})
}

func (tr *TestRunner) MustQueryStringCodeWith(fullCmd string, opts *Options) (string, int) {
	parts, err := shellquote.Split(fullCmd)
	asserts.NoError(err)
	cmd, args := parts[0], parts[1:]
	output, exitCode, err := tr.QueryWithCode(opts, cmd, args...)
	asserts.NoError(err)
	return output, exitCode
}

// MustQueryWith provides the output of the given command and didn't encounter any form of error.
func (tr *TestRunner) MustQueryWith(opts *Options, cmd string, args ...string) string {
	output, err := tr.QueryWith(opts, cmd, args...)
	asserts.NoError(err)
	return output
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (tr *TestRunner) MustRun(name string, arguments ...string) {
	output, err := tr.Query(name, arguments...)
	if err != nil {
		panic(fmt.Sprintf("Error executing \"%s %v\": %v\n%s", name, arguments, err, output))
	}
}

func (tr *TestRunner) MustRunMany(commands [][]string) {
	asserts.NoError(tr.RunMany(commands))
}

// Query provides the output of the given command.
// Overrides will be used and removed when done.
func (tr *TestRunner) Query(name string, arguments ...string) (string, error) {
	return tr.QueryWith(&Options{}, name, arguments...)
}

// Query provides the output of the given command.
// Overrides will be used and removed when done.
func (tr *TestRunner) QueryTrim(name string, arguments ...string) (string, error) {
	output, err := tr.QueryWith(&Options{}, name, arguments...)
	return strings.TrimSpace(output), err
}

// QueryString runs the given command (including possible arguments).
// Overrides will be used and removed when done.
func (tr *TestRunner) QueryString(fullCmd string) (string, error) {
	return tr.QueryStringWith(fullCmd, &Options{})
}

// QueryStringWith runs the given command (including possible arguments) using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Overrides will be used and removed when done.
func (tr *TestRunner) QueryStringWith(fullCmd string, opts *Options) (string, error) {
	parts, err := shellquote.Split(fullCmd)
	asserts.NoError(err)
	cmd, args := parts[0], parts[1:]
	return tr.QueryWith(opts, cmd, args...)
}

// QueryWith provides the output of the given command and ensures it exited with code 0.
func (tr *TestRunner) QueryWith(opts *Options, cmd string, args ...string) (string, error) {
	output, exitCode, err := tr.QueryWithCode(opts, cmd, args...)
	if exitCode != 0 {
		err = fmt.Errorf("process \"%s %s\" failed with code %d, output:\n%s", cmd, strings.Join(args, " "), exitCode, output)
	}
	return output, err
}

// QueryWith runs the given command with the given options in this ShellRunner's directory.
func (tr *TestRunner) QueryWithCode(opts *Options, cmd string, args ...string) (string, int, error) {
	currentBranchText := ""
	if tr.Debug {
		getBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		getBranchCmd.Dir = tr.WorkingDir
		currentBranch, _ := getBranchCmd.Output()
		currentBranchText = strings.TrimSpace(string(currentBranch))
	}

	// create an environment with the temp Overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	opts.Env = envvars.Replace(opts.Env, "HOME", tr.HomeDir)
	// add the custom origin
	if tr.testOrigin != "" {
		opts.Env = envvars.Replace(opts.Env, "GIT_TOWN_REMOTE", tr.testOrigin)
	}
	// add the custom bin dir to the PATH
	if tr.usesBinDir {
		opts.Env = envvars.PrependPath(opts.Env, tr.BinDir)
	}
	// add the custom GIT_EDITOR
	if tr.gitEditor != "" {
		opts.Env = envvars.Replace(opts.Env, "GIT_EDITOR", filepath.Join(tr.BinDir, "git_editor"))
	}
	// set the working dir
	opts.Dir = filepath.Join(tr.WorkingDir, opts.Dir)
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
			fmt.Printf("\nERROR: can't write %q to subprocess '%s %s': %v\n\n", userInput, cmd, strings.Join(args, " "), err)
			fmt.Printf("OUTPUT: %s\n", output.String())
			return "", 0, errors.New("subprocess crashed")
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
	if tr.Debug {
		fmt.Printf("\n\n%s@%s > %s %s\n\n", strings.ToUpper(filepath.Base(tr.WorkingDir)), currentBranchText, cmd, strings.Join(args, " "))
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
func (tr *TestRunner) Run(name string, arguments ...string) error {
	_, err := tr.QueryWith(&Options{IgnoreOutput: true}, name, arguments...)
	return err
}

func (tr *TestRunner) RunMany(commands [][]string) error {
	for _, argv := range commands {
		err := tr.Run(argv[0], argv[1:]...)
		if err != nil {
			return fmt.Errorf("error running command %q: %w", argv, err)
		}
	}
	return nil
}

// SetTestOrigin adds the given environment variable to subsequent runs of commands.
func (tr *TestRunner) SetTestOrigin(content string) {
	tr.testOrigin = content
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
