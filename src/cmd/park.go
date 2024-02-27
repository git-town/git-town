package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config"
	"github.com/git-town/git-town/v12/src/config/commandconfig"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const parkDesc = "Suspends syncing of some feature branches"

const parkHelp = `
Parks the given local feature branches.
If no branch is provided, parks the current branch.

Git Town does not sync parked branches.
The currently checked out branch gets synced even if parked.
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
	config, err := determineParkConfig(args, repo)
	if err != nil {
		return err
	}
	err = validateParkConfig(config)
	if err != nil {
		return err
	}
	if err = repo.Runner.Config.AddToParkedBranches(config.branchesToPark.Keys()...); err != nil {
		return err
	}
	if err = removeNonParkBranchTypes(config.branchesToPark, repo.Runner.Config); err != nil {
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

type parkConfig struct {
	allBranches    gitdomain.BranchInfos
	branchesToPark commandconfig.BranchesToMark
}

func removeNonParkBranchTypes(branches map[gitdomain.LocalBranchName]configdomain.BranchType, config *config.Config) error {
	for branchName, branchType := range branches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			if err := config.RemoveFromContributionBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeObservedBranch:
			if err := config.RemoveFromObservedBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determineParkConfig(args []string, repo *execute.OpenRepoResult) (parkConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return parkConfig{}, err
	}
	branchesToPark := commandconfig.BranchesToMark{}
	if len(args) == 0 {
		branchesToPark.Add(branchesSnapshot.Active, &repo.Runner.Config.FullConfig)
	} else {
		for _, branch := range args {
			branchesToPark.Add(gitdomain.NewLocalBranchName(branch), &repo.Runner.Config.FullConfig)
		}
	}
	return parkConfig{
		allBranches:    branchesSnapshot.Branches,
		branchesToPark: branchesToPark,
	}, nil
}

func validateParkConfig(config parkConfig) error {
	for branchName, branchType := range config.branchesToPark {
		if !config.allBranches.HasLocalBranch(branchName) {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotPark)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotPark)
		case configdomain.BranchTypeParkedBranch:
			return fmt.Errorf(messages.BranchIsAlreadyParked, branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch:
		}
	}
	return nil
}
