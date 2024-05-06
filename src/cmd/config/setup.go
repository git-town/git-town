package config

import (
	"os"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
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
		config:        configdomain.DefaultConfig(),
		configStorage: dialog.ConfigStorageOptionFile,
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
	data, exit, err := loadSetupData(repo)
	if err != nil || exit {
		return err
	}
	aborted, err := enterData(repo.UnvalidatedConfig, repo.Backend, data)
	if err != nil || aborted {
		return err
	}
	err = saveAll(data.userInput, repo.UnvalidatedConfig, repo.Frontend)
	if err != nil {
		return err
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             repo.Backend,
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "setup",
		CommandsCounter:     repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       &repo.FinalMessages,
		RootDir:             repo.RootDir,
		Verbose:             verbose,
	})
}

type setupData struct {
	config        config.UnvalidatedConfig
	dialogInputs  components.TestInputs
	hasConfigFile bool
	localBranches gitdomain.BranchInfos
	userInput     userInput
}

type userInput struct {
	config        configdomain.UnvalidatedConfig
	configStorage dialog.ConfigStorageOption
}

func determineHostingPlatform(config config.UnvalidatedConfig, userChoice Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userChoice.IsSome() {
		return userChoice
	}
	if originURL, hasOriginURL := config.OriginURL().Get(); hasOriginURL {
		return hosting.Detect(originURL, userChoice)
	}
	return None[configdomain.HostingPlatform]()
}

