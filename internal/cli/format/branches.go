package format

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

func BranchNames(branches gitdomain.LocalBranchNames) string {
	if len(branches) == 0 {
		return messages.DialogResultNone
	}
	return branches.Join(", ")
}
