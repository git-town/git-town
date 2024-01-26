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
	aborted, err := setupAliases(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupMainBranch(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupPerennialBranches(repo.Runner, &config)
	if err != nil || aborted {
		return err
	}
	aborted, err = setupHostingPlatform(repo.Runner, &config)
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
	return nil
}

type setupConfig struct {
	localBranches gitdomain.BranchInfos
	inputs        dialog.TestInputs
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
	}, exit, err
}

func setupAliases(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	aliasableCommands := configdomain.AllAliasableCommands()
	newAliases, aborted, err := enter.Aliases(aliasableCommands, runner.Aliases, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	for _, aliasableCommand := range aliasableCommands {
		newAlias, hasNew := newAliases[aliasableCommand]
		oldAlias, hasOld := runner.Aliases[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err := runner.Frontend.RemoveGitAlias(aliasableCommand)
			if err != nil {
				return aborted, err
			}
		case newAlias != oldAlias:
			err := runner.Frontend.SetGitAlias(aliasableCommand)
			if err != nil {
				return aborted, err
			}
		}
	}
	return aborted, nil
}

func setupHostingPlatform(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	existingValue := runner.HostingPlatform
	newValue, aborted, err := enter.HostingPlatform(existingValue, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	switch {
	case existingValue == "" && newValue == configdomain.HostingPlatformNone:
		// no changes --> do nothing
	case existingValue != "" && newValue == configdomain.HostingPlatformNone:
		return aborted, runner.Frontend.DeleteHostingPlatform()
	case existingValue != newValue:
		return aborted, runner.Frontend.SetHostingPlatform(newValue)
	}
	return aborted, nil
}

func setupMainBranch(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	existingValue := runner.MainBranch
	if existingValue.IsEmpty() {
		existingValue, _ = runner.Backend.DefaultBranch()
	}
	newMainBranch, aborted, err := enter.MainBranch(config.localBranches.Names(), existingValue, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetMainBranch(newMainBranch)
}

func setupPerennialBranches(runner *git.ProdRunner, config *setupConfig) (bool, error) {
	newValue, aborted, err := enter.PerennialBranches(config.localBranches.Names(), runner.PerennialBranches, runner.MainBranch, config.inputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	if slices.Compare(runner.PerennialBranches, newValue) != 0 || runner.LocalGitConfig.PerennialBranches == nil {
		err = runner.SetPerennialBranches(newValue)
	}
	return aborted, err
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