func enterData(config config.UnvalidatedConfig, backend git.BackendCommands, data *setupData) (aborted bool, err error) {
	aborted, err = dialog.Welcome(data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.Aliases, aborted, err = dialog.Aliases(configdomain.AllAliasableCommands(), config.Config.Aliases, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	existingMainBranch := config.Config.MainBranch
	if existingMainBranch.IsNone() {
		existingMainBranch = backend.DefaultBranch()
	}
	if existingMainBranch.IsNone() {
		existingMainBranch = backend.OriginHead()
	}
	mainBranch, aborted, err := dialog.MainBranch(data.localBranches.Names(), existingMainBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.MainBranch = Some(mainBranch)
	data.userInput.config.PerennialBranches, aborted, err = dialog.PerennialBranches(data.localBranches.Names(), config.Config.PerennialBranches, mainBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.PerennialRegex, aborted, err = dialog.PerennialRegex(config.Config.PerennialRegex, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.HostingPlatform, aborted, err = dialog.HostingPlatform(config.Config.HostingPlatform, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	if platform, has := determineHostingPlatform(config, data.userInput.config.HostingPlatform).Get(); has {
		switch platform {
		case configdomain.HostingPlatformBitbucket:
			// BitBucket API isn't supported yet
		case configdomain.HostingPlatformGitea:
			data.userInput.config.GiteaToken, aborted, err = dialog.GiteaToken(config.Config.GiteaToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitHub:
			data.userInput.config.GitHubToken, aborted, err = dialog.GitHubToken(config.Config.GitHubToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitLab:
			data.userInput.config.GitLabToken, aborted, err = dialog.GitLabToken(config.Config.GitLabToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		}
	}
	data.userInput.config.HostingOriginHostname, aborted, err = dialog.OriginHostname(config.Config.HostingOriginHostname, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.SyncFeatureStrategy, aborted, err = dialog.SyncFeatureStrategy(config.Config.SyncFeatureStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.SyncPerennialStrategy, aborted, err = dialog.SyncPerennialStrategy(config.Config.SyncPerennialStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.SyncUpstream, aborted, err = dialog.SyncUpstream(config.Config.SyncUpstream, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.PushNewBranches, aborted, err = dialog.PushNewBranches(config.Config.PushNewBranches, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.PushHook, aborted, err = dialog.PushHook(config.Config.PushHook, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.SyncBeforeShip, aborted, err = dialog.SyncBeforeShip(config.Config.SyncBeforeShip, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.ShipDeleteTrackingBranch, aborted, err = dialog.ShipDeleteTrackingBranch(config.Config.ShipDeleteTrackingBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.configStorage, aborted, err = dialog.ConfigStorage(data.hasConfigFile, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return false, nil
}

func loadSetupData(repo execute.OpenRepoResult) (setupData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               &repo.Backend,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		Frontend:              &repo.Frontend,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
	})
	return &setupData{
		config:        repo.UnvalidatedConfig,
		dialogInputs:  dialogTestInputs,
		hasConfigFile: repo.UnvalidatedConfig.ConfigFile.IsSome(),
		localBranches: branchesSnapshot.Branches,
		userInput:     defaultUserInput(),
	}, exit, err
}

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, frontend git.FrontendCommands) error {
	err := saveAliases(oldConfig.Config.Aliases, userInput.config.Aliases, frontend)
	if err != nil {
		return err
	}
	err = saveGiteaToken(oldConfig.Config.GiteaToken, userInput.config.GiteaToken, frontend)
	if err != nil {
		return err
	}
	err = saveGitHubToken(oldConfig.Config.GitHubToken, userInput.config.GitHubToken, frontend)
	if err != nil {
		return err
	}
	err = saveGitLabToken(oldConfig.Config.GitLabToken, userInput.config.GitLabToken, frontend)
	if err != nil {
		return err
	}
	switch userInput.configStorage {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, &oldConfig)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, oldConfig, frontend)
	}
	panic("unknown configStorage: " + userInput.configStorage)
}

func saveToGit(userInput userInput, oldConfig config.UnvalidatedConfig, frontend git.FrontendCommands) error {
	fc := execute.FailureCollector{}
	fc.Check(saveHostingPlatform(oldConfig.Config.HostingPlatform, userInput.config.HostingPlatform, frontend))
	fc.Check(saveOriginHostname(oldConfig.Config.HostingOriginHostname, userInput.config.HostingOriginHostname, frontend))
	fc.Check(saveMainBranch(oldConfig.Config.MainBranch, userInput.config.MainBranch.GetOrPanic(), oldConfig))
	fc.Check(savePerennialBranches(oldConfig.Config.PerennialBranches, userInput.config.PerennialBranches, oldConfig))
	fc.Check(savePerennialRegex(oldConfig.Config.PerennialRegex, userInput.config.PerennialRegex, oldConfig))
	fc.Check(savePushHook(oldConfig.Config.PushHook, userInput.config.PushHook, oldConfig))
	fc.Check(savePushNewBranches(oldConfig.Config.PushNewBranches, userInput.config.PushNewBranches, oldConfig))
	fc.Check(saveShipDeleteTrackingBranch(oldConfig.Config.ShipDeleteTrackingBranch, userInput.config.ShipDeleteTrackingBranch, oldConfig))
	fc.Check(saveSyncFeatureStrategy(oldConfig.Config.SyncFeatureStrategy, userInput.config.SyncFeatureStrategy, oldConfig))
	fc.Check(saveSyncPerennialStrategy(oldConfig.Config.SyncPerennialStrategy, userInput.config.SyncPerennialStrategy, oldConfig))
	fc.Check(saveSyncUpstream(oldConfig.Config.SyncUpstream, userInput.config.SyncUpstream, oldConfig))
	fc.Check(saveSyncBeforeShip(oldConfig.Config.SyncBeforeShip, userInput.config.SyncBeforeShip, oldConfig))
	return fc.Err
}

func saveAliases(oldAliases, newAliases configdomain.Aliases, frontend git.FrontendCommands) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := oldAliases[aliasableCommand]
		newAlias, hasNew := newAliases[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err = frontend.RemoveGitAlias(aliasableCommand)
		case newAlias != oldAlias:
			err = frontend.SetGitAlias(aliasableCommand)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveGiteaToken(oldToken, newToken Option[configdomain.GiteaToken], frontend git.FrontendCommands) error {
	if newToken == oldToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return frontend.SetGiteaToken(value)
	}
	return frontend.RemoveGiteaToken()
}

func saveGitHubToken(oldToken, newToken Option[configdomain.GitHubToken], frontend git.FrontendCommands) error {
	if newToken == oldToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return frontend.SetGitHubToken(value)
	}
	return frontend.RemoveGitHubToken()
}

func saveGitLabToken(oldToken, newToken Option[configdomain.GitLabToken], frontend git.FrontendCommands) error {
	if newToken == oldToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return frontend.SetGitLabToken(value)
	}
	return frontend.RemoveGitLabToken()
}

func saveHostingPlatform(oldHostingPlatform, newHostingPlatform Option[configdomain.HostingPlatform], frontend git.FrontendCommands) (err error) {
	oldValue, oldHas := oldHostingPlatform.Get()
	newValue, newHas := newHostingPlatform.Get()
	if !oldHas && !newHas {
		return nil
	}
	if oldValue == newValue {
		return nil
	}
	if newHas {
		return frontend.SetHostingPlatform(newValue)
	}
	return frontend.DeleteHostingPlatform()
}

func saveMainBranch(oldValue Option[gitdomain.LocalBranchName], newValue gitdomain.LocalBranchName, oldConfig config.UnvalidatedConfig) error {
	if Some(newValue) == oldValue {
		return nil
	}
	return oldConfig.SetMainBranch(newValue)
}

func saveOriginHostname(oldValue, newValue Option[configdomain.HostingOriginHostname], frontend git.FrontendCommands) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return frontend.SetOriginHostname(value)
	}
	return frontend.DeleteOriginHostname()
}

func savePerennialBranches(oldValue, newValue gitdomain.LocalBranchNames, oldConfig config.UnvalidatedConfig) error {
	if slices.Compare(oldValue, newValue) != 0 || oldConfig.LocalGitConfig.PerennialBranches == nil {
		return oldConfig.SetPerennialBranches(newValue)
	}
	return nil
}

func savePerennialRegex(oldValue, newValue Option[configdomain.PerennialRegex], oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return oldConfig.SetPerennialRegexLocally(value)
	}
	oldConfig.RemovePerennialRegex()
	return nil
}

func savePushHook(oldValue, newValue configdomain.PushHook, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetPushHookLocally(newValue)
}

func savePushNewBranches(oldValue, newValue configdomain.PushNewBranches, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetPushNewBranches(newValue, false)
}

func saveShipDeleteTrackingBranch(oldValue, newValue configdomain.ShipDeleteTrackingBranch, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetShipDeleteTrackingBranch(newValue, false)
}

func saveSyncFeatureStrategy(oldValue, newValue configdomain.SyncFeatureStrategy, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetSyncFeatureStrategy(newValue)
}

func saveSyncPerennialStrategy(oldValue, newValue configdomain.SyncPerennialStrategy, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetSyncPerennialStrategy(newValue)
}

func saveSyncUpstream(oldValue, newValue configdomain.SyncUpstream, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetSyncUpstream(newValue, false)
}

func saveSyncBeforeShip(oldValue, newValue configdomain.SyncBeforeShip, oldConfig config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return oldConfig.SetSyncBeforeShip(newValue, false)
}

func saveToFile(userInput userInput, config *config.UnvalidatedConfig) error {
	err := configfile.Save(&userInput.config)
	if err != nil {
		return err
	}
	config.RemoveMainBranch()
	config.RemovePerennialBranches()
	config.RemovePerennialRegex()
	config.RemovePushNewBranches()
	config.RemovePushHook()
	config.RemoveSyncBeforeShip()
	config.RemoveShipDeleteTrackingBranch()
	config.RemoveSyncFeatureStrategy()
	config.RemoveSyncPerennialStrategy()
	config.RemoveSyncUpstream()
	return nil
}
