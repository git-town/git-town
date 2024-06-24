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

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/subshell"
	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/git-town/git-town/v14/test/envvars"
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

	// the directory that contains the global Git configuration
	HomeDir string

	// whether to log the output of subshell commands
	Verbose bool

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
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(self.BinDir, name))
	self.createMockBinary("which", content)
	// write custom command
	content = "#!/usr/bin/env bash\n\nexit 1"
	self.createMockBinary(name, content)
}

// MockCommand adds a mock for the command with the given name.
func (self *TestRunner) MockCommand(name string) {
	// write custom "which" command
	content := fmt.Sprintf("#!/usr/bin/env bash\n\nif [ \"$1\" == %q ]; then\n  echo %q\nelse\n  exit 1\nfi", name, filepath.Join(self.BinDir, name))
	self.createMockBinary("which", content)
	// write custom command
	content = fmt.Sprintf("#!/usr/bin/env bash\n\necho %s called with: \"$@\"\n", name)
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
	if err != nil {
		log.Fatalf("cannot locate the git executable: %v", err)
	}
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
	return self.MustQueryWith(&Options{}, name, arguments...)
}

func (self *TestRunner) MustQueryStringCode(fullCmd string) (output string, exitCode int) {
	return self.MustQueryStringCodeWith(fullCmd, &Options{})
}

func (self *TestRunner) MustQueryStringCodeWith(fullCmd string, opts *Options) (string, int) {
	parts, err := shellquote.Split(fullCmd)
	asserts.NoError(err)
	cmd, args := parts[0], parts[1:]
	output, exitCode, err := self.QueryWithCode(opts, cmd, args...)
	asserts.NoError(err)
	return output, exitCode
}

// MustQueryWith provides the output of the given command and didn't encounter any form of error.
func (self *TestRunner) MustQueryWith(opts *Options, cmd string, args ...string) string {
	output, err := self.QueryWith(opts, cmd, args...)
	asserts.NoError(err)
	return output
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
	return self.QueryWith(&Options{}, name, arguments...)
}

// QueryString runs the given command (including possible arguments).
// Overrides will be used and removed when done.
func (self *TestRunner) QueryString(fullCmd string) (string, error) {
	return self.QueryStringWith(fullCmd, &Options{})
}

// QueryStringWith runs the given command (including possible arguments) using the given options.
// opts.Dir is a relative path inside the working directory of this ShellRunner.
// Overrides will be used and removed when done.
func (self *TestRunner) QueryStringWith(fullCmd string, opts *Options) (string, error) {
	parts, err := shellquote.Split(fullCmd)
	asserts.NoError(err)
	cmd, args := parts[0], parts[1:]
	return self.QueryWith(opts, cmd, args...)
}

// Query provides the output of the given command.
// Overrides will be used and removed when done.
func (self *TestRunner) QueryTrim(name string, arguments ...string) (string, error) {
	output, err := self.QueryWith(&Options{}, name, arguments...)
	return strings.TrimSpace(output), err
}

// QueryWith provides the output of the given command and ensures it exited with code 0.
func (self *TestRunner) QueryWith(opts *Options, cmd string, args ...string) (string, error) {
	output, exitCode, err := self.QueryWithCode(opts, cmd, args...)
	if exitCode != 0 {
		err = fmt.Errorf("process \"%s %s\" failed with code %d, output:\n%s", cmd, strings.Join(args, " "), exitCode, output)
	}
	return output, err
}

// QueryWith runs the given command with the given options in this ShellRunner's directory.
func (self *TestRunner) QueryWithCode(opts *Options, cmd string, args ...string) (string, int, error) {
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
	// set the working dir
	opts.Dir = filepath.Join(self.WorkingDir, opts.Dir)
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
	err := subProcess.Run()
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
	if self.Verbose {
		fmt.Printf("\n\n%s@%s > %s %s\n\n", strings.ToUpper(filepath.Base(self.WorkingDir)), currentBranchText, cmd, stringslice.JoinArgs(args))
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
func (self *TestRunner) Run(name string, arguments ...string) error {
	_, err := self.QueryWith(&Options{IgnoreOutput: true}, name, arguments...)
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
	//nolint:gosec // intentionally creating an executable here
	asserts.NoError(os.WriteFile(filepath.Join(self.BinDir, name), []byte(content), 0o744))
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
}
