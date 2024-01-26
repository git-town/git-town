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
	return saveUserInput(repo.Runner, config.userInput)
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
	aborted, err = setupSyncFeatureStrategy(runner, config)
	if err != nil || aborted {
		return aborted, err
	}
	aborted, err = setupSyncPerennialStrategy(runner, config)
	if err != nil || aborted {
		return aborted, err
	}
	aborted, err = setupSyncUpstream(runner, config)
	if err != nil || aborted {
		return aborted, err
	}
	aborted, err = setupPushNewBranches(runner, config)
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PushHook, aborted, err = enter.PushHook(runner.PushHook, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	aborted, err = setupSyncBeforeShip(runner, config)
	if err != nil || aborted {
		return aborted, err
	}
	aborted, err = setupShipDeleteTrackingBranch(runner, config)
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
		userInput:     configdomain.FullConfig{},
	}, exit, err
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
	if newConfig.PushHook != runner.PushHook {
		return runner.SetPushHookLocally(newConfig.PushHook)
	}
	return nil
}

func saveUserInput(runner *git.ProdRunner, newConfig configdomain.FullConfig) error {
	err := saveAliases(runner, newConfig)
	if err != nil {
		return err
	}
	err = saveHostingPlatform(runner, newConfig)
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
	return nil
}

func setupPushNewBranches(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.PushNewBranches(runner.NewBranchPush, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetNewBranchPush(newValue, false)
}

func setupShipDeleteTrackingBranch(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.ShipDeleteTrackingBranch(runner.ShipDeleteTrackingBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetShipDeleteTrackingBranch(newValue, false)
}

func setupSyncBeforeShip(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncBeforeShip(runner.SyncBeforeShip, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncBeforeShip(newValue, false)
}

func setupSyncFeatureStrategy(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncFeatureStrategy(runner.SyncFeatureStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncFeatureStrategy(newValue)
}

func setupSyncPerennialStrategy(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncPerennialStrategy(runner.SyncPerennialStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncPerennialStrategy(newValue)
}

func setupSyncUpstream(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncUpstream(runner.SyncUpstream, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncUpstream(newValue, false)
}
