package dialog

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// EnterMainBranch lets the user select a new main branch for this repo.
// This includes asking the user and updating the respective setting.
func EnterMainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName, backend *git.BackendCommands) (gitdomain.LocalBranchName, error) {
	newMainBranchName, err := Select(SelectArgs{
		Options: localBranches.Strings(),
		Message: mainBranchPrompt(oldMainBranch),
		Default: oldMainBranch.String(),
	})
	if err != nil {
		return gitdomain.EmptyLocalBranchName(), err
	}
	newMainBranch := gitdomain.NewLocalBranchName(newMainBranchName)
	err = backend.GitTown.SetMainBranch(newMainBranch)
	return newMainBranch, err
}

func mainBranchPrompt(mainBranch gitdomain.LocalBranchName) string {
	result := "Please specify the main development branch:"
	if !mainBranch.IsEmpty() {
		coloredBranch := color.New(color.Bold).Add(color.FgCyan).Sprintf(mainBranch.String())
		result += fmt.Sprintf(" (current value: %s)", coloredBranch)
	}
	return result
}
