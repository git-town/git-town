// Package subshelldomain defines types around subshells.
package subshelldomain

import "github.com/git-town/git-town/v23/internal/gohacks"

type Runner interface {
	Run(executable string, args ...string) error
	RunWithEnv(env []string, executable string, args ...string) error
}

type Querier interface {
	Query(executable string, args ...string) (string, error)
	QueryZ(executable string, args ...string) (gohacks.ZString, error)
	QueryTrim(executable string, args ...string) (string, error)
}

type RunnerQuerier interface {
	Runner
	Querier
}
