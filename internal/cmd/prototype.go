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
	interpreterConfig "github.com/git-town/git-town/v16/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const prototypeDesc = "Make an existing branch a prototype branch"

const prototypeHelp = `
A prototype branch is for local-only development. It incorporates updates from its parent branch and is not pushed to the remote repository until you run "git town propose" on it.

You can create new prototype branches using git town hack, append, or prepend with the --prototype option.
`

func prototypeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "prototype [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: "types",
		Short:   prototypeDesc,
		Long:    cmdhelpers.Long(prototypeDesc, prototypeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePrototype(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrototype(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := determinePrototypeData(args, repo)
	if err != nil {
		return err
	}
	if err = validatePrototypeData(data); err != nil {
		return err
	}
	branchNames := data.branchesToPrototype.Keys()
	if err = repo.UnvalidatedConfig.NormalConfig.AddToPrototypeBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonPrototypeBranchTypes(data.branchesToPrototype, repo.UnvalidatedConfig); err != nil {
		return err
	}
	if checkout, hasCheckout := data.checkout.Get(); hasCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, checkout, false); err != nil {
			return err
		}
	}
	printPrototypeBranches(branchNames)
	return interpreterConfig.Finished(interpreterConfig.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.branchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "prototype",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               verbose,
	})
}

type prototypeData struct {
	branchInfos         gitdomain.BranchInfos
	branchesSnapshot    gitdomain.BranchesSnapshot
	branchesToPrototype configdomain.BranchesAndTypes
	checkout            Option[gitdomain.LocalBranchName]
}

func printPrototypeBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.PrototypeBranchIsNowPrototype, branch)
	}
}

func removeNonPrototypeBranchTypes(branches configdomain.BranchesAndTypes, config config.UnvalidatedConfig) error {
	for branchName, branchType := range branches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			if err := config.NormalConfig.RemoveFromContributionBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeObservedBranch:
			if err := config.NormalConfig.RemoveFromObservedBranches(branchName); err != nil {
				return err
			}
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypePrototypeBranch,
			configdomain.BranchTypeMainBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determinePrototypeData(args []string, repo execute.OpenRepoResult) (prototypeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return prototypeData{}, err
	}
	branchesToPrototype, branchToCheckout, err := execute.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig.UnvalidatedConfig)
	return prototypeData{
		branchInfos:         branchesSnapshot.Branches,
		branchesSnapshot:    branchesSnapshot,
		branchesToPrototype: branchesToPrototype,
		checkout:            branchToCheckout,
	}, err
}

func validatePrototypeData(data prototypeData) error {
	for branchName, branchType := range data.branchesToPrototype {
		if !data.branchInfos.HasLocalBranch(branchName) && !data.branchInfos.HasMatchingTrackingBranchFor(branchName) {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotPrototype)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotPrototype)
		case configdomain.BranchTypePrototypeBranch:
			return fmt.Errorf(messages.BranchIsAlreadyPrototype, branchName)
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypeObservedBranch:
		default:
			panic("unhandled branch type" + branchType.String())
		}
	}
	return nil
}
