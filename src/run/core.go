// Package run provides facilities to execute CLI commands in subshells.
package run

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

// Exec executes the command given in argv notation.
func Exec(cmd string, args ...string) (*Result, error) {
	return WithOptions(&Options{}, cmd, args...)
}

// InDir executes the given command in the given directory.
func InDir(dir string, cmd string, args ...string) (*Result, error) {
	return WithOptions(&Options{Dir: dir}, cmd, args...)
}

// WithOptions runs the command with the given RunOptions.
func WithOptions(opts *Options, cmd string, args ...string) (*Result, error) {
	return &result, err
}
