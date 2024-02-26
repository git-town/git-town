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
	for _, branchToPark := range config.branchesToPark {
		if !config.branches.HasLocalBranch(branchToPark) {
			return fmt.Errorf(messages.BranchDoesntExist, branchToPark)
		}
		if err = validateIsParkableBranch(branchToPark, &repo.Runner.Config.FullConfig); err != nil {
			return err
		}
	}
	if err = repo.Runner.Config.AddToParkedBranches(config.branchesToPark...); err != nil {
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
	branches       gitdomain.BranchInfos
	branchesToPark gitdomain.LocalBranchNames
}

func determineParkConfig(args []string, repo *execute.OpenRepoResult) (parkConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return parkConfig{}, err
	}
	var branchesToPark gitdomain.LocalBranchNames
	if len(args) == 0 {
		branchesToPark = gitdomain.LocalBranchNames{branchesSnapshot.Active}
	} else {
		branchesToPark = gitdomain.NewLocalBranchNames(args...)
	}
	return parkConfig{
		branches:       branchesSnapshot.Branches,
		branchesToPark: branchesToPark,
	}, nil
}

func validateIsParkableBranch(branch gitdomain.LocalBranchName, config *configdomain.FullConfig) error {
	if config.IsMainBranch(branch) {
		return errors.New(messages.MainBranchCannotPark)
	}
	if config.IsPerennialBranch(branch) {
		return errors.New(messages.PerennialBranchCannotPark)
	}
	return nil
}
