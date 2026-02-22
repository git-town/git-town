package validate

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/setup"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func Config(args ConfigArgs) (config.ValidatedConfig, dialogdomain.Exit, error) {
	// enter and save main and perennials
	mainBranch, hasMain := args.Unvalidated.Value.UnvalidatedConfig.MainBranch.Get()
	var userInput setup.UserInput
	if !hasMain {
		setupData := setup.Data{
			Backend:       args.Backend,
			Config:        args.Unvalidated.Immutable(),
			Git:           args.Git,
			Inputs:        args.Inputs,
			LocalBranches: args.LocalBranches,
			Remotes:       args.Remotes,
			Snapshot:      args.ConfigSnapshot,
		}
		var exit dialogdomain.Exit
		userInput, exit, enterAll, err := setup.Enter(setupData, args.ConfigDir)
		if err != nil {
			if errors.Is(err, dialogcomponents.ErrNoTTY) {
				return config.EmptyValidatedConfig(), false, errors.New(messages.NoTTYMainBranchMissing) //lint:ignore ST1005 This error contains user-visible guidance, and therefore needs to end with a period.
			}
			return config.EmptyValidatedConfig(), exit, err
		}
		if exit {
			return config.EmptyValidatedConfig(), exit, nil
		}
		err = setup.Save(userInput, args.Unvalidated.Immutable(), setupData, enterAll, args.Frontend)
		if err != nil {
			return config.EmptyValidatedConfig(), exit, err
		}
		mainBranch = userInput.ValidatedConfig.MainBranch
	}
	// enter and save missing parent branches
	lineageResult, exit, err := dialog.Lineage(dialog.LineageArgs{
		BranchInfos:      args.BranchInfos,
		BranchesAndTypes: args.BranchesAndTypes,
		BranchesToVerify: args.BranchesToValidate,
		Config:           args.Unvalidated.Immutable(),
		Connector:        args.Connector,
		DefaultChoice:    mainBranch,
		Inputs:           args.Inputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       mainBranch,
	})
	if err != nil || exit {
		return config.EmptyValidatedConfig(), exit, err
	}
	for _, entry := range lineageResult.AdditionalLineage.Entries() {
		if err = args.Unvalidated.Value.NormalConfig.SetParent(args.Backend, entry.Child, entry.Parent); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}
	if len(lineageResult.AdditionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Value.NormalConfig.PerennialBranches, lineageResult.AdditionalPerennials...)
		if err = args.Unvalidated.Value.NormalConfig.SetPerennialBranches(args.Backend, newPerennials); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}

	// store the entered data
	if !hasMain {
		args.Unvalidated.Value.NormalConfig = args.Unvalidated.Value.NormalConfig.OverwriteWith(userInput.Data)
		args.Unvalidated.Value.UnvalidatedConfig.MainBranch = Some(mainBranch)
		args.BranchesAndTypes[mainBranch] = configdomain.BranchTypeMainBranch
	}
	validatedConfig := config.ValidatedConfig{
		ValidatedConfigData: configdomain.ValidatedConfigData{
			MainBranch: mainBranch,
		},
		NormalConfig: args.Unvalidated.Value.NormalConfig,
	}
	return validatedConfig, false, err
}

type ConfigArgs struct {
	Backend            subshelldomain.RunnerQuerier
	BranchInfos        gitdomain.BranchInfos
	BranchesAndTypes   configdomain.BranchesAndTypes
	BranchesToValidate gitdomain.LocalBranchNames
	ConfigDir          configdomain.RepoConfigDir
	ConfigSnapshot     configdomain.BeginConfigSnapshot
	Connector          Option[forgedomain.Connector]
	Frontend           subshelldomain.Runner
	Git                git.Commands
	Inputs             dialogcomponents.Inputs
	LocalBranches      gitdomain.LocalBranchNames
	Remotes            gitdomain.Remotes
	RepoStatus         gitdomain.RepoStatus
	Unvalidated        Mutable[config.UnvalidatedConfig]
}
