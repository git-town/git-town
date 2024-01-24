package config

import (
	"slices"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const setupConfigDesc = "Prompts to setup your Git Town configuration"

func setupCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: setupConfigDesc,
		Long:  cmdhelpers.Long(setupConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeConfigSetup(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigSetup(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := loadSetupConfig(repo, verbose)
	if err != nil || exit {
		return err
	}

	// ALIASES
	allAliasableCommands := configdomain.AllAliasableCommands()
	newAliases, aborted, err := dialog.Aliases(allAliasableCommands, repo.Runner.FullConfig.Aliases, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	for _, aliasableCommand := range allAliasableCommands {
		newAlias, hasNew := newAliases[aliasableCommand]
		oldAlias, hasOld := config.FullConfig.Aliases[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err := repo.Runner.Frontend.RemoveGitAlias(aliasableCommand)
			if err != nil {
				return err
			}
		case newAlias != oldAlias:
			err := repo.Runner.Frontend.SetGitAlias(aliasableCommand)
			if err != nil {
				return err
			}
		}
	}

	// MAIN BRANCH
	defaultMainBranch := repo.Runner.MainBranch
	if defaultMainBranch.IsEmpty() {
		defaultMainBranch, _ = repo.Runner.Backend.DefaultBranch()
	}
	newMainBranch, aborted, err := dialog.EnterMainBranch(config.localBranches.Names(), defaultMainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetMainBranch(newMainBranch)
	if err != nil {
		return err
	}

	// PERENNIAL BRANCHES
	newPerennialBranches, aborted, err := dialog.EnterPerennialBranches(config.localBranches.Names(), repo.Runner.PerennialBranches, repo.Runner.MainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	if slices.Compare(repo.Runner.PerennialBranches, newPerennialBranches) != 0 || repo.Runner.LocalGitConfig.PerennialBranches == nil {
		err = repo.Runner.SetPerennialBranches(newPerennialBranches)
		if err != nil {
			return err
		}
	}

	// CODE HOSTING
	newCodeHostingPlatform, aborted, err := dialog.EnterHostingPlatform(config.CodeHostingPlatform, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	switch {
	case config.CodeHostingPlatform == "" && newCodeHostingPlatform == configdomain.CodeHostingPlatformAutoDetect:
		// no changes --> do nothing
	case config.CodeHostingPlatform != "" && newCodeHostingPlatform == configdomain.CodeHostingPlatformAutoDetect:
		err = repo.Runner.Frontend.DeleteCodeHostingPlatform()
		if err != nil {
			return err
		}
	case config.CodeHostingPlatform.String() != newCodeHostingPlatform:
		err = repo.Runner.Frontend.SetCodeHostingPlatform(newCodeHostingPlatform)
		if err != nil {
			return err
		}
	}

	// SYNC-FEATURE-STRATEGY
	newSyncFeatureStrategy, aborted, err := dialog.EnterSyncFeatureStrategy(config.SyncFeatureStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetSyncFeatureStrategy(newSyncFeatureStrategy)
	if err != nil {
		return err
	}

	// SYNC-PERENNIAL-STRATEGY
	newSyncPerennialStrategy, aborted, err := dialog.EnterSyncPerennialStrategy(config.SyncPerennialStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetSyncPerennialStrategy(newSyncPerennialStrategy)
	if err != nil {
		return err
	}

	// SYNC UPSTREAM
	newSyncUpstream, aborted, err := dialog.EnterSyncUpstream(config.SyncUpstream, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetSyncUpstream(newSyncUpstream, false)
	if err != nil {
		return err
	}

	// PUSH NEW BRANCHES
	newPushNewBranches, aborted, err := dialog.EnterPushNewBranches(config.NewBranchPush, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetNewBranchPush(newPushNewBranches, false)
	if err != nil {
		return err
	}

	// PUSH HOOK
	newPushHook, aborted, err := dialog.EnterPushHook(config.PushHook, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetPushHookLocally(newPushHook)
	if err != nil {
		return err
	}

	// SYNC BEFORE SHIP
	newSyncBeforeShip, aborted, err := dialog.EnterSyncBeforeShip(config.SyncBeforeShip, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetSyncBeforeShip(newSyncBeforeShip, false)
	if err != nil {
		return err
	}

	// SHIP DELETE TRACKING BRANCH
	newShipDeleteTrackingBranch, aborted, err := dialog.EnterShipDeleteTrackingBranch(config.ShipDeleteTrackingBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	err = repo.Runner.SetShipDeleteTrackingBranch(newShipDeleteTrackingBranch, false)
	if err != nil {
		return err
	}

	return nil
}

type setupConfig struct {
	*configdomain.FullConfig
	localBranches gitdomain.BranchInfos
	dialogInputs  dialog.TestInputs
}

func loadSetupConfig(repo *execute.OpenRepoResult, verbose bool) (setupConfig, bool, error) {
	branchesSnapshot, _, dialogInputs, exit, err := execute.LoadRepoSnapshot(execute.LoadBranchesArgs{
		FullConfig:            &repo.Runner.FullConfig,
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateNoOpenChanges: false,
	})
	return setupConfig{
		FullConfig:    &repo.Runner.FullConfig,
		localBranches: branchesSnapshot.Branches,
		dialogInputs:  dialogInputs,
	}, exit, err
}
