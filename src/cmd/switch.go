package cmd

import (
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/dialog"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
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
		OmitBranchNames:       false,
		ValidateGitversion:    true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	currentBranch, err := run.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	newBranch, err := queryBranch(currentBranch, run.Config.Lineage())
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
func queryBranch(currentBranch string, lineage config.Lineage) (selection *string, err error) {
	entries, err := createEntries(lineage)
	if err != nil {
		return nil, err
	}
	return dialog.ModalSelect(entries, currentBranch)
}

// createEntries provides all the entries for the branch dialog.
func createEntries(lineage config.Lineage) (dialog.ModalEntries, error) {
	entries := dialog.ModalEntries{}
	var err error
	for _, root := range lineage.Roots() {
		entries, err = addEntryAndChildren(entries, root, 0, lineage)
		if err != nil {
			return nil, err
		}
	}
	return entries, nil
}

// addEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func addEntryAndChildren(entries dialog.ModalEntries, branch string, indent int, lineage config.Lineage) (dialog.ModalEntries, error) {
	entries = append(entries, dialog.ModalEntry{
		Text:  strings.Repeat("  ", indent) + branch,
		Value: branch,
	})
	var err error
	for _, child := range lineage.Children(branch) {
		entries, err = addEntryAndChildren(entries, child, indent+1, lineage)
		if err != nil {
			return entries, err
		}
	}
	return entries, nil
}
