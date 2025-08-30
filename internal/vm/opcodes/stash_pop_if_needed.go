package opcodes

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

// StashPopIfExists restores stashed away changes into the workspace.
type StashPopIfNeeded struct {
	InitialStashSize        gitdomain.StashSize
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *StashPopIfNeeded) Run(args shared.RunArgs) error {
	stashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	if stashSize <= self.InitialStashSize && !args.Config.Value.NormalConfig.DryRun {
		return nil
	}
	args.PrependOpcodes(
		&StashPopIfExists{},
	)
	return nil
}
