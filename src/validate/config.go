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

func Config(args ConfigArgs) (*config.ValidatedConfig, *git.ProdRunner, bool, error) {
	// check Git user data
	if args.Unvalidated.Config.GitUserEmail.IsNone() {
		return nil, nil, false, errors.New(messages.GitUserEmailMissing)
	}
	if args.Unvalidated.Config.GitUserName.IsNone() {
		return nil, nil, false, errors.New(messages.GitUserNameMissing)
	}

	// enter and save main and perennials
	validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
		DialogInputs:          args.TestInputs,
		GetDefaultBranch:      args.Backend.DefaultBranch,
		HasConfigFile:         args.Unvalidated.ConfigFile.IsSome(),
		LocalBranches:         args.LocalBranches,
		UnvalidatedMain:       args.Unvalidated.Config.MainBranch,
		UnvalidatedPerennials: args.Unvalidated.Config.PerennialBranches,
	})
	if err != nil || aborted {
		return nil, nil, aborted, err
	}
	if err = args.Unvalidated.SetMainBranch(validatedMain); err != nil {
		return nil, nil, false, err
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return nil, nil, false, err
		}
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, abort, err := dialog.Lineage(dialog.LineageArgs{
		BranchesToVerify: args.BranchesToValidate,
		Config:           args.Unvalidated.Config,
		DefaultChoice:    validatedMain,
		DialogTestInputs: args.TestInputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       validatedMain,
	})
	if err != nil || abort {
		return nil, nil, abort, err
	}
	for branch, parent := range additionalLineage {
		if err = args.Unvalidated.SetParent(branch, parent); err != nil {
			return nil, nil, abort, err
		}
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return nil, nil, false, err
		}
	}

	// remove outdated lineage
	err = args.Unvalidated.RemoveOutdatedConfiguration(args.LocalBranches)
	if err != nil {
		return nil, nil, false, err
	}
	err = cleanupPerennialParentEntries(args.Unvalidated.Config.Lineage, args.Unvalidated.Config.PerennialBranches, args.Unvalidated.GitConfig, args.FinalMessages)
	if err != nil {
		return nil, nil, false, err
	}

	// create validated configuration
	validatedConfig := config.ValidatedConfig{
		Config: configdomain.ValidatedConfig{
			UnvalidatedConfig: args.Unvalidated.Config,
			MainBranch:        validatedMain,
		},
	}

	runner := git.ProdRunner{
		Config:          &validatedConfig,
		Backend:         *args.Backend,
		Frontend:        args.Frontend,
		CommandsCounter: args.CommandsCounter,
		FinalMessages:   args.FinalMessages,
	}

	// handle unfinished state
	exit, err := HandleUnfinishedState(UnfinishedStateArgs{
		Connector:               nil,
		CurrentBranch:           args.BranchesSnapshot.Active,
		DialogTestInputs:        args.DialogTestInputs,
		HasOpenChanges:          args.RepoStatus.OpenChanges,
		InitialBranchesSnapshot: args.BranchesSnapshot,
		InitialConfigSnapshot:   args.ConfigSnapshot,
		InitialStashSize:        args.StashSize,
		Lineage:                 validatedConfig.Config.Lineage,
		PushHook:                validatedConfig.Config.PushHook,
		RootDir:                 args.RootDir,
		Run:                     &runner,
		Verbose:                 args.Verbose,
	})
	if err != nil || exit {
		return nil, &runner, false, err
	}

	return &validatedConfig, &runner, false, err
}

type ConfigArgs struct {
	Backend            *git.BackendCommands
	BranchesSnapshot   gitdomain.BranchesSnapshot
	BranchesToValidate gitdomain.LocalBranchNames
	CommandsCounter    *gohacks.Counter
	ConfigSnapshot     undoconfig.ConfigSnapshot
	DialogTestInputs   components.TestInputs
	FinalMessages      *stringslice.Collector
	Frontend           git.FrontendCommands
	LocalBranches      gitdomain.LocalBranchNames
	RepoStatus         gitdomain.RepoStatus
	RootDir            gitdomain.RepoRootDir
	StashSize          gitdomain.StashSize
	TestInputs         *components.TestInputs
	Unvalidated        config.UnvalidatedConfig
	Verbose            bool
}

// cleanupPerennialParentEntries removes outdated entries from the configuration.
func cleanupPerennialParentEntries(lineage configdomain.Lineage, perennialBranches gitdomain.LocalBranchNames, access gitconfig.Access, finalMessages *stringslice.Collector) error {
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
