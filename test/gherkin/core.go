package gherkin

type Runner interface {
	ShaForCommit(string) (string, error)
}
