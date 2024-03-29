package datatable

import "github.com/git-town/git-town/v13/src/git/gitdomain"

type runner interface {
	SHAsForCommit(name string) gitdomain.SHAs
}
