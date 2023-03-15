package subshell

// Runner defines operations to run shell commands.
type Runner interface {
	Run(string, ...string) (*Result, error)
	RunMany([][]string) error
	RunString(string) (*Result, error)
	WorkingDir() string
}
