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
	aborted, err := enterAliases(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = enterMainBranch(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = enterPerennialBranches(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = enterHostingPlatform(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncFeatureStrategy(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncPerennialStrategy(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncUpstream(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupPushNewBranches(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupPushHook(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncBeforeShip(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupShipDeleteTrackingBranch(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	return saveUserInput(repo.Runner, config.userInput)
}

type setupConfig struct {
	localBranches gitdomain.BranchInfos
	inputs        dialog.TestInputs
	userInput     configdomain.FullConfig
}

func enterAliases(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	config.userInput.Aliases, aborted, err = enter.Aliases(configdomain.AllAliasableCommands(), runner.Aliases, config.inputs.Next())
	return aborted, err
}

func enterHostingPlatform(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	config.userInput.HostingPlatform, aborted, err = enter.HostingPlatform(runner.HostingPlatform, config.inputs.Next())
	return aborted, err
}

func enterMainBranch(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	existingValue := runner.MainBranch
	if existingValue.IsEmpty() {
		existingValue, _ = runner.Backend.DefaultBranch()
	}
	config.userInput.MainBranch, aborted, err = enter.MainBranch(config.localBranches.Names(), existingValue, config.inputs.Next())
	return aborted, err
}

func enterPerennialBranches(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	config.userInput.PerennialBranches, aborted, err = enter.PerennialBranches(config.localBranches.Names(), runner.PerennialBranches, config.userInput.MainBranch, config.inputs.Next())
	return aborted, err
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
		inputs:        dialogInputs,
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
	return nil
}

func setupPushHook(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newPushHook, aborted, err := enter.PushHook(runner.PushHook, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetPushHookLocally(newPushHook)
}

func setupPushNewBranches(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.PushNewBranches(runner.NewBranchPush, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetNewBranchPush(newValue, false)
}

func setupShipDeleteTrackingBranch(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.ShipDeleteTrackingBranch(runner.ShipDeleteTrackingBranch, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetShipDeleteTrackingBranch(newValue, false)
}

func setupSyncBeforeShip(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncBeforeShip(runner.SyncBeforeShip, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncBeforeShip(newValue, false)
}

func setupSyncFeatureStrategy(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncFeatureStrategy(runner.SyncFeatureStrategy, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncFeatureStrategy(newValue)
}

func setupSyncPerennialStrategy(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncPerennialStrategy(runner.SyncPerennialStrategy, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncPerennialStrategy(newValue)
}

func setupSyncUpstream(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.SyncUpstream(runner.SyncUpstream, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncUpstream(newValue, false)
}
