package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func switchCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "switch",
		Short: "Displays the local branches visually and allows switching between them",
		Run: func(cmd *cobra.Command, args []string) {
			currentBranch, err := repo.Silent.CurrentBranch()
			if err != nil {
				cli.Exit(err)
			}
			newBranch, err := queryBranch(currentBranch, repo)
			if err != nil {
				cli.Exit(err)
			}
			if newBranch != nil && *newBranch != currentBranch {
				err = repo.Silent.CheckoutBranch(*newBranch)
				if err != nil {
					cli.Exit(err)
				}
			}
		},
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateIsRepository(repo); err != nil {
				return err
			}
			return validateIsConfigured(repo)
		},
		GroupID: "basic",
	}
}

// queryBranch lets the user select a new branch via a visual dialog.
// Returns the selected branch or nil if the user aborted.
func queryBranch(currentBranch string, repo *git.ProdRepo) (selection *string, err error) { //nolint:nonamedreturns
	entries, err := createEntries(repo)
	if err != nil {
		return nil, err
	}
	return dialog.ModalSelect(entries, currentBranch)
}

// createEntries provides all the entries for the branch dialog.
func createEntries(repo *git.ProdRepo) (dialog.ModalEntries, error) {
	entries := dialog.ModalEntries{}
	var err error
	for _, root := range repo.Config.BranchAncestryRoots() {
		entries, err = addEntryAndChildren(entries, root, 0, repo)
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

// addEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func addEntryAndChildren(entries dialog.ModalEntries, branch string, indent int, repo *git.ProdRepo) (dialog.ModalEntries, error) {
	entries = append(entries, dialog.ModalEntry{
		Text:  strings.Repeat("  ", indent) + branch,
		Value: branch,
	})
	var err error
	for _, child := range repo.Silent.Config.ChildBranches(branch) {
		entries, err = addEntryAndChildren(entries, child, indent+1, repo)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}
