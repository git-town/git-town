package subshell

// Runner defines operations to run shell commands.
type Runner interface {
	Run(string, ...string) (*Output, error)
	RunMany([][]string) error
	RunString(string) (*Output, error)
	WorkingDir() string
}
