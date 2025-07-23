package validate

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/setup"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func Config(args ConfigArgs) (config.ValidatedConfig, dialogdomain.Exit, error) {
	// check Git user data
	gitUserEmail, gitUserName, err := GitUser(args.Unvalidated.Value.UnvalidatedConfig)
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}

	// enter and save main and perennials
	mainBranch, hasMain := args.Unvalidated.Value.UnvalidatedConfig.MainBranch.Get()
	if !hasMain {
		setupData := setup.Data{
			Backend:        args.Backend,
			Config:         args.Unvalidated.Immutable(),
			ConfigSnapshot: args.ConfigSnapshot,
			Git:            args.Git,
			Inputs:         args.Inputs,
			LocalBranches:  args.LocalBranches,
			Remotes:        args.Remotes,
		}
		userInput, exit, err := setup.Enter(setupData)
		if err != nil || exit {
			return config.EmptyValidatedConfig(), exit, err
		}
		setup.Save(userInput, args.Unvalidated.Immutable(), setupData, args.Frontend)
		mainBranch = userInput.ValidatedConfig.MainBranch
		args.BranchesAndTypes[mainBranch] = configdomain.BranchTypeMainBranch // TODO: don't modify this here?
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, exit, err := dialog.Lineage(dialog.LineageArgs{
		BranchesAndTypes:  args.BranchesAndTypes,
		BranchesToVerify:  args.BranchesToValidate,
		Connector:         args.Connector,
		DefaultChoice:     mainBranch,
		Inputs:            args.Inputs,
		Lineage:           args.Unvalidated.Value.NormalConfig.Lineage,
		LocalBranches:     args.LocalBranches,
		MainBranch:        mainBranch,
		PerennialBranches: args.Unvalidated.Value.NormalConfig.PerennialBranches,
	})
	if err != nil || exit {
		return config.EmptyValidatedConfig(), exit, err
	}
	for _, entry := range additionalLineage.Entries() {
		if err = args.Unvalidated.Value.NormalConfig.SetParent(args.Backend, entry.Child, entry.Parent); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Value.NormalConfig.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.Value.NormalConfig.SetPerennialBranches(args.Backend, newPerennials); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}

	// create validated configuration
	validatedConfig := config.ValidatedConfig{
		ValidatedConfigData: configdomain.ValidatedConfigData{
			GitUserEmail: gitUserEmail,
			GitUserName:  gitUserName,
			MainBranch:   mainBranch,
		},
		// TODO: merge userInput.Data into this
		NormalConfig: args.Unvalidated.Value.NormalConfig,
	}
	return validatedConfig, false, err
}

type ConfigArgs struct {
	Backend            subshelldomain.RunnerQuerier
	BranchesAndTypes   configdomain.BranchesAndTypes
	BranchesSnapshot   gitdomain.BranchesSnapshot
	BranchesToValidate gitdomain.LocalBranchNames
	ConfigSnapshot     undoconfig.ConfigSnapshot
	Connector          Option[forgedomain.Connector]
	Frontend           subshelldomain.Runner
	Git                git.Commands
	Inputs             dialogcomponents.Inputs
	LocalBranches      gitdomain.LocalBranchNames
	Remotes            gitdomain.Remotes
	RepoStatus         gitdomain.RepoStatus
	Unvalidated        Mutable[config.UnvalidatedConfig]
}
