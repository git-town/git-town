package commands

type Shell interface {
	Run(string, ...string) (string, error)
	RunMany([][]string) error
}
