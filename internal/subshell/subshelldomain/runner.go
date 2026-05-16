// Package subshelldomain defines types around subshells.
package subshelldomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

type Runner interface {
	Run(executable string, args ...string) error
	RunWithEnv(env []string, executable string, args ...string) error
}

type Querier interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (string, error)
	QueryZ(executable string, args ...string) (stringss.ZeroDelineated, error)
}

type RunnerQuerier interface {
	Runner
	Querier
}
