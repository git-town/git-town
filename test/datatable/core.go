package datatable

type Runner interface {
	ShaForCommit(string) (string, error)
}
