package cmd

import (
	"errors"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const observeDesc = "Stops your collaboration on specific feature branches"

const parkHelp = `
Your local copy of an observed feature branch only reflects the
and doesn't pull in updates from the parent.


`

func parkCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "park",
		Args:    cobra.NoArgs,
		GroupID: "types",
		Short:   parkDesc,
		Long:    cmdhelpers.Long(parkDesc, parkHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePark(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePark(_ []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	currentBranch, err := repo.Runner.Backend.CurrentBranch()
	if err != nil {
		return err
	}
	if repo.Runner.Config.FullConfig.IsContributionBranch(currentBranch) {
		return errors.New(messages.ContributionBranchCannotPark)
	}
	if repo.Runner.Config.FullConfig.IsMainBranch(currentBranch) {
		return errors.New(messages.MainBranchCannotPark)
	}
	if repo.Runner.Config.FullConfig.IsObservedBranch(currentBranch) {
		return errors.New(messages.ObservedBranchCannotPark)
	}
	if repo.Runner.Config.FullConfig.IsPerennialBranch(currentBranch) {
		return errors.New(messages.PerennialBranchCannotPark)
	}
	err = repo.Runner.Config.AddToParkedBranches(currentBranch)
	if err != nil {
		return err
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "park",
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		RootDir:             repo.RootDir,
		Runner:              repo.Runner,
		Verbose:             verbose,
	})
}
