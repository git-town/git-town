package cmd

import (
	"errors"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const setParentDesc = "Prompts to set the parent branch for the current branch"

func setParentCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "set-parent",
		GroupID: "lineage",
		Args:    cobra.NoArgs,
		Short:   setParentDesc,
		Long:    long(setParentDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setParent(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func setParent(debug bool) error {
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
	if !run.Config.IsFeatureBranch(currentBranch) {
		return errors.New("only feature branches can have parent branches")
	}
	existingParent := run.Config.Lineage().Parent(currentBranch)
	if existingParent != "" {
		// TODO: delete the old parent only when the user has entered a new parent
		err = run.Config.RemoveParent(currentBranch)
		if err != nil {
			return err
		}
	} else {
		existingParent = run.Config.MainBranch()
	}
	err = validate.KnowsBranchAncestors(currentBranch, existingParent, &run.Backend)
	if err != nil {
		return err
	}
	run.Stats.PrintAnalysis()
	return nil
}
