package validate

import (
	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
)

func Config(args ConfigArgs) (config.ValidatedConfig, bool, error) {
	// check Git user data
	gitUserEmail, gitUserName, err := GitUser(*args.Unvalidated.Config)
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}

	// enter and save main and perennials
	mainBranch, hasMain := args.Unvalidated.Config.MainBranch.Get()
	if !hasMain {
		validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
			Backend:               args.Backend,
			DialogInputs:          args.TestInputs,
			GetDefaultBranch:      args.Git.DefaultBranch,
			HasConfigFile:         args.Unvalidated.ConfigFile.IsSome(),
			LocalBranches:         args.LocalBranches,
			UnvalidatedMain:       args.Unvalidated.Config.MainBranch,
			UnvalidatedPerennials: args.Unvalidated.Config.PerennialBranches,
		})
		if err != nil || aborted {
			return config.EmptyValidatedConfig(), aborted, err
		}
		mainBranch = validatedMain
		if err = args.Unvalidated.SetMainBranch(validatedMain); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		if len(additionalPerennials) > 0 {
			newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
			if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
				return config.EmptyValidatedConfig(), false, err
			}
		}
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, exit, err := dialog.Lineage(dialog.LineageArgs{
		BranchesToVerify: args.BranchesToValidate,
		Config:           *args.Unvalidated.Config,
		DefaultChoice:    mainBranch,
		DialogTestInputs: args.TestInputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       mainBranch,
	})
	if err != nil || exit {
		return config.EmptyValidatedConfig(), exit, err
	}
	for branch, parent := range additionalLineage {
		if err = args.Unvalidated.SetParent(branch, parent); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
	}

	// create validated configuration
	validatedConfig := config.ValidatedConfig{
		Config: configdomain.ValidatedConfig{
			UnvalidatedConfig: args.Unvalidated.Config,
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
