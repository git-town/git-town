package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ConflictResolve struct {
	FilePath   string
	Resolution gitdomain.ConflictResolution
}

func (self *ConflictResolve) Run(args shared.RunArgs) error {
	if err := args.Git.ResolveConflict(args.Frontend, self.FilePath, self.Resolution); err != nil {
		return err
	}
	return args.Git.StageFiles(args.Frontend, self.FilePath)
}
