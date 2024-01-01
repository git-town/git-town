package opcode

import "github.com/git-town/git-town/v11/src/vm/shared"

// PullCurrentBranch updates the branch with the given name with commits from its remote.
type PullCurrentBranch struct {
	undeclaredOpcodeMethods
}

func (self *PullCurrentBranch) Run(args shared.RunArgs) error {
	return args.Runner.Frontend.Pull()
}
