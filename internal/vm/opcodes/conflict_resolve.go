package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConflictResolve struct {
	FilePath                string
	Resolution              gitdomain.ConflictResolution
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictResolve) Run(args shared.RunArgs) error {
	if err := args.Git.ResolveConflict(args.Frontend, self.FilePath, self.Resolution); err != nil {
		return err
	}
	return args.Git.StageFiles(args.Frontend, self.FilePath)
}
