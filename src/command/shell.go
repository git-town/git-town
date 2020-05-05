package command

// Shell allows running commands in a subshell.
type Shell interface {
	MustRun(string, ...string) *Result
	Run(string, ...string) (*Result, error)
	RunMany([][]string) error
	RunString(string) (*Result, error)
	RunStringWith(string, Options) (*Result, error)
}
