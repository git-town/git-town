package dialog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(localBranches domain.LocalBranchNames, oldMainBranch domain.LocalBranchName, backend *git.BackendCommands) (domain.LocalBranchName, error) {
	newMainBranchName, err := Select(SelectArgs{
		Options: localBranches.Strings(),
		Message: mainBranchPrompt(oldMainBranch),
		Default: oldMainBranch.String(),
	})
	if err != nil {
		return domain.LocalBranchName{}, err
	}
	newMainBranch := domain.NewLocalBranchName(newMainBranchName)
	return newMainBranch, backend.Config.SetMainBranch(newMainBranch)
}

func mainBranchPrompt(mainBranch domain.LocalBranchName) string {
	result := "Please specify the main development branch:"
	if !mainBranch.IsEmpty() {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch.String())
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
