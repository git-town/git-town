// Package datatable supports comparing Gherkin tables in test code.
package datatable

import "github.com/git-town/git-town/v17/internal/git/gitdomain"

type runner interface {
	SHAsForCommit(name string) gitdomain.SHAs
}
