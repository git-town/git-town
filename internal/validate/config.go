package validate

import (
	"github.com/git-town/git-town/v17/internal/cli/dialog"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/hosting/hostingdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

func Config(args ConfigArgs) (config.ValidatedConfig, bool, error) {
	// check Git user data
	gitUserEmail, gitUserName, err := GitUser(args.Unvalidated.Value.UnvalidatedConfig)
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}

	// enter and save main and perennials
	mainBranch, hasMain := args.Unvalidated.Value.UnvalidatedConfig.MainBranch.Get()
	if !hasMain {
		validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
			Backend:               args.Backend,
			DialogInputs:          args.TestInputs,
			GetDefaultBranch:      args.Git.DefaultBranch,
			HasConfigFile:         args.Unvalidated.Value.NormalConfig.ConfigFile.IsSome(),
			LocalBranches:         args.LocalBranches,
			UnvalidatedMain:       args.Unvalidated.Value.UnvalidatedConfig.MainBranch,
			UnvalidatedPerennials: args.Unvalidated.Value.NormalConfig.PerennialBranches,
		})
		if err != nil || aborted {
			return config.EmptyValidatedConfig(), aborted, err
		}
		mainBranch = validatedMain
		args.BranchesAndTypes[validatedMain] = configdomain.BranchTypeMainBranch
		if err = args.Unvalidated.Value.SetMainBranch(validatedMain); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		if len(additionalPerennials) > 0 {
			newPerennials := append(args.Unvalidated.Value.NormalConfig.PerennialBranches, additionalPerennials...)
			if err = args.Unvalidated.Value.NormalConfig.SetPerennialBranches(newPerennials); err != nil {
				return config.EmptyValidatedConfig(), false, err
			}
		}
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, exit, err := dialog.Lineage(dialog.LineageArgs{
		BranchesAndTypes:  args.BranchesAndTypes,
		BranchesToVerify:  args.BranchesToValidate,
		Connector:         args.Connector,
		DefaultChoice:     mainBranch,
		DialogTestInputs:  args.TestInputs,
		Lineage:           args.Unvalidated.Value.NormalConfig.Lineage,
		LocalBranches:     args.LocalBranches,
		MainBranch:        mainBranch,
		PerennialBranches: args.Unvalidated.Value.NormalConfig.PerennialBranches,
	})
	if err != nil || exit {
		return config.EmptyValidatedConfig(), exit, err
	}
	for _, entry := range additionalLineage.Entries() {
		if err = args.Unvalidated.Value.NormalConfig.SetParent(entry.Child, entry.Parent); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Value.NormalConfig.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.Value.NormalConfig.SetPerennialBranches(newPerennials); err != nil {
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
		NormalConfig: args.Unvalidated.Value.NormalConfig,
	}

	return validatedConfig, false, err
}

type ConfigArgs struct {
	Backend            gitdomain.RunnerQuerier
	BranchesAndTypes   configdomain.BranchesAndTypes
	BranchesSnapshot   gitdomain.BranchesSnapshot
	BranchesToValidate gitdomain.LocalBranchNames
	Connector          Option[hostingdomain.Connector]
	DialogTestInputs   components.TestInputs
	Frontend           gitdomain.Runner
	Git                git.Commands
	LocalBranches      gitdomain.LocalBranchNames
	RepoStatus         gitdomain.RepoStatus
	TestInputs         components.TestInputs
	Unvalidated        Mutable[config.UnvalidatedConfig]
}
