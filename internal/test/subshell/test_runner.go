package subshell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/bytestream"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/test/envvars"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/kballard/go-shellquote"
)

// TestRunner runs shell commands using a customizable environment.
// This is useful in tests. Possible customizations:
//   - override environment variables
//   - Temporarily override certain shell commands with mock implementations.
//     Temporary mocks are only valid for the next command being run.
type TestRunner struct {
	// the directory that contains mock executables, ignored if empty
	BinDir string

	// the directory that contains the global Git configuration
	HomeDir string

	// whether to log the output of subshell commands
	Verbose configdomain.Verbose

	// the directory in which this runner executes shell commands
	WorkingDir string

	// name of the binary to use as the custom editor during "git commit"
	gitEditor Option[string]

	// content of the GIT_TOWN_REMOTE environment variable
	testOrigin Option[string]

	// indicates whether the current test has created the binDir
	usesBinDir bool
}

// MockBrokenCommand adds a mock for the given command that returns an error.
func (self *TestRunner) MockBrokenCommand(name string) {
	content := "#!/usr/bin/env bash\n\nexit 1"
	self.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (self *TestRunner) MockCommand(name string) {
	content := fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
	self.createMockBinary(name, content)
}

// MockCommitMessage sets up this runner with an editor that enters the given commit message.
func (self *TestRunner) MockCommitMessage(message string) {
	editorPath := "git_editor"
	self.gitEditor = Some(editorPath)
	self.createMockBinary(editorPath, fmt.Sprintf("#!/usr/bin/env bash\n\necho %q > $1", message))
}

// MockGit pretends that this repo has Git in the given version installed.
func (self *TestRunner) MockGit(version string) {
	if runtime.GOOS == "windows" {
		// create Windows binary
		content := fmt.Sprintf("echo git version %s\n", version)
		self.createMockBinary("git.cmd", content)
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
	asserts.NoError(err)
	content := fmt.Sprintf(mockGit, version, gitPath)
	self.createMockBinary("git", content)
}

// MockNoCommandsInstalled pretends that no commands are installed.
func (self *TestRunner) MockNoCommandsInstalled() {
	content := "#!/usr/bin/env bash\n\nexit 1\n"
	self.createMockBinary("which", content)
}

// MustQuery provides the output of the given command with the given arguments.
// Overrides will be used and removed when done.
func (self *TestRunner) MustQuery(name string, arguments ...string) string {
	return self.MustQueryWith(&Options{TTY: true}, name, arguments...)
}

func (self *TestRunner) MustQueryStringCode(fullCmd string) RunResult {
	return self.MustQueryStringCodeWith(fullCmd, &Options{TTY: true})
}

func (self *TestRunner) MustQueryStringCodeWith(fullCmd string, opts *Options) RunResult {
	parts := asserts.NoError1(shellquote.Split(fullCmd))
	cmd, args := parts[0], parts[1:]
	return asserts.NoError1(self.QueryWithCode(opts, cmd, args...))
}

// MustQueryWith provides the output of the given command and didn't encounter any form of error.
func (self *TestRunner) MustQueryWith(opts *Options, cmd string, args ...string) string {
	return asserts.NoError1(self.QueryWith(opts, cmd, args...))
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (self *TestRunner) MustRun(name string, arguments ...string) {
	output, err := self.Query(name, arguments...)
	if err != nil {
		panic(fmt.Sprintf("Error executing \"%s %v\": %v\n%s", name, arguments, err, output))
	}
}

// Query provides the output of the given command.
// Overrides will be used and removed when done.
func (self *TestRunner) Query(name string, arguments ...string) (string, error) {
	return self.QueryWith(&Options{TTY: true}, name, arguments...)
}

// QueryString runs the given command (including possible arguments).
// Overrides will be used and removed when done.
func (self *TestRunner) QueryString(fullCmd string) (string, error) {
	return self.QueryStringWith(fullCmd, &Options{TTY: true})
}

// QueryStringWith runs the given command (including possible arguments) using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Overrides will be used and removed when done.
func (self *TestRunner) QueryStringWith(fullCmd string, opts *Options) (string, error) {
	parts := asserts.NoError1(shellquote.Split(fullCmd))
	cmd, args := parts[0], parts[1:]
	return self.QueryWith(opts, cmd, args...)
}

// Query provides the output of the given command.
// Overrides will be used and removed when done.
func (self *TestRunner) QueryTrim(name string, arguments ...string) (string, error) {
	output, err := self.QueryWith(&Options{TTY: true}, name, arguments...)
	return strings.TrimSpace(output), err
}

// QueryWith provides the output of the given command and ensures it exited with code 0.
func (self *TestRunner) QueryWith(opts *Options, cmd string, args ...string) (string, error) {
	runResult, err := self.QueryWithCode(opts, cmd, args...)
	if runResult.ExitCode != 0 {
		err = fmt.Errorf("process \"%s %s\" failed with code %d.\nOUTPUT START\n%s\nOUTPUT END", cmd, strings.Join(args, " "), runResult.ExitCode, runResult.Output)
	}
	return runResult.Output, err
}

// QueryWith runs the given command with the given options in this ShellRunner's directory.
func (self *TestRunner) QueryWithCode(opts *Options, cmd string, args ...string) (RunResult, error) {
	emptyResult := RunResult{ExitCode: 0, Output: ""}
	currentBranchText := ""
	if self.Verbose {
		getBranchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		getBranchCmd.Dir = self.WorkingDir
		currentBranch, _ := getBranchCmd.Output()
		currentBranchText = strings.TrimSpace(string(currentBranch))
	}

	// create an environment with the temp Overrides directory added to the PATH
	if opts.Env == nil {
		opts.Env = os.Environ()
	}
	// set HOME to the given global directory so that Git puts the global configuration there.
	opts.Env = envvars.Replace(opts.Env, "HOME", self.HomeDir)
	// add the custom origin
	if testOrigin, hasTestOrigin := self.testOrigin.Get(); hasTestOrigin {
		opts.Env = envvars.Replace(opts.Env, "GIT_TOWN_REMOTE", testOrigin)
	}
	// add the custom bin dir to the PATH
	if self.usesBinDir {
		opts.Env = envvars.PrependPath(opts.Env, self.BinDir)
	}
	// add the custom GIT_EDITOR
	if gitEditor, hasGitEditor := self.gitEditor.Get(); hasGitEditor {
		opts.Env = envvars.Replace(opts.Env, "GIT_EDITOR", filepath.Join(self.BinDir, gitEditor))
	}
	// mark as test run
	opts.Env = append(opts.Env, subshell.TestToken+"=1")
	// run the command inside the custom environment
	subProcess := exec.Command(cmd, args...) // #nosec
	subProcess.Dir = filepath.Join(self.WorkingDir, opts.Dir)
	subProcess.Env = opts.Env
	var outputBuf bytes.Buffer
	subProcess.Stdout = &outputBuf
	subProcess.Stderr = &outputBuf
	var err error
	if input, hasInput := opts.Input.Get(); hasInput {
		var stdin io.WriteCloser
		stdin, err = subProcess.StdinPipe()
		if err != nil {
			return emptyResult, fmt.Errorf("cannot create stdin pipe: %w", err)
		}
		if err = subProcess.Start(); err != nil {
			return emptyResult, fmt.Errorf("cannot start command: %w", err)
		}
		_, err = stdin.Write([]byte(input))
		if err != nil {
			return emptyResult, fmt.Errorf("cannot write to stdin: %w", err)
		}
		if err = stdin.Close(); err != nil {
			return emptyResult, fmt.Errorf("cannot close stdin pipe: %w", err)
		}
		if err = subProcess.Wait(); err != nil {
			fmt.Println("cannot wait for command to finish:", err)
			return emptyResult, err
		}
	} else {
		subProcess.Stdin = nil
		if !opts.TTY {
			// NOTE: We only disable the TTY when a test explicitly requests it.
			// This keeps TTY behavior explicit and predictable.
			//
			// Some tests simulate dialog input via environment variables rather
			// than providing opts.Input. If we implicitly disabled the TTY when opts.Input is not set,
			// those tests would fail.
			disableTTY(subProcess)
		}
		err = subProcess.Run()
	}
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
			err = nil
		} else {
			err = subshell.ErrorDetails(cmd, args, err, outputBuf.Bytes())
		}
	}
	if self.Verbose {
		fmt.Printf("\n\n%s@%s > %s %s\n\n", strings.ToUpper(filepath.Base(self.WorkingDir)), currentBranchText, cmd, stringslice.JoinArgs(args))
		os.Stdout.Write(bytestream.NullDelineated(outputBuf.Bytes()).ToNewlines())
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
		}
	}
	var output string
	if opts.IgnoreOutput {
		output = ""
	} else {
		output = strings.TrimRight(outputBuf.String(), "\n")
	}
	return RunResult{
		ExitCode: exitCode,
		Output:   output,
	}, err
}

