package opcode

import (
	"fmt"
)

// CommitOpenChanges commits all open changes as a new commit.
// It does not ask the user for a commit message, but chooses one automatically.
type CommitOpenChanges struct {
	BaseOpcode
}

func (step *CommitOpenChanges) Run(args RunArgs) error {
	err := args.Runner.Frontend.StageFiles("-A")
	if err != nil {
		return err
	}
	currentBranch, err := args.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	return args.Runner.Frontend.CommitStagedChanges(fmt.Sprintf("WIP on %s", currentBranch))
}
