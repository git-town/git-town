package config

import (
	"os"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/configfile"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeConfigSetup(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

// the config settings to be used if the user accepts all default options
func defaultUserInput() userInput {
	return userInput{
		UnvalidatedConfig: configdomain.DefaultConfig(),
		configStorage:     dialog.ConfigStorageOptionFile,
	}
}

func executeConfigSetup(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, exit, err := loadSetupConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	aborted, err := enterData(repo.Runner, config)
	if err != nil || aborted {
		return err
	}
	err = saveAll(repo.Runner, config.userInput)
	if err != nil {
		return err
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "setup",
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		RootDir:             repo.RootDir,
		Runner:              repo.Runner,
		Verbose:             verbose,
	})
}

type setupConfig struct {
	dialogInputs  components.TestInputs
	hasConfigFile bool
	localBranches gitdomain.BranchInfos
	userInput     userInput
}

type userInput struct {
	configdomain.UnvalidatedConfig
	configStorage dialog.ConfigStorageOption
}

func determineHostingPlatform(runner *git.ProdRunner, userChoice Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userChoice.IsSome() {
		return userChoice
	}
	return hosting.Detect(runner.Config.OriginURL(), userChoice)
}

func enterData(runner *git.ProdRunner, config *setupConfig) (aborted bool, err error) {
	aborted, err = dialog.Welcome(config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.Aliases, aborted, err = dialog.Aliases(configdomain.AllAliasableCommands(), runner.Config.FullConfig.Aliases, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	existingMainBranch := runner.Config.FullConfig.MainBranch
	if existingMainBranch.IsEmpty() {
		existingMainBranch = runner.Backend.DefaultBranch()
	}
	if existingMainBranch.IsEmpty() {
		existingMainBranch = runner.Backend.OriginHead()
	}
	config.userInput.MainBranch, aborted, err = dialog.MainBranch(config.localBranches.Names(), existingMainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PerennialBranches, aborted, err = dialog.PerennialBranches(config.localBranches.Names(), runner.Config.FullConfig.PerennialBranches, config.userInput.MainBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PerennialRegex, aborted, err = dialog.PerennialRegex(runner.Config.FullConfig.PerennialRegex, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.HostingPlatform, aborted, err = dialog.HostingPlatform(runner.Config.FullConfig.HostingPlatform, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}

	if platform, has := determineHostingPlatform(runner, config.userInput.HostingPlatform).Get(); has {
		switch platform {
		case configdomain.HostingPlatformBitbucket:
			// BitBucket API isn't supported yet
		case configdomain.HostingPlatformGitea:
			config.userInput.GiteaToken, aborted, err = dialog.GiteaToken(runner.Config.FullConfig.GiteaToken, config.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitHub:
			config.userInput.GitHubToken, aborted, err = dialog.GitHubToken(runner.Config.FullConfig.GitHubToken, config.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitLab:
			config.userInput.GitLabToken, aborted, err = dialog.GitLabToken(runner.Config.FullConfig.GitLabToken, config.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		}
	}
	config.userInput.HostingOriginHostname, aborted, err = dialog.OriginHostname(runner.Config.FullConfig.HostingOriginHostname, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncFeatureStrategy, aborted, err = dialog.SyncFeatureStrategy(runner.Config.FullConfig.SyncFeatureStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncPerennialStrategy, aborted, err = dialog.SyncPerennialStrategy(runner.Config.FullConfig.SyncPerennialStrategy, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncUpstream, aborted, err = dialog.SyncUpstream(runner.Config.FullConfig.SyncUpstream, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PushNewBranches, aborted, err = dialog.PushNewBranches(runner.Config.FullConfig.PushNewBranches, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.PushHook, aborted, err = dialog.PushHook(runner.Config.FullConfig.PushHook, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.SyncBeforeShip, aborted, err = dialog.SyncBeforeShip(runner.Config.FullConfig.SyncBeforeShip, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.ShipDeleteTrackingBranch, aborted, err = dialog.ShipDeleteTrackingBranch(runner.Config.FullConfig.ShipDeleteTrackingBranch, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	config.userInput.configStorage, aborted, err = dialog.ConfigStorage(config.hasConfigFile, config.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return false, nil
}

func loadSetupConfig(repo *execute.OpenRepoResult, verbose bool) (*setupConfig, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  false,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	return &setupConfig{
		dialogInputs:  dialogTestInputs,
		hasConfigFile: repo.Runner.Config.ConfigFile != nil,
		localBranches: branchesSnapshot.Branches,
		userInput:     defaultUserInput(),
	}, exit, err
}

func saveAll(runner *git.ProdRunner, userInput userInput) error {
	err := saveAliases(runner, userInput.Aliases)
	if err != nil {
		return err
	}
	err = saveGiteaToken(runner, userInput.GiteaToken)
	if err != nil {
		return err
	}
	err = saveGitHubToken(runner, userInput.GitHubToken)
	if err != nil {
		return err
	}
	err = saveGitLabToken(runner, userInput.GitLabToken)
	if err != nil {
		return err
	}
	switch userInput.configStorage {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, runner)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(runner, userInput)
	}
	panic("unknown configStorage: " + userInput.configStorage)
}

func saveToGit(runner *git.ProdRunner, userInput userInput) error {
	err := saveHostingPlatform(runner, userInput.HostingPlatform)
	if err != nil {
		return err
	}
	err = saveOriginHostname(runner, userInput.HostingOriginHostname)
	if err != nil {
		return err
	}
	err = saveMainBranch(runner, userInput.MainBranch)
	if err != nil {
		return err
	}
	err = savePerennialBranches(runner, userInput.PerennialBranches)
	if err != nil {
		return err
	}
	err = savePerennialRegex(runner, userInput.PerennialRegex)
	if err != nil {
		return err
	}
	err = savePushHook(runner, userInput.PushHook)
	if err != nil {
		return err
	}
	err = savePushNewBranches(runner, userInput.PushNewBranches)
	if err != nil {
		return err
	}
	err = saveShipDeleteTrackingBranch(runner, userInput.ShipDeleteTrackingBranch)
	if err != nil {
		return err
	}
	err = saveSyncFeatureStrategy(runner, userInput.SyncFeatureStrategy)
	if err != nil {
		return err
	}
	err = saveSyncPerennialStrategy(runner, userInput.SyncPerennialStrategy)
	if err != nil {
		return err
	}
	err = saveSyncUpstream(runner, userInput.SyncUpstream)
	if err != nil {
		return err
	}
	err = saveSyncBeforeShip(runner, userInput.SyncBeforeShip)
	if err != nil {
		return err
	}
	return nil
}

func saveAliases(runner *git.ProdRunner, newAliases configdomain.Aliases) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := runner.Config.FullConfig.Aliases[aliasableCommand]
		newAlias, hasNew := newAliases[aliasableCommand]
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

func saveGiteaToken(runner *git.ProdRunner, newToken Option[configdomain.GiteaToken]) error {
	if newToken == runner.Config.FullConfig.GiteaToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return runner.Frontend.SetGiteaToken(value)
	}
	return runner.Frontend.RemoveGiteaToken()
}

func saveGitHubToken(runner *git.ProdRunner, newToken Option[configdomain.GitHubToken]) error {
	if newToken == runner.Config.FullConfig.GitHubToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return runner.Frontend.SetGitHubToken(value)
	}
	return runner.Frontend.RemoveGitHubToken()
}

func saveGitLabToken(runner *git.ProdRunner, newToken Option[configdomain.GitLabToken]) error {
	if newToken == runner.Config.FullConfig.GitLabToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return runner.Frontend.SetGitLabToken(value)
	}
	return runner.Frontend.RemoveGitLabToken()
}

func saveHostingPlatform(runner *git.ProdRunner, newHostingPlatform Option[configdomain.HostingPlatform]) (err error) {
	oldValue, oldHas := runner.Config.FullConfig.HostingPlatform.Get()
	newValue, newHas := newHostingPlatform.Get()
	if !oldHas && !newHas {
		return nil
	}
	if oldValue == newValue {
		return nil
	}
	if newHas {
		return runner.Frontend.SetHostingPlatform(newValue)
	}
	return runner.Frontend.DeleteHostingPlatform()
}

func saveMainBranch(runner *git.ProdRunner, newValue gitdomain.LocalBranchName) error {
	if newValue == runner.Config.FullConfig.MainBranch {
		return nil
	}
	return runner.Config.SetMainBranch(newValue)
}

func saveOriginHostname(runner *git.ProdRunner, newValue Option[configdomain.HostingOriginHostname]) error {
	if newValue == runner.Config.FullConfig.HostingOriginHostname {
		return nil
	}
	if value, has := newValue.Get(); has {
		return runner.Frontend.SetOriginHostname(value)
	}
	return runner.Frontend.DeleteOriginHostname()
}

func savePerennialBranches(runner *git.ProdRunner, newValue gitdomain.LocalBranchNames) error {
	oldValue := runner.Config.FullConfig.PerennialBranches
	if slices.Compare(oldValue, newValue) != 0 || runner.Config.LocalGitConfig.PerennialBranches == nil {
		return runner.Config.SetPerennialBranches(newValue)
	}
	return nil
}

func savePerennialRegex(runner *git.ProdRunner, newValue Option[configdomain.PerennialRegex]) error {
	if newValue == runner.Config.FullConfig.PerennialRegex {
		return nil
	}
	if value, has := newValue.Get(); has {
		return runner.Config.SetPerennialRegexLocally(value)
	}
	runner.Config.RemovePerennialRegex()
	return nil
}

func savePushHook(runner *git.ProdRunner, newValue configdomain.PushHook) error {
	if newValue == runner.Config.FullConfig.PushHook {
		return nil
	}
	return runner.Config.SetPushHookLocally(newValue)
}

func savePushNewBranches(runner *git.ProdRunner, newValue configdomain.PushNewBranches) error {
	if newValue == runner.Config.FullConfig.PushNewBranches {
		return nil
	}
	return runner.Config.SetPushNewBranches(newValue, false)
}

func saveShipDeleteTrackingBranch(runner *git.ProdRunner, newValue configdomain.ShipDeleteTrackingBranch) error {
	if newValue == runner.Config.FullConfig.ShipDeleteTrackingBranch {
		return nil
	}
	return runner.Config.SetShipDeleteTrackingBranch(newValue, false)
}

func saveSyncFeatureStrategy(runner *git.ProdRunner, newValue configdomain.SyncFeatureStrategy) error {
	if newValue == runner.Config.FullConfig.SyncFeatureStrategy {
		return nil
	}
	return runner.Config.SetSyncFeatureStrategy(newValue)
}

func saveSyncPerennialStrategy(runner *git.ProdRunner, newValue configdomain.SyncPerennialStrategy) error {
	if newValue == runner.Config.FullConfig.SyncPerennialStrategy {
		return nil
	}
	return runner.Config.SetSyncPerennialStrategy(newValue)
}

func saveSyncUpstream(runner *git.ProdRunner, newValue configdomain.SyncUpstream) error {
	if newValue == runner.Config.FullConfig.SyncUpstream {
		return nil
	}
	return runner.Config.SetSyncUpstream(newValue, false)
}

func saveSyncBeforeShip(runner *git.ProdRunner, newValue configdomain.SyncBeforeShip) error {
	if newValue == runner.Config.FullConfig.SyncBeforeShip {
		return nil
	}
	return runner.Config.SetSyncBeforeShip(newValue, false)
}

func saveToFile(userInput userInput, runner *git.ProdRunner) error {
	err := configfile.Save(&userInput.UnvalidatedConfig)
	if err != nil {
		return err
	}
	runner.Config.RemoveMainBranch()
	runner.Config.RemovePerennialBranches()
	runner.Config.RemovePerennialRegex()
	runner.Config.RemovePushNewBranches()
	runner.Config.RemovePushHook()
	runner.Config.RemoveSyncBeforeShip()
	runner.Config.RemoveShipDeleteTrackingBranch()
	runner.Config.RemoveSyncFeatureStrategy()
	runner.Config.RemoveSyncPerennialStrategy()
	runner.Config.RemoveSyncUpstream()
	return nil
}
