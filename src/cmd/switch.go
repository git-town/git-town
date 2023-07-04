package cmd

import (
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/spf13/cobra"
)

const switchDesc = "Displays the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   switchDesc,
		Long:    long(switchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSwitch(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func runSwitch(debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: true,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
		ValidateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	currentBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	branchParents := run.Config.ParentBranchMap()
	newBranch, err := queryBranch(currentBranch, branchParents, &run)
	if err != nil {
		return err
	}
	if newBranch != nil && *newBranch != currentBranch {
		err = run.Backend.CheckoutBranch(*newBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

// queryBranch lets the user select a new branch via a visual dialog.
// Returns the selected branch or nil if the user aborted.
func queryBranch(currentBranch string, branchParents config.BranchParents, run *git.ProdRunner) (selection *string, err error) { //nolint:nonamedreturns
	entries, err := createEntries(branchParents, run)
	if err != nil {
		return nil, err
	}
	return dialog.ModalSelect(entries, currentBranch)
}

// createEntries provides all the entries for the branch dialog.
func createEntries(branchParents config.BranchParents, run *git.ProdRunner) (dialog.ModalEntries, error) {
	entries := dialog.ModalEntries{}
	var err error
	for _, root := range run.Config.BranchAncestryRoots(branchParents) {
		entries, err = addEntryAndChildren(entries, root, 0, run)
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

// addEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func addEntryAndChildren(entries dialog.ModalEntries, branch string, indent int, run *git.ProdRunner) (dialog.ModalEntries, error) {
	entries = append(entries, dialog.ModalEntry{
		Text:  strings.Repeat("  ", indent) + branch,
		Value: branch,
	})
	var err error
	for _, child := range run.Config.ChildBranches(branch) {
		entries, err = addEntryAndChildren(entries, child, indent+1, run)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}
