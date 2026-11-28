// Package subshelldomain defines types around subshells.
package subshelldomain

import "github.com/git-town/git-town/v24/internal/gohacks/stringss"

type Runner interface {
	Run(executable string, args ...string) error
	RunWithEnv(env []string, executable string, args ...string) error
}

// RunnerWithInput is implemented by runners that can pipe input into the command they execute.
type RunnerWithInput interface {
	RunWithInput(input string, executable string, args ...string) error
}

type Querier interface {
	Query(executable string, args ...string) (string, error)
	QueryTrim(executable string, args ...string) (stringss.Trimmed, error)
	QueryZ(executable string, args ...string) (stringss.ZeroDelineated, error)
}

type RunnerQuerier interface {
	Runner
	Querier
}
