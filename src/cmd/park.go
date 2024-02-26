package cmd

import (
	"errors"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const parkDesc = "Parks the current branch"

const parkHelp = `
Git Town does not sync a parked branch unless it is checked when the sync starts.`

func parkCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "park",
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
	config, err := determineParkConfig(args, repo.Runner)
	for _, branchToPark := range config.branchesToPark {
		err = validateIsParkableBranch(branchToPark, &repo.Runner.Config.FullConfig)
		if err != nil {
			return err
		}
	}
	for _, branchToPark := range config.branchesToPark {
		err = repo.Runner.Config.AddToParkedBranches(branchToPark)
		if err != nil {
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
	currentBranch  gitdomain.LocalBranchName
}

func determineParkConfig(args []string, runner *git.ProdRunner) (parkConfig, error) {
	currentBranch, err := runner.Backend.CurrentBranch()
	if err != nil {
		return parkConfig{}, err
	}
	branchesToPark := make(gitdomain.LocalBranchNames, len(args))
	for b, branchName := range args {
		branchesToPark[b] = gitdomain.NewLocalBranchName(branchName)
	}
	return parkConfig{
		branchesToPark: branchesToPark,
		currentBranch:  currentBranch,
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
