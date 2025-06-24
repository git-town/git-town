// Package subshelldomain defines types around subshells.
package subshelldomain

type Runner interface {
	Run(executable string, args ...string) error
	RunWithEnv(env []string, executable string, args ...string) error
}

type Querier interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
}

type RunnerQuerier interface {
	Runner
	Querier
}
