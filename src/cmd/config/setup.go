package config

import (
	"slices"

	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/configfile"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/spf13/cobra"
)

const setupConfigDesc = "Prompts to setup your Git Town configuration"

func SetupCommand() *cobra.Command {
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

// the config settings to be used if the user accepts all default options
func defaultUserInput() userInput {
	return userInput{
		FullConfig:    configdomain.DefaultConfig(),
		configStorage: dialog.ConfigStorageOptionFile,
	}
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
	dialogInputs  components.TestInputs
	hasConfigFile bool
	localBranches gitdomain.BranchInfos
	userInput     userInput
}

type userInput struct {
	configdomain.FullConfig
	configStorage dialog.ConfigStorageOption
}

func enterData(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	aborted, err = dialog.Welcome(config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.Aliases, aborted, err = dialog.Aliases(configdomain.AllAliasableCommands(), runner.Aliases, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	existingMainBranch := runner.MainBranch
	if existingMainBranch.IsEmpty() {
		existingMainBranch, _ = runner.Backend.DefaultBranch()
	}
	config.userInput.MainBranch, aborted, err = dialog.MainBranch(config.localBranches.Names(), existingMainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PerennialBranches, aborted, err = dialog.PerennialBranches(config.localBranches.Names(), runner.PerennialBranches, config.userInput.MainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.HostingPlatform, aborted, err = dialog.HostingPlatform(runner.HostingPlatform, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	switch config.userInput.HostingPlatform {
	case configdomain.HostingPlatformBitbucket:
		// BitBucket API isn't supported yet
	case configdomain.HostingPlatformGitea:
		config.userInput.GiteaToken, aborted, err = dialog.GiteaToken(runner.GiteaToken, config.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, err
		}
	case configdomain.HostingPlatformGitHub:
		config.userInput.GitHubToken, aborted, err = dialog.GitHubToken(runner.GitHubToken, config.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, err
		}
	case configdomain.HostingPlatformGitLab:
		config.userInput.GitLabToken, aborted, err = dialog.GitLabToken(runner.GitLabToken, config.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, err
		}
	case configdomain.HostingPlatformNone:
	}
	config.userInput.HostingOriginHostname, aborted, err = dialog.OriginHostname(runner.HostingOriginHostname, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncFeatureStrategy, aborted, err = dialog.SyncFeatureStrategy(runner.SyncFeatureStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncPerennialStrategy, aborted, err = dialog.SyncPerennialStrategy(runner.SyncPerennialStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncUpstream, aborted, err = dialog.SyncUpstream(runner.SyncUpstream, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.NewBranchPush, aborted, err = dialog.PushNewBranches(runner.NewBranchPush, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PushHook, aborted, err = dialog.PushHook(runner.PushHook, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncBeforeShip, aborted, err = dialog.SyncBeforeShip(runner.SyncBeforeShip, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.ShipDeleteTrackingBranch, aborted, err = dialog.ShipDeleteTrackingBranch(runner.ShipDeleteTrackingBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.configStorage, aborted, err = dialog.ConfigStorage(config.hasConfigFile, config.dialogInputs.Next())
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
		dialogInputs:  dialogInputs,
		hasConfigFile: repo.Runner.Config.ConfigFile != nil,
		localBranches: branchesSnapshot.Branches,
		userInput:     defaultUserInput(),
	}, exit, err
}

func saveAll(runner *git.ProdRunner, userInput userInput) error {
	err := saveAliases(runner, userInput)
	if err != nil {
		return err
	}
	err = saveGiteaToken(runner, userInput)
	if err != nil {
		return err
	}
	err = saveGitHubToken(runner, userInput)
	if err != nil {
		return err
	}
	err = saveGitLabToken(runner, userInput)
	if err != nil {
		return err
	}
	switch userInput.configStorage {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(runner, userInput)
	}
	panic("unknown configStorage: " + userInput.configStorage)
}

func saveToGit(runner *git.ProdRunner, userInput userInput) error {
	err := saveHostingPlatform(runner, userInput)
	if err != nil {
		return err
	}
	err = saveOriginHostname(runner, userInput)
	if err != nil {
		return err
	}
	err = saveMainBranch(runner, userInput)
	if err != nil {
		return err
	}
	err = savePerennialBranches(runner, userInput)
	if err != nil {
		return err
	}
	err = savePushHook(runner, userInput)
	if err != nil {
		return err
	}
	err = savePushNewBranches(runner, userInput)
	if err != nil {
		return err
	}
	err = saveShipDeleteTrackingBranch(runner, userInput)
	if err != nil {
		return err
	}
	err = saveSyncFeatureStrategy(runner, userInput)
	if err != nil {
		return err
	}
	err = saveSyncPerennialStrategy(runner, userInput)
	if err != nil {
		return err
	}
	err = saveSyncUpstream(runner, userInput)
	if err != nil {
		return err
	}
	err = saveSyncBeforeShip(runner, userInput)
	if err != nil {
		return err
	}
	return nil
}

func saveAliases(runner *git.ProdRunner, userInput userInput) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := runner.Aliases[aliasableCommand]
		newAlias, hasNew := userInput.Aliases[aliasableCommand]
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

func saveGiteaToken(runner *git.ProdRunner, userInput userInput) error {
	if userInput.GiteaToken == runner.GiteaToken {
		return nil
	}
	return runner.Frontend.SetGiteaToken(userInput.GiteaToken)
}

func saveGitHubToken(runner *git.ProdRunner, userInput userInput) error {
	if userInput.GitHubToken == runner.GitHubToken {
		return nil
	}
	return runner.Frontend.SetGitHubToken(userInput.GitHubToken)
}

func saveGitLabToken(runner *git.ProdRunner, userInput userInput) error {
	if userInput.GitLabToken == runner.GitLabToken {
		return nil
	}
	return runner.Frontend.SetGitLabToken(userInput.GitLabToken)
}

func saveHostingPlatform(runner *git.ProdRunner, userInput userInput) (err error) {
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

func saveMainBranch(runner *git.ProdRunner, userInput userInput) error {
	if userInput.MainBranch == runner.MainBranch {
		return nil
	}
	return runner.SetMainBranch(userInput.MainBranch)
}

func saveOriginHostname(runner *git.ProdRunner, userInput userInput) error {
	if userInput.HostingOriginHostname == runner.HostingOriginHostname {
		return nil
	}
	if runner.HostingOriginHostname != "" && userInput.HostingOriginHostname == "" {
		return runner.Frontend.DeleteOriginHostname()
	}
	return runner.Frontend.SetOriginHostname(userInput.HostingOriginHostname)
}

func savePerennialBranches(runner *git.ProdRunner, config userInput) error {
	oldSetting := runner.PerennialBranches
	newSetting := config.PerennialBranches
	if slices.Compare(oldSetting, newSetting) != 0 || runner.LocalGitConfig.PerennialBranches == nil {
		return runner.SetPerennialBranches(newSetting)
	}
	return nil
}

func savePushHook(runner *git.ProdRunner, userInput userInput) error {
	if userInput.PushHook == runner.PushHook {
		return nil
	}
	return runner.SetPushHookLocally(userInput.PushHook)
}

func savePushNewBranches(runner *git.ProdRunner, userInput userInput) error {
	if userInput.NewBranchPush == runner.NewBranchPush {
		return nil
	}
	return runner.SetNewBranchPush(userInput.NewBranchPush, false)
}

func saveShipDeleteTrackingBranch(runner *git.ProdRunner, userInput userInput) error {
	if userInput.ShipDeleteTrackingBranch == runner.ShipDeleteTrackingBranch {
		return nil
	}
	return runner.SetShipDeleteTrackingBranch(userInput.ShipDeleteTrackingBranch, false)
}

func saveSyncFeatureStrategy(runner *git.ProdRunner, userInput userInput) error {
	if userInput.SyncFeatureStrategy == runner.SyncFeatureStrategy {
		return nil
	}
	return runner.SetSyncFeatureStrategy(userInput.SyncFeatureStrategy)
}

func saveSyncPerennialStrategy(runner *git.ProdRunner, userInput userInput) error {
	if userInput.SyncPerennialStrategy == runner.SyncPerennialStrategy {
		return nil
	}
	return runner.SetSyncPerennialStrategy(userInput.SyncPerennialStrategy)
}

func saveSyncUpstream(runner *git.ProdRunner, userInput userInput) error {
	if userInput.SyncUpstream == runner.SyncUpstream {
		return nil
	}
	return runner.SetSyncUpstream(userInput.SyncUpstream, false)
}

func saveSyncBeforeShip(runner *git.ProdRunner, userInput userInput) error {
	if userInput.SyncBeforeShip == runner.SyncBeforeShip {
		return nil
	}
	return runner.SetSyncBeforeShip(userInput.SyncBeforeShip, false)
}

func saveToFile(userInput userInput) error {
	return configfile.Save(&userInput.FullConfig)
}
