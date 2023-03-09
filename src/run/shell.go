package run

// Shell defines operations to run commands in a subshell.
type Shell interface {
	Run(string, ...string) (*Result, error)
	RunMany([][]string) error
	RunString(string) (*Result, error)
	WorkingDir() string
}
