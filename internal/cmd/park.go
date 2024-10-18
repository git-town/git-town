package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	configInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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

func executePark(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := determineParkData(args, repo)
	if err != nil {
		return err
	}
	if err = validateParkData(data); err != nil {
		return err
	}
	branchNames := data.branchesToPark.Keys()
	if err = repo.UnvalidatedConfig.AddToParkedBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonParkBranchTypes(data.branchesToPark, repo.UnvalidatedConfig); err != nil {
		return err
	}
	printParkedBranches(branchNames)
	if branchToCheckout, hasBranchToCheckout := data.branchToCheckout.Get(); hasBranchToCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, false); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.beginBranchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "park",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               verbose,
	})
}

type parkData struct {
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	branchInfos           gitdomain.BranchInfos
	branchToCheckout      Option[gitdomain.LocalBranchName]
	branchesToPark        configdomain.BranchesAndTypes
}

func printParkedBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.ParkedBranchIsNowParked, branch)
	}
}

func removeNonParkBranchTypes(branches map[gitdomain.LocalBranchName]configdomain.BranchType, config config.UnvalidatedConfig) error {
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
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypeMainBranch,
			configdomain.BranchTypePerennialBranch,
			configdomain.BranchTypePrototypeBranch:
		}
	}
	return nil
}

func determineParkData(args []string, repo execute.OpenRepoResult) (parkData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return parkData{}, err
	}
	branchesToPark, branchToCheckout, err := execute.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig.UnvalidatedConfig.Get())
	return parkData{
		beginBranchesSnapshot: branchesSnapshot,
		branchInfos:           branchesSnapshot.Branches,
		branchToCheckout:      branchToCheckout,
		branchesToPark:        branchesToPark,
	}, err
}

func validateParkData(data parkData) error {
	for branchName, branchType := range data.branchesToPark {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotPark)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotPark)
		case configdomain.BranchTypeParkedBranch:
			return fmt.Errorf(messages.BranchIsAlreadyParked, branchName)
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeObservedBranch,
			configdomain.BranchTypePrototypeBranch:
		}
		hasLocalBranch := data.branchInfos.HasLocalBranch(branchName)
		hasRemoteBranch := data.branchInfos.HasMatchingTrackingBranchFor(branchName)
		if !hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
	}
	return nil
}
