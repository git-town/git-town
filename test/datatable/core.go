package datatable

import "github.com/git-town/git-town/v11/src/git/gitdomain"

type runner interface {
	SHAForCommit(name string) gitdomain.SHA
}
