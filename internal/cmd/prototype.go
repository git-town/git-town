package cmd

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/mapstools"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	prototypeDesc = "Make an existing branch a prototype branch"
	prototypeHelp = `
A prototype branch is for local-only development.
It incorporates updates from its parent branch
and is not pushed to the remote repository
until you run "git town propose" on it.

You can create new prototype branches
using git town hack, append, or prepend
with the --prototype option.
`
)

func prototypeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "prototype [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDTypes,
		Short:   prototypeDesc,
		Long:    cmdhelpers.Long(prototypeDesc, prototypeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executePrototype(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrototype(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, err := determinePrototypeData(args, repo)
	if err != nil {
		return err
	}
	if err = validatePrototypeData(data, repo); err != nil {
		return err
	}
	branchNames := data.branchesToPrototype.Keys()
	if err = gitconfig.SetBranchTypeOverride(repo.Backend, configdomain.BranchTypePrototypeBranch, branchNames...); err != nil {
		return err
	}
	if checkout, hasCheckout := data.checkout.Get(); hasCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, checkout, false); err != nil {
			return err
		}
	}
	printPrototypeBranches(branchNames)
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.branchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "prototype",
		CommandsCounter:       repo.CommandsCounter,
		ConfigDir:             repo.ConfigDir,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
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
		fmt.Printf(messages.BranchIsNowPrototype, branch)
	}
}

func determinePrototypeData(args []string, repo execute.OpenRepoResult) (prototypeData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return prototypeData{}, err
	}
	if branchesSnapshot.DetachedHead {
		return prototypeData{}, errors.New(messages.PrototypeDetachedHead)
	}
	branchesToPrototype, err := config.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	return prototypeData{
		branchInfos:         branchesSnapshot.Branches,
		branchesSnapshot:    branchesSnapshot,
		branchesToPrototype: branchesToPrototype.BranchesToMark,
		checkout:            branchesToPrototype.BranchToCheckout,
	}, err
}

func validatePrototypeData(data prototypeData, repo execute.OpenRepoResult) error {
	for branchName, branchType := range mapstools.SortedKeyValues(data.branchesToPrototype) {
		if !data.branchesSnapshot.Branches.HasLocalBranch(branchName) && !data.branchesSnapshot.Branches.HasMatchingTrackingBranchFor(branchName) {
			return fmt.Errorf(messages.BranchDoesntExist, branchName)
		}
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotPrototype)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotPrototype)
		case configdomain.BranchTypePrototypeBranch:
			repo.FinalMessages.Addf(messages.BranchIsAlreadyPrototype, branchName)
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypeObservedBranch:
			data.branchesToPrototype.Add(branchName, branchType)
		}
	}
	return nil
}
