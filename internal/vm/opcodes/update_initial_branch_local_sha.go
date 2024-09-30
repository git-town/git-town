package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/shared"
)

type UpdateInitialBranchLocalSHA struct {
	Branch                  gitdomain.LocalBranchName
	SHA                     gitdomain.SHA
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *UpdateInitialBranchLocalSHA) Run(args shared.RunArgs) error {
	return args.UpdateInitialBranchLocalSHA(self.Branch, self.SHA)
}
