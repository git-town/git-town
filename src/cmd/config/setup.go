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
	aborted, err := setupAliases(repo.Runner.FullConfig.Aliases, configdomain.AllAliasableCommands(), repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupMainBranch(config.MainBranch, config.localBranches.Names(), repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupPerennialBranches(repo.Runner.PerennialBranches, repo.Runner.MainBranch, config.localBranches.Names(), repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupHostingPlatform(config.HostingPlatform, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupHostingToken()
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncFeatureStrategy(config.SyncFeatureStrategy, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncPerennialStrategy(config.SyncPerennialStrategy, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncUpstream(config.SyncUpstream, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupPushNewBranches(config.NewBranchPush, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupPushHook(config.PushHook, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupSyncBeforeShip(config.SyncBeforeShip, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
		return err
	}
	aborted, err = setupShipDeleteTrackingBranch(config.ShipDeleteTrackingBranch, repo.Runner, config.dialogInputs.Next())
	if err != nil || aborted {
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

func setupAliases(existingValue configdomain.Aliases, allAliasableCommands configdomain.AliasableCommands, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newAliases, aborted, err := enter.Aliases(allAliasableCommands, existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	for _, aliasableCommand := range allAliasableCommands {
		newAlias, hasNew := newAliases[aliasableCommand]
		oldAlias, hasOld := existingValue[aliasableCommand]
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

func setupHostingPlatform(existingValue configdomain.HostingPlatform, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.HostingPlatform(existingValue, inputs)
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

func setupHostingToken() (bool, error) {
	// switch newCodeHostingPlatform {
	// case configdomain.CodeHostingPlatformAutoDetect:
	// }
	return false, nil
}

func setupGitHubToken() (bool, error) {
	// newToken, err := dialog.EnterGitHubToken(newCodeHostingPlatform, config.dialogInputs.Next())
	// if err != nil {
	// 	return err
	// }
	// err = repo.Runner.Frontend.SetCodeHostingToken(newCodeHostingToken)
	// if err != nil {
	// 	return err
	// }
	return false, nil
}

func setupMainBranch(existingValue gitdomain.LocalBranchName, allBranches gitdomain.LocalBranchNames, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	if existingValue.IsEmpty() {
		existingValue, _ = runner.Backend.DefaultBranch()
	}
	newMainBranch, aborted, err := enter.MainBranch(allBranches, existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetMainBranch(newMainBranch)
}

func setupPerennialBranches(existingValue gitdomain.LocalBranchNames, mainBranch gitdomain.LocalBranchName, allBranches gitdomain.LocalBranchNames, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.PerennialBranches(allBranches, existingValue, mainBranch, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	if slices.Compare(runner.PerennialBranches, newValue) != 0 || runner.LocalGitConfig.PerennialBranches == nil {
		err = runner.SetPerennialBranches(newValue)
	}
	return aborted, err
}

func setupPushHook(existingValue configdomain.PushHook, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newPushHook, aborted, err := enter.PushHook(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetPushHookLocally(newPushHook)
}

func setupPushNewBranches(existingValue configdomain.NewBranchPush, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.PushNewBranches(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetNewBranchPush(newValue, false)
}

func setupShipDeleteTrackingBranch(existingValue configdomain.ShipDeleteTrackingBranch, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.ShipDeleteTrackingBranch(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetShipDeleteTrackingBranch(newValue, false)
}

func setupSyncBeforeShip(existingValue configdomain.SyncBeforeShip, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.SyncBeforeShip(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncBeforeShip(newValue, false)
}

func setupSyncFeatureStrategy(existingValue configdomain.SyncFeatureStrategy, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.SyncFeatureStrategy(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncFeatureStrategy(newValue)
}

func setupSyncPerennialStrategy(existingValue configdomain.SyncPerennialStrategy, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.SyncPerennialStrategy(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncPerennialStrategy(newValue)
}

func setupSyncUpstream(existingValue configdomain.SyncUpstream, runner *git.ProdRunner, inputs dialog.TestInput) (bool, error) {
	newValue, aborted, err := enter.SyncUpstream(existingValue, inputs)
	if err != nil || aborted {
		return aborted, err
	}
	return aborted, runner.SetSyncUpstream(newValue, false)
}
