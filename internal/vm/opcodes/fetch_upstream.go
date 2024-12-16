package opcodes

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/vm/shared"
)

// FetchUpstream brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstream struct {
	Branch                  gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *FetchUpstream) Run(args shared.RunArgs) error {
	return args.Git.FetchUpstream(args.Frontend, self.Branch)
}
