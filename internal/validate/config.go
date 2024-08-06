package validate

import (
	"github.com/git-town/git-town/v14/internal/cli/dialog"
	"github.com/git-town/git-town/v14/internal/cli/dialog/components"
	"github.com/git-town/git-town/v14/internal/config"
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/git"
	"github.com/git-town/git-town/v14/internal/git/gitdomain"
)

func Config(args ConfigArgs) (config.ValidatedConfig, bool, error) {
	// check Git user data
	gitUserEmail, gitUserName, err := GitUser(args.Unvalidated.Config.Get())
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}

	// enter and save main and perennials
	mainBranch, hasMain := args.Unvalidated.Config.Value.MainBranch.Get()
	if !hasMain {
		validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
			Backend:               args.Backend,
			DialogInputs:          args.TestInputs,
			GetDefaultBranch:      args.Git.DefaultBranch,
			HasConfigFile:         args.Unvalidated.ConfigFile.IsSome(),
			LocalBranches:         args.LocalBranches,
			UnvalidatedMain:       args.Unvalidated.Config.Value.MainBranch,
			UnvalidatedPerennials: args.Unvalidated.Config.Value.PerennialBranches,
		})
		if err != nil || aborted {
			return config.EmptyValidatedConfig(), aborted, err
		}
		mainBranch = validatedMain
		if err = args.Unvalidated.SetMainBranch(validatedMain); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		if len(additionalPerennials) > 0 {
			newPerennials := append(args.Unvalidated.Config.Value.PerennialBranches, additionalPerennials...)
			if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
				return config.EmptyValidatedConfig(), false, err
			}
		}
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, exit, err := dialog.Lineage(dialog.LineageArgs{
		BranchesToVerify: args.BranchesToValidate,
		Config:           args.Unvalidated.Config.Get(),
		DefaultChoice:    mainBranch,
		DialogTestInputs: args.TestInputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       mainBranch,
	})
	if err != nil || exit {
		return config.EmptyValidatedConfig(), exit, err
	}
	for _, entry := range additionalLineage.Entries() {
		if err = args.Unvalidated.SetParent(entry.Child, entry.Parent); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.Value.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}

	// create validated configuration
	validatedConfig := config.ValidatedConfig{
		Config: configdomain.ValidatedConfig{
			UnvalidatedConfig: args.Unvalidated.Config.Value,
			GitUserEmail:      gitUserEmail,
			GitUserName:       gitUserName,
			MainBranch:        mainBranch,
		},
		UnvalidatedConfig: &args.Unvalidated,
	}

	return validatedConfig, false, err
}

type ConfigArgs struct {
	Backend            gitdomain.RunnerQuerier
	BranchesSnapshot   gitdomain.BranchesSnapshot
	BranchesToValidate gitdomain.LocalBranchNames
	DialogTestInputs   components.TestInputs
	Frontend           gitdomain.Runner
	Git                git.Commands
	LocalBranches      gitdomain.LocalBranchNames
	RepoStatus         gitdomain.RepoStatus
	TestInputs         components.TestInputs
	Unvalidated        config.UnvalidatedConfig
}