// Run runs the given command with the given arguments.
// Overrides will be used and removed when done.
func (self *TestRunner) Run(name string, arguments ...string) error {
	_, err := self.QueryWith(&Options{IgnoreOutput: true, TTY: true}, name, arguments...)
	return err
}

func (self *TestRunner) RunWithEnv(env []string, name string, arguments ...string) error {
	_, err := self.QueryWith(&Options{Env: env, IgnoreOutput: true, TTY: true}, name, arguments...)
	return err
}

// SetTestOrigin adds the given environment variable to subsequent runs of commands.
func (self *TestRunner) SetTestOrigin(content string) {
	self.testOrigin = Some(content)
}

// createBinDir creates the directory that contains mock executables.
// This method is idempotent.
func (self *TestRunner) createBinDir() {
	if self.usesBinDir {
		// binDir already created --> nothing to do here
		return
	}
	asserts.NoError(os.Mkdir(self.BinDir, 0o700))
	self.usesBinDir = true
}

// createMockBinary creates an executable with the given name and content in ms.binDir.
func (self *TestRunner) createMockBinary(name string, content string) {
	self.createBinDir()
	binaryPath := filepath.Join(self.BinDir, name)
	//nolint:gosec // intentionally creating an executable here
	asserts.NoError(os.WriteFile(binaryPath, []byte(content), 0o744))
}

// Options defines optional arguments for ShellRunner.RunWith().
type Options struct {
	// Dir contains the directory in which to execute the command.
	// If empty, runs in the current directory.
	Dir string `exhaustruct:"optional"`

	// Env allows to override the environment variables to use in the subshell, in the format provided by os.Environ()
	// If empty, uses the environment variables of this process.
	Env []string `exhaustruct:"optional"`

	// when set, captures the output and returns it
	IgnoreOutput bool `exhaustruct:"optional"`

	// input to pipe into STDIN
	Input Option[string] `exhaustruct:"optional"`

	// whether to provide a TTY to the subshell
	TTY bool
}

type RunResult struct {
	ExitCode int
	Output   string
}
