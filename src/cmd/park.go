package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const parkDesc = "Suspends syncing some feature branches"

const parkHelp = `
Git Town does not sync parked branches.
The only exception is the currently checked out branch.

If branches are given, parks the given branches.
If no branch is given, parks the current branch.
`

func parkCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "park [branches]",
		Args:    cobra.ArbitraryArgs,
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

func executePark(args []string, verbose bool) error {
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
	config, exit, err := determineParkConfig(args, repo, verbose)
	if err != nil && exit {
		return err
	}
	for _, branchToPark := range config.branchesToPark {
		if !config.branches.HasLocalBranch(branchToPark) {
			return fmt.Errorf(messages.BranchDoesntExist, &branchToPark)
		}
		if err = validateIsParkableBranch(branchToPark, &repo.Runner.Config.FullConfig); err != nil {
			return err
		}
	}
	for _, branchToPark := range config.branchesToPark {
		if err = repo.Runner.Config.AddToParkedBranches(branchToPark); err != nil {
			return err
		}
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

type parkConfig struct {
	branchesToPark gitdomain.LocalBranchNames
	branches       gitdomain.BranchInfos
}

func determineParkConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (parkConfig, bool, error) {
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      nil,
		Fetch:                 false,
		FullConfig:            &repo.Runner.Config.FullConfig,
		HandleUnfinishedState: false,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return parkConfig{}, exit, err
	}
	var branchesToPark gitdomain.LocalBranchNames
	if len(args) == 0 {
		currentBranch, err := repo.Runner.Backend.CurrentBranch()
		if err != nil {
			return parkConfig{}, false, err
		}
		branchesToPark = gitdomain.LocalBranchNames{currentBranch}
	} else {
		branchesToPark = make(gitdomain.LocalBranchNames, len(args))
		for b, branchName := range args {
			branchesToPark[b] = gitdomain.NewLocalBranchName(branchName)
		}
	}
	return parkConfig{
		branchesToPark: branchesToPark,
		branches:       branchesSnapshot.Branches,
	}, false, nil
}

func validateIsParkableBranch(branch gitdomain.LocalBranchName, config *configdomain.FullConfig) error {
	if config.IsContributionBranch(branch) {
		return errors.New(messages.ContributionBranchCannotPark)
	}
	if config.IsMainBranch(branch) {
		return errors.New(messages.MainBranchCannotPark)
	}
	if config.IsObservedBranch(branch) {
		return errors.New(messages.ObservedBranchCannotPark)
	}
	if config.IsPerennialBranch(branch) {
		return errors.New(messages.PerennialBranchCannotPark)
	}
	return nil
}
