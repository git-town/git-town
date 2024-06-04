package gitdomain

type Runner interface {
	Run(executable string, args ...string) error
}

type ManyRunner interface {
	RunMany(commands [][]string) error
}

type RunnerQuerier interface {
	Runner
	ManyRunner
	Querier
}
