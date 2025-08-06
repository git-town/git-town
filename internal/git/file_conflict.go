package git

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// FileConflict contains information about a file with conflicts, as provided by "git ls-files --unmerged".
type FileConflict struct {
	BaseChange          Option[Blob] // info about the base version of the file (when 3-way merging)
	CurrentBranchChange Option[Blob] // info about the content of the file on the branch where the merge conflict occurs, None == file is deleted here
	IncomingChange      Option[Blob] // info about the content of the file on the branch being merged in, None == file is being deleted here
}

// prints debug information
func (quickInfo FileConflict) Debug(querier subshelldomain.Querier) {
	base, hasBase := quickInfo.BaseChange.Get()
	current, hasCurrent := quickInfo.CurrentBranchChange.Get()
	incoming, hasIncoming := quickInfo.IncomingChange.Get()
	fmt.Print("BASE CHANGE: ")
	if hasBase {
		base.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
	fmt.Print("CURRENT CHANGE: ")
	if hasCurrent {
		current.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
	fmt.Print("INCOMING CHANGE: ")
	if hasIncoming {
		incoming.Debug(querier)
	} else {
		fmt.Println("(none)")
	}
}
