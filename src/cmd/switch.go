package cmd

import (
	"strings"

	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func switchCmd() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   "Displays the local branches visually and allows switching between them",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSwitch(debug)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runSwitch(debug bool) error {
	repo, err := Repo(RepoArgs{
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
		validateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return err
	}
	newBranch, err := queryBranch(currentBranch, &repo)
	if err != nil {
		return err
	}
	if newBranch != nil && *newBranch != currentBranch {
		err = repo.CheckoutBranch(*newBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

// queryBranch lets the user select a new branch via a visual dialog.
// Returns the selected branch or nil if the user aborted.
func queryBranch(currentBranch string, repo *git.PublicRepo) (selection *string, err error) { //nolint:nonamedreturns
	entries, err := createEntries(repo)
	if err != nil {
		return nil, err
	}
	return dialog.ModalSelect(entries, currentBranch)
}

// createEntries provides all the entries for the branch dialog.
func createEntries(repo *git.PublicRepo) (dialog.ModalEntries, error) {
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
func addEntryAndChildren(entries dialog.ModalEntries, branch string, indent int, repo *git.PublicRepo) (dialog.ModalEntries, error) {
	entries = append(entries, dialog.ModalEntry{
		Text:  strings.Repeat("  ", indent) + branch,
		Value: branch,
	})
	var err error
	for _, child := range repo.Config.ChildBranches(branch) {
		entries, err = addEntryAndChildren(entries, child, indent+1, repo)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}
