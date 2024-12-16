package opcodes

import "github.com/git-town/git-town/v17/internal/vm/shared"

// PullCurrentBranch updates the branch with the given name with commits from its remote.
type PullCurrentBranch struct {
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *PullCurrentBranch) Run(args shared.RunArgs) error {
	return args.Git.Pull(args.Frontend)
}
