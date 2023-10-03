package steps

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/messages"
)

// RemoveFromPerennialBranchesStep removes the branch with the given name as a perennial branch.
type RemoveFromPerennialBranchesStep struct {
	Branch domain.LocalBranchName
	EmptyStep
}

func (step *RemoveFromPerennialBranchesStep) Run(args RunArgs) error {
	fmt.Println(fmt.Sprintf(messages.PerennialBranchRemoved, step.Branch))
	return args.Runner.Config.RemoveFromPerennialBranches(step.Branch)
}
