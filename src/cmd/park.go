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

const parkDesc = "Suspends syncing of selected feature branches"

const parkHelp = `
Parks the given branches.
If no branch is provided, parks the current branch.

Git Town does not sync a parked branch
unless it is currently checked out.
Only feature branches can be parked.
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
	config, err := determineParkConfig(args, repo, verbose)
	if err != nil {
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

func determineParkConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (parkConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return parkConfig{}, err
	}
	var branchesToPark gitdomain.LocalBranchNames
	if len(args) == 0 {
		branchesToPark = gitdomain.LocalBranchNames{branchesSnapshot.Active}
	} else {
		branchesToPark = make(gitdomain.LocalBranchNames, len(args))
		for b, branchName := range args {
			branchesToPark[b] = gitdomain.NewLocalBranchName(branchName)
		}
	}
	return parkConfig{
		branchesToPark: branchesToPark,
		branches:       branchesSnapshot.Branches,
	}, nil
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
