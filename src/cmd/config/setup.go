package config

import (
	"slices"

	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
	"github.com/git-town/git-town/v11/src/cli/dialogs/enter"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
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
	aborted, err := enterData(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	return saveAll(repo.Runner, config.userInput)
}

type setupConfig struct {
	localBranches gitdomain.BranchInfos
	dialogInputs  dialog.TestInputs
	userInput     configdomain.FullConfig
}

func enterData(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	config.userInput.Aliases, aborted, err = enter.Aliases(configdomain.AllAliasableCommands(), runner.Aliases, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	existingMainBranch := runner.MainBranch
	if existingMainBranch.IsEmpty() {
		existingMainBranch, _ = runner.Backend.DefaultBranch()
	}
	config.userInput.MainBranch, aborted, err = enter.MainBranch(config.localBranches.Names(), existingMainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PerennialBranches, aborted, err = enter.PerennialBranches(config.localBranches.Names(), runner.PerennialBranches, config.userInput.MainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.HostingPlatform, aborted, err = enter.HostingPlatform(runner.HostingPlatform, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	switch config.userInput.HostingPlatform {
	case configdomain.HostingPlatformGitea:
		config.userInput.GiteaToken, aborted, err = enter.GiteaToken(runner.GiteaToken, config.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, err
		}
	case configdomain.HostingPlatformGitHub:
		config.userInput.GitHubToken, aborted, err = enter.GitHubToken(runner.GitHubToken, config.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, err
		}
	case configdomain.HostingPlatformGitLab:
		config.userInput.GitLabToken, aborted, err = enter.GitLabToken(runner.GitLabToken, config.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, err
		}
	}
	config.userInput.SyncFeatureStrategy, aborted, err = enter.SyncFeatureStrategy(runner.SyncFeatureStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncPerennialStrategy, aborted, err = enter.SyncPerennialStrategy(runner.SyncPerennialStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncUpstream, aborted, err = enter.SyncUpstream(runner.SyncUpstream, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.NewBranchPush, aborted, err = enter.PushNewBranches(runner.NewBranchPush, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PushHook, aborted, err = enter.PushHook(runner.PushHook, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncBeforeShip, aborted, err = enter.SyncBeforeShip(runner.SyncBeforeShip, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.ShipDeleteTrackingBranch, aborted, err = enter.ShipDeleteTrackingBranch(runner.ShipDeleteTrackingBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return false, nil
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
		localBranches: branchesSnapshot.Branches,
		dialogInputs:  dialogInputs,
		userInput:     configdomain.FullConfig{}, //nolint:exhaustruct
	}, exit, err
}

func saveAll(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	err := saveAliases(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveHostingPlatform(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveGiteaToken(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveGitHubToken(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveGitLabToken(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveMainBranch(runner, newConfig)
	if err != nil {
		return err
	}
	err = savePerennialBranches(runner, newConfig)
	if err != nil {
		return err
	}
	err = savePushHook(runner, newConfig)
	if err != nil {
		return err
	}
	err = savePushNewBranches(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveShipDeleteTrackingBranch(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveSyncFeatureStrategy(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveSyncPerennialStrategy(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveSyncUpstream(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveSyncBeforeShip(runner, newConfig)
	if err != nil {
		return err
	}
	return nil
}

func saveAliases(runner *git.ProdRunner, newConfig configdomain.FullConfig) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := runner.Aliases[aliasableCommand]
		newAlias, hasNew := newConfig.Aliases[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err = runner.Frontend.RemoveGitAlias(aliasableCommand)
		case newAlias != oldAlias:
			err = runner.Frontend.SetGitAlias(aliasableCommand)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveGiteaToken(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.GiteaToken == runner.GiteaToken {
		return nil
	}
	return runner.Frontend.SetGiteaToken(newConfig.GiteaToken)
}

func saveGitHubToken(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.GitHubToken == runner.GitHubToken {
		return nil
	}
	return runner.Frontend.SetGitHubToken(newConfig.GitHubToken)
}

func saveGitLabToken(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.GitLabToken == runner.GitLabToken {
		return nil
	}
	return runner.Frontend.SetGitLabToken(newConfig.GitLabToken)
}

func saveHostingPlatform(runner *git.ProdRunner, userInput configdomain.FullConfig) (err error) {
	oldValue := runner.HostingPlatform
	newValue := userInput.HostingPlatform
	switch {
	case oldValue == "" && newValue == configdomain.HostingPlatformNone:
		// no changes --> do nothing
	case oldValue != "" && newValue == configdomain.HostingPlatformNone:
		return runner.Frontend.DeleteHostingPlatform()
	case oldValue != newValue:
		return runner.Frontend.SetHostingPlatform(newValue)
	}
	return nil
}

func saveMainBranch(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.MainBranch == runner.MainBranch {
		return nil
	}
	return runner.SetMainBranch(newConfig.MainBranch)
}

func savePerennialBranches(runner *git.ProdRunner, config configdomain.FullConfig) error {
	oldSetting := runner.PerennialBranches
	newSetting := config.PerennialBranches
	if slices.Compare(oldSetting, newSetting) != 0 || runner.LocalGitConfig.PerennialBranches == nil {
		return runner.SetPerennialBranches(newSetting)
	}
	return nil
}

func savePushHook(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.PushHook == runner.PushHook {
		return nil
	}
	return runner.SetPushHookLocally(newConfig.PushHook)
}

func savePushNewBranches(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.NewBranchPush == runner.NewBranchPush {
		return nil
	}
	return runner.SetNewBranchPush(newConfig.NewBranchPush, false)
}

func saveShipDeleteTrackingBranch(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.ShipDeleteTrackingBranch == runner.ShipDeleteTrackingBranch {
		return nil
	}
	return runner.SetShipDeleteTrackingBranch(newConfig.ShipDeleteTrackingBranch, false)
}

func saveSyncFeatureStrategy(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.SyncFeatureStrategy == runner.SyncFeatureStrategy {
		return nil
	}
	return runner.SetSyncFeatureStrategy(newConfig.SyncFeatureStrategy)
}

func saveSyncPerennialStrategy(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.SyncPerennialStrategy == runner.SyncPerennialStrategy {
		return nil
	}
	return runner.SetSyncPerennialStrategy(newConfig.SyncPerennialStrategy)
}

func saveSyncUpstream(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.SyncUpstream == runner.SyncUpstream {
		return nil
	}
	return runner.SetSyncUpstream(newConfig.SyncUpstream, false)
}

func saveSyncBeforeShip(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	if newConfig.SyncBeforeShip == runner.SyncBeforeShip {
		return nil
	}
	return runner.SetSyncBeforeShip(newConfig.SyncBeforeShip, false)
}
