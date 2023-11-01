package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v10/src/cli/dialog"
	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/spf13/cobra"
)

const switchDesc = "Displays the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   switchDesc,
		Long:    long(switchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeSwitch(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwitch(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return err
	}
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	newBranch, validChoice, err := queryBranch(branches.Initial, lineage)
	if err != nil {
		return err
	}
	if validChoice && newBranch != branches.Initial {
		fmt.Println()
		err = repo.Runner.Frontend.CheckoutBranch(newBranch)
		if err != nil {
			exitCode := 1
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				exitCode = exitErr.ExitCode()
			}
			os.Exit(exitCode)
		}
	}
	return nil
}

// queryBranch lets the user select a new branch via a visual dialog.
// Indicates via `validSelection` whether the user made a valid selection.
func queryBranch(currentBranch domain.LocalBranchName, lineage config.Lineage) (selection domain.LocalBranchName, validSelection bool, err error) {
	entries, err := createEntries(lineage, currentBranch)
	if err != nil {
		return domain.EmptyLocalBranchName(), false, err
	}
	choice, err := dialog.ModalSelect(entries, currentBranch.String())
	if err != nil {
		return domain.EmptyLocalBranchName(), false, err
	}
	if choice == nil {
		return domain.EmptyLocalBranchName(), false, nil
	}
	return domain.NewLocalBranchName(*choice), true, nil
}

// createEntries provides all the entries for the branch dialog.
func createEntries(lineage config.Lineage, currentBranch domain.LocalBranchName) (dialog.ModalSelectEntries, error) {
	entries := dialog.ModalSelectEntries{}
	var err error
	for _, root := range lineage.Roots() {
		entries, err = addEntryAndChildren(entries, root, 0, lineage)
		if err != nil {
			return nil, err
		}
	}
	if len(entries) == 0 {
		entries = append(entries, dialog.ModalSelectEntry{
			Text:  string(currentBranch),
			Value: string(currentBranch),
		})
	}
	return entries, nil
}

// addEntryAndChildren adds the given branch and all its child branches to the given entries collection.
func addEntryAndChildren(entries dialog.ModalSelectEntries, branch domain.LocalBranchName, indent int, lineage config.Lineage) (dialog.ModalSelectEntries, error) {
	entries = append(entries, dialog.ModalSelectEntry{
		Text:  strings.Repeat("  ", indent) + branch.String(),
		Value: branch.String(),
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
