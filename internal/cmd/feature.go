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
	featureDesc = "Convert branches to feature branches"
	featureHelp = `
Marks the given local branches as feature branches.
If no branch is provided, makes the current branch a feature branch.

Feauture branches are branches that you own and use to make code changes in.
`
)

func featureCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "feature [branches]",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDTypes,
		Short:   featureDesc,
		Long:    cmdhelpers.Long(featureDesc, featureHelp),
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
			return executeFeature(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeFeature(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, err := determineFeatureData(args, repo)
	if err != nil {
		return err
	}
	if err = validateFeatureData(data, repo); err != nil {
		return err
	}
	branchNames := data.branchesToFeature.Keys()
	if err = gitconfig.SetBranchTypeOverride(repo.Backend, configdomain.BranchTypeFeatureBranch, branchNames...); err != nil {
		return err
	}
	printFeatureBranches(branchNames)
	if checkout, hasCheckout := data.checkout.Get(); hasCheckout {
		if err = repo.Git.CheckoutBranch(repo.Frontend, checkout, false); err != nil {
			return err
		}
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: Some(data.branchesSnapshot),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "feature",
		CommandsCounter:       repo.CommandsCounter,
		ConfigDir:             repo.ConfigDir,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		TouchedBranches:       branchNames.BranchNames(),
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
	})
}

type featureData struct {
	branchInfos       gitdomain.BranchInfos
	branchesSnapshot  gitdomain.BranchesSnapshot
	branchesToFeature configdomain.BranchesAndTypes
	checkout          Option[gitdomain.LocalBranchName]
}

func printFeatureBranches(branches gitdomain.LocalBranchNames) {
	for _, branch := range branches {
		fmt.Printf(messages.BranchIsNowFeature, branch)
	}
}

func determineFeatureData(args []string, repo execute.OpenRepoResult) (featureData, error) {
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return featureData{}, err
	}
	if branchesSnapshot.DetachedHead {
		return featureData{}, errors.New(messages.FeatureDetachedHead)
	}
	branchesToFeature, err := config.BranchesToMark(args, branchesSnapshot, repo.UnvalidatedConfig)
	return featureData{
		branchInfos:       branchesSnapshot.Branches,
		branchesSnapshot:  branchesSnapshot,
		branchesToFeature: branchesToFeature.BranchesToMark,
		checkout:          branchesToFeature.BranchToCheckout,
	}, err
}

func validateFeatureData(data featureData, repo execute.OpenRepoResult) error {
	for branchName, branchType := range mapstools.SortedKeyValues(data.branchesToFeature) {
		switch branchType {
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.MainBranchCannotMakeFeature)
		case configdomain.BranchTypePerennialBranch:
			return errors.New(messages.PerennialBranchCannotMakeFeature)
		case configdomain.BranchTypeFeatureBranch:
			repo.FinalMessages.Addf(messages.HackBranchIsAlreadyFeature, branchName)
		case
			configdomain.BranchTypeObservedBranch,
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
			hasLocalBranch := data.branchesSnapshot.Branches.HasLocalBranch(branchName)
			hasRemoteBranch := data.branchesSnapshot.Branches.HasMatchingTrackingBranchFor(branchName)
			if !hasLocalBranch && !hasRemoteBranch {
				return fmt.Errorf(messages.BranchDoesntExist, branchName)
			}
		}
	}
	return nil
}
