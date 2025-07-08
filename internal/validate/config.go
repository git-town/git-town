package validate

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configsetup"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
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
		// run setup assistant here
		setupData :=configsetup.SetupData{
			Backend:           args.Backend,
			ConfigSnapshot:    args.ConfigSnapshot,
			CommandName:       "",
			CommandsCounter:   Mutable[gohacks.Counter]{},
			ConfigFile:        Option[configdomain.PartialConfig]{},
			DialogInputs:      components.TestInputs{},
			FinalMessages:     stringslice.Collector{},
			Frontend:          nil,
			Git:               git.Commands{},
			LocalBranches:     gitdomain.BranchInfos{},
			Remotes:           gitdomain.Remotes{},
			RootDir:           "",
			UnvalidatedConfig: config.UnvalidatedConfig{},
			UserInput:         configsetup.UserInput{},
			Verbose:           false,
		}
		mainBranch, err = configsetup.Execute(&setupData)
		args.Unvalidated =
		args.BranchesAndTypes[validatedMain] = configdomain.BranchTypeMainBranch
		if err = args.Unvalidated.Value.SetMainBranch(validatedMain, args.Backend); err != nil {
			return config.EmptyValidatedConfig(), false, err
		}
		if len(additionalPerennials) > 0 {
			newPerennials := append(args.Unvalidated.Value.NormalConfig.PerennialBranches, additionalPerennials...)
			if err = args.Unvalidated.Value.NormalConfig.SetPerennialBranches(args.Backend, newPerennials); err != nil {
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
	DialogTestInputs   dialogcomponents.TestInputs
	Frontend           subshelldomain.Runner
	Git                git.Commands
	LocalBranches      gitdomain.LocalBranchNames
	RepoStatus         gitdomain.RepoStatus
	TestInputs         dialogcomponents.TestInputs
	Unvalidated        Mutable[config.UnvalidatedConfig]
}
