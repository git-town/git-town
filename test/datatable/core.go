package datatable

type runner interface {
	SHAForCommit(name string) string
}
