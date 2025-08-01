package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	parkDesc = "Suspend syncing of some feature branches"
	parkHelp = `
Parks the given local feature branches.
If no branch is provided, parks the current branch.

Git Town does not sync parked branches.
The currently checked out branch gets synced even if parked.
`
)

func parkCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "park [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDTypes,
		Short:   parkDesc,
		Long:    cmdhelpers.Long(parkDesc, parkHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				DryRun:      None[configdomain.DryRun](),
				Verbose:     verbose,
			})
			return executePark(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePark(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, err := determineParkData(args, repo)
	if err != nil {
		return err
	}
	if err = validateParkData(data, repo); err != nil {
		return err
	}
	branchNames := data.branchesToPark.Keys()
	if err = gitconfig.SetBranchTypeOverride(repo.Backend, configdomain.BranchTypeParkedBranch, branchNames...); err != nil {
		return err
	}
	printParkedBranches(branchNames)
	if branchToCheckout, hasBranchToCheckout := data.branchToCheckout.Get(); hasBranchToCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, false); err != nil {
			return err
		}
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.beginBranchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "park",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
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
		fmt.Printf(messages.BranchIsNowParked, branch)
	}
}

func determineParkData(args []string, repo execute.OpenRepoResult) (parkData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return parkData{}, err
	}
	branchesToPark, branchToCheckout, err := config.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	return parkData{
		beginBranchesSnapshot: branchesSnapshot,
		branchInfos:           branchesSnapshot.Branches,
		branchToCheckout:      branchToCheckout,
		branchesToPark:        branchesToPark,
	}, err
}

func validateParkData(data parkData, repo execute.OpenRepoResult) error {
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
		hasRemoteBranch := data.branchInfos.HasMatchingTrackingBranchFor(branchName, repo.UnvalidatedConfig.NormalConfig.DevRemote)
		if !hasLocalBranch && !hasRemoteBranch {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
	}
	return nil
}
