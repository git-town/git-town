// Package git provides access to Git.
package git

type Runner interface {
	Run(executable string, args ...string) error
}

type ManyRunner interface {
	RunMany(commands [][]string) error
}

type Querier interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
}

type RunnerQuerier interface {
	Runner
	ManyRunner
	Querier
}
