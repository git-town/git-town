package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/commandconfig"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const prototypeDesc = "Create a new prototype branch"

const prototypeHelp = `
Marks the given local branches as prototype.
If no branch is provided, marks the current branch.

A prototype branch is a local-only feature branch that incorporates updates from its parent branch but is not pushed to the remote repository.
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

func executePrototype(args []string, verbose bool) error {
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
	data, err := determinePrototypeData(args, repo)
	if err != nil {
		return err
	}
	err = validatePrototypeData(data)
	if err != nil {
		return err
	}
	branchNames := data.branchesToPrototype.Keys()
	if err = repo.UnvalidatedConfig.AddToPrototypeBranches(branchNames...); err != nil {
		return err
	}
	if err = removeNonPrototypeBranchTypes(data.branchesToPrototype, repo.UnvalidatedConfig); err != nil {
		return err
	}
	printPrototypeBranches(branchNames)
	if checkout, hasCheckout := data.checkout.Get(); hasCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, checkout, false); err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             repo.Backend,
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "prototype",
		CommandsCounter:     repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       repo.FinalMessages,
		RootDir:             repo.RootDir,
		Verbose:             verbose,
	})
}

type prototypeData struct {
	allBranches         gitdomain.BranchInfos
	branchesToPrototype commandconfig.BranchesAndTypes
	checkout            Option[gitdomain.LocalBranchName]
}

func printPrototypeBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.PrototypeBranchIsNowPrototype, branch)
	}
}

func removeNonPrototypeBranchTypes(branches map[gitdomain.LocalBranchName]configdomain.BranchType, config config.UnvalidatedConfig) error {
	for branchName, branchType := range branches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			if err := config.RemoveFromContributionBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeParkedBranch:
			if err := config.RemoveFromParkedBranches(branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypePrototypeBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		}
	}
	return nil
}

func determinePrototypeData(args []string, repo execute.OpenRepoResult) (prototypeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return prototypeData{}, err
	}
	branchesToPrototype := commandconfig.BranchesAndTypes{}
	checkout := None[gitdomain.LocalBranchName]()
	currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		return prototypeData{}, errors.New(messages.CurrentBranchCannotDetermine)
	}
	switch len(args) {
	case 0:
		branchesToPrototype.Add(currentBranch, *repo.UnvalidatedConfig.Config)
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToPrototype.Add(branch, *repo.UnvalidatedConfig.Config)
		trackingBranchName := branch.TrackingBranch()
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByRemoteName(trackingBranchName).Get()
		if !hasBranchInfo {
			return prototypeData{}, fmt.Errorf(messages.BranchDoesntExist, branch.String())
		}
		if branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			checkout = Some(branch)
		}
	default:
		branchesToPrototype.AddMany(gitdomain.NewLocalBranchNames(args...), *repo.UnvalidatedConfig.Config)
	}
	return prototypeData{
		allBranches:         branchesSnapshot.Branches,
		branchesToPrototype: branchesToPrototype,
		checkout:            checkout,
	}, nil
}

func validatePrototypeData(data prototypeData) error {
	for branchName, branchType := range data.branchesToPrototype {
		if !data.allBranches.HasLocalBranch(branchName) && !data.allBranches.HasMatchingTrackingBranchFor(branchName) {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotPrototype)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotPrototype)
		case configdomain.BranchTypePrototypeBranch:
			return fmt.Errorf(messages.BranchIsAlreadyPrototype, branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeContributionBranch, configdomain.BranchTypeParkedBranch:
		}
	}
	return nil
}
