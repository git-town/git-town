package cmd

import (
	"errors"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/messages"
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
	run, _, _, exit, err := execute.OpenShell(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	_, currentBranch, err := execute.LoadBranches(&run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	if !run.Config.IsFeatureBranch(currentBranch) {
		return errors.New(messages.SetParentNoFeatureBranch)
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
