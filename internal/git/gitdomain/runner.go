package gitdomain

type Runner interface {
	Run(executable string, args ...string) error
	RunWithEnv(env []string, executable string, args ...string) error
}

type RunnerQuerier interface {
	Runner
	Querier
}
