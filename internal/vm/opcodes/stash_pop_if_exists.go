package opcodes

import (
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

// StashPopIfExists restores stashed away changes into the workspace.
type StashPopIfExists struct{}

func (self *StashPopIfExists) Run(args shared.RunArgs) error {
	stashSize, err := args.Git.StashSize(args.Backend)
	if err != nil {
		return err
	}
	if stashSize == 0 && !args.Config.Value.NormalConfig.DryRun {
		return nil
	}
	args.PrependOpcodes(
		&StashPop{},
		&ChangesUnstageAll{},
	)
	return nil
}
