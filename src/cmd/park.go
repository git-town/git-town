package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/commandconfig"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const parkDesc = "Suspend syncing of some feature branches"

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
	data, exit, err := determineParkData(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	err = validateParkData(data)
	if err != nil {
		return err
	}
	branchNames := data.branchesToPark.Keys()
	if err = data.config.AddToParkedBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonParkBranchTypes(data.branchesToPark, data.config); err != nil {
		return err
	}
	printParkedBranches(branchNames)
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             repo.Backend,
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "park",
		CommandsCounter:     repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       repo.FinalMessages,
		RootDir:             repo.RootDir,
		Verbose:             verbose,
	})
}

type parkData struct {
	allBranches    gitdomain.BranchInfos
	branchesToPark commandconfig.BranchesAndTypes
	config         config.ValidatedConfig
}

func printParkedBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ParkedBranchIsNowParked, branch)
	}
}

func removeNonParkBranchTypes(branches map[gitdomain.LocalBranchName]configdomain.BranchType, config config.ValidatedConfig) error {
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

func determineParkData(args []string, repo execute.OpenRepoResult, verbose bool) (parkData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return parkData{}, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return parkData{}, exit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesToPark := commandconfig.BranchesAndTypes{}
	if len(args) == 0 {
		branchesToPark.Add(branchesSnapshot.Active, *repo.UnvalidatedConfig.Config)
	} else {
		branchesToPark.AddMany(gitdomain.NewLocalBranchNames(args...), *repo.UnvalidatedConfig.Config)
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return parkData{}, exit, err
	}
	return parkData{
		allBranches:    branchesSnapshot.Branches,
		branchesToPark: branchesToPark,
		config:         validatedConfig,
	}, false, nil
}

func validateParkData(data parkData) error {
	for branchName, branchType := range data.branchesToPark {
		if !data.allBranches.HasLocalBranch(branchName) {
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
