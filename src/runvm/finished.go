package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/persistence"
)

// finished is called when executing all steps has successfully finished.
func finished(args ExecuteArgs) error {
	args.RunState.MarkAsFinished()
	if args.RunState.IsAbort || args.RunState.IsUndo {
		err := persistence.Delete(args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateDeleteProblem, err)
		}
	} else {
		err := persistence.Save(args.RunState, args.RootDir)
		if err != nil {
			return fmt.Errorf(messages.RunstateSaveProblem, err)
		}
	}
	fmt.Println()
	args.Run.Stats.PrintAnalysis()
	return nil
}
