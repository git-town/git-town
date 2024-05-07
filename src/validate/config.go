package validate

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
)

func Config(args ConfigArgs) (config.ValidatedConfig, bool, error) {
	// check Git user data
	gitUserEmail, hasGitUserEmail := args.Unvalidated.Config.GitUserEmail.Get()
	if !hasGitUserEmail {
		return config.EmptyValidatedConfig(), false, errors.New(messages.GitUserEmailMissing)
	}
	gitUserName, hasGitUserName := args.Unvalidated.Config.GitUserName.Get()
	if !hasGitUserName {
		return config.EmptyValidatedConfig(), false, errors.New(messages.GitUserNameMissing)
	}

	// enter and save main and perennials
	mainBranch, hasMain := args.Unvalidated.Config.MainBranch.Get()
	if !hasMain {
		validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
			DialogInputs:          args.TestInputs,
			GetDefaultBranch:      args.Backend.DefaultBranch,
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

	// remove outdated lineage
	err = args.Unvalidated.RemoveOutdatedConfiguration(args.LocalBranches)
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
	}
	err = cleanupPerennialParentEntries(args.Unvalidated.Config.Lineage, args.Unvalidated.Config.PerennialBranches, args.Unvalidated.GitConfig, args.FinalMessages)
	if err != nil {
		return config.EmptyValidatedConfig(), false, err
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

	// handle unfinished state
	if args.HandleUnfinishedState {
		exit, err = HandleUnfinishedState(UnfinishedStateArgs{
			Backend:                 args.Backend,
			CommandsCounter:         args.CommandsCounter,
			Config:                  validatedConfig,
			Connector:               nil,
			CurrentBranch:           args.BranchesSnapshot.Active,
			DialogTestInputs:        args.DialogTestInputs,
			FinalMessages:           args.FinalMessages,
			Frontend:                args.Frontend,
			HasOpenChanges:          args.RepoStatus.OpenChanges,
			InitialBranchesSnapshot: args.BranchesSnapshot,
			InitialConfigSnapshot:   args.ConfigSnapshot,
			InitialStashSize:        args.StashSize,
			Lineage:                 validatedConfig.Config.Lineage,
			PushHook:                validatedConfig.Config.PushHook,
			RootDir:                 args.RootDir,
			Verbose:                 args.Verbose,
		})
		if err != nil || exit {
			return config.EmptyValidatedConfig(), exit, err
		}
	}

	return validatedConfig, false, err
}

type ConfigArgs struct {
	Backend               git.BackendCommands
	BranchesSnapshot      gitdomain.BranchesSnapshot
	BranchesToValidate    gitdomain.LocalBranchNames
	CommandsCounter       gohacks.Counter
	ConfigSnapshot        undoconfig.ConfigSnapshot
	DialogTestInputs      components.TestInputs
	FinalMessages         stringslice.Collector
	Frontend              git.FrontendCommands
	HandleUnfinishedState bool
	LocalBranches         gitdomain.LocalBranchNames
	RepoStatus            gitdomain.RepoStatus
	RootDir               gitdomain.RepoRootDir
	StashSize             gitdomain.StashSize
	TestInputs            components.TestInputs
	Unvalidated           config.UnvalidatedConfig
	Verbose               bool
}

// cleanupPerennialParentEntries removes outdated entries from the configuration.
func cleanupPerennialParentEntries(lineage configdomain.Lineage, perennialBranches gitdomain.LocalBranchNames, access gitconfig.Access, finalMessages stringslice.Collector) error {
	for _, perennialBranch := range perennialBranches {
		if lineage.Parent(perennialBranch).IsSome() {
			if err := access.RemoveLocalConfigValue(gitconfig.NewParentKey(perennialBranch)); err != nil {
				return err
			}
			lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
	return nil
}
