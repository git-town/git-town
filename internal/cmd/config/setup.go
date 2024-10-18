package config

import (
	"os"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/configfile"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	configInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/config"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
		configStorage:     dialog.ConfigStorageOptionFile,
		normalConfig:      configdomain.DefaultNormalConfig(),
		unvalidatedConfig: configdomain.DefaultUnvalidatedConfig(),
	}
}

func executeConfigSetup(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := loadSetupData(repo, verbose)
	if err != nil || exit {
		return err
	}
	aborted, err := enterData(repo.NormalConfig, repo.UnvalidatedConfig, repo.Git, repo.Backend, &data)
	if err != nil || aborted {
		return err
	}
	err = saveAll(data.userInput, repo.UnvalidatedConfig, repo.Git, repo.Frontend)
	if err != nil {
		return err
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               "setup",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       []gitdomain.BranchName(nil),
		Verbose:               verbose,
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
	configStorage     dialog.ConfigStorageOption
	normalConfig      configdomain.NormalConfig
	unvalidatedConfig configdomain.UnvalidatedConfig
}

func determineHostingPlatform(config config.NormalConfig, userChoice Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userChoice.IsSome() {
		return userChoice
	}
	if originURL, hasOriginURL := config.OriginURL().Get(); hasOriginURL {
		return hosting.Detect(originURL, userChoice)
	}
	return None[configdomain.HostingPlatform]()
}

func enterData(normalConfig config.NormalConfig, unvalidatedConfig config.UnvalidatedConfig, gitCommands git.Commands, backend gitdomain.RunnerQuerier, data *setupData) (aborted bool, err error) {
	aborted, err = dialog.Welcome(data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.Aliases, aborted, err = dialog.Aliases(configdomain.AllAliasableCommands(), normalConfig.Config.Value.Aliases, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	existingMainBranch := unvalidatedConfig.Config.Value.MainBranch
	if existingMainBranch.IsNone() {
		existingMainBranch = gitCommands.DefaultBranch(backend)
	}
	if existingMainBranch.IsNone() {
		existingMainBranch = gitCommands.OriginHead(backend)
	}
	mainBranch, aborted, err := dialog.MainBranch(data.localBranches.Names(), existingMainBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.unvalidatedConfig.MainBranch = Some(mainBranch)
	data.userInput.normalConfig.PerennialBranches, aborted, err = dialog.PerennialBranches(data.localBranches.Names(), normalConfig.Config.Value.PerennialBranches, mainBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.PerennialRegex, aborted, err = dialog.PerennialRegex(normalConfig.Config.Value.PerennialRegex, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.DefaultBranchType, aborted, err = dialog.DefaultBranchType(normalConfig.Config.Value.DefaultBranchType, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.FeatureRegex, aborted, err = dialog.FeatureRegex(normalConfig.Config.Value.FeatureRegex, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.HostingPlatform, aborted, err = dialog.HostingPlatform(normalConfig.Config.Value.HostingPlatform, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	if platform, has := determineHostingPlatform(normalConfig, data.userInput.normalConfig.HostingPlatform).Get(); has {
		switch platform {
		case configdomain.HostingPlatformBitbucket:
			data.userInput.normalConfig.BitbucketUsername, aborted, err = dialog.BitbucketUsername(normalConfig.Config.Value.BitbucketUsername, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
			data.userInput.normalConfig.BitbucketAppPassword, aborted, err = dialog.BitbucketAppPassword(normalConfig.Config.Value.BitbucketAppPassword, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitea:
			data.userInput.normalConfig.GiteaToken, aborted, err = dialog.GiteaToken(normalConfig.Config.Value.GiteaToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitHub:
			data.userInput.normalConfig.GitHubToken, aborted, err = dialog.GitHubToken(normalConfig.Config.Value.GitHubToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitLab:
			data.userInput.normalConfig.GitLabToken, aborted, err = dialog.GitLabToken(normalConfig.Config.Value.GitLabToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		}
	}
	data.userInput.normalConfig.HostingOriginHostname, aborted, err = dialog.OriginHostname(normalConfig.Config.Value.HostingOriginHostname, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.SyncFeatureStrategy, aborted, err = dialog.SyncFeatureStrategy(normalConfig.Config.Value.SyncFeatureStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.SyncPerennialStrategy, aborted, err = dialog.SyncPerennialStrategy(normalConfig.Config.Value.SyncPerennialStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.SyncUpstream, aborted, err = dialog.SyncUpstream(normalConfig.Config.Value.SyncUpstream, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.SyncTags, aborted, err = dialog.SyncTags(normalConfig.Config.Value.SyncTags, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.PushNewBranches, aborted, err = dialog.PushNewBranches(normalConfig.Config.Value.PushNewBranches, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.PushHook, aborted, err = dialog.PushHook(normalConfig.Config.Value.PushHook, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.CreatePrototypeBranches, aborted, err = dialog.CreatePrototypeBranches(normalConfig.Config.Value.CreatePrototypeBranches, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.ShipStrategy, aborted, err = dialog.ShipStrategy(normalConfig.Config.Value.ShipStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.normalConfig.ShipDeleteTrackingBranch, aborted, err = dialog.ShipDeleteTrackingBranch(normalConfig.Config.Value.ShipDeleteTrackingBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.configStorage, aborted, err = dialog.ConfigStorage(data.hasConfigFile, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	return false, nil
}

func loadSetupData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data setupData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	return setupData{
		config:        repo.UnvalidatedConfig,
		dialogInputs:  dialogTestInputs,
		hasConfigFile: repo.UnvalidatedConfig.ConfigFile.IsSome(),
		localBranches: branchesSnapshot.Branches,
		userInput:     defaultUserInput(),
	}, exit, err
}

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, gitCommands git.Commands, frontend gitdomain.Runner) error {
	err := saveAliases(oldConfig.Config.Value.Aliases, userInput.normalConfig.Aliases, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveBitbucketUsername(oldConfig.Config.Value.BitbucketUsername, userInput.normalConfig.BitbucketUsername, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveBitbucketAppPassword(oldConfig.Config.Value.BitbucketAppPassword, userInput.normalConfig.BitbucketAppPassword, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveGiteaToken(oldConfig.Config.Value.GiteaToken, userInput.normalConfig.GiteaToken, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveGitHubToken(oldConfig.Config.Value.GitHubToken, userInput.normalConfig.GitHubToken, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveGitLabToken(oldConfig.Config.Value.GitLabToken, userInput.normalConfig.GitLabToken, gitCommands, frontend)
	if err != nil {
		return err
	}
	switch userInput.configStorage {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, oldConfig)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, oldConfig, gitCommands, frontend)
	}
	panic("unknown configStorage: " + userInput.configStorage)
}

func saveToGit(userInput userInput, oldConfig config.UnvalidatedConfig, gitCommands git.Commands, frontend gitdomain.Runner) error {
	fc := execute.FailureCollector{}
	fc.Check(saveCreatePrototypeBranches(oldConfig.Config.Value.CreatePrototypeBranches, userInput.normalConfig.CreatePrototypeBranches, oldConfig))
	fc.Check(saveHostingPlatform(oldConfig.Config.Value.HostingPlatform, userInput.normalConfig.HostingPlatform, gitCommands, frontend))
	fc.Check(saveOriginHostname(oldConfig.Config.Value.HostingOriginHostname, userInput.normalConfig.HostingOriginHostname, gitCommands, frontend))
	fc.Check(saveMainBranch(oldConfig.Config.Value.MainBranch, userInput.normalConfig.MainBranch.GetOrPanic(), oldConfig))
	fc.Check(savePerennialBranches(oldConfig.Config.Value.PerennialBranches, userInput.normalConfig.PerennialBranches, oldConfig))
	fc.Check(savePerennialRegex(oldConfig.Config.Value.PerennialRegex, userInput.normalConfig.PerennialRegex, oldConfig))
	fc.Check(saveDefaultBranchType(oldConfig.Config.Value.DefaultBranchType, userInput.normalConfig.DefaultBranchType, oldConfig))
	fc.Check(saveFeatureRegex(oldConfig.Config.Value.FeatureRegex, userInput.normalConfig.FeatureRegex, oldConfig))
	fc.Check(savePushHook(oldConfig.Config.Value.PushHook, userInput.normalConfig.PushHook, oldConfig))
	fc.Check(savePushNewBranches(oldConfig.Config.Value.PushNewBranches, userInput.normalConfig.PushNewBranches, oldConfig))
	fc.Check(saveShipStrategy(oldConfig.Config.Value.ShipStrategy, userInput.normalConfig.ShipStrategy, oldConfig))
	fc.Check(saveShipDeleteTrackingBranch(oldConfig.Config.Value.ShipDeleteTrackingBranch, userInput.normalConfig.ShipDeleteTrackingBranch, oldConfig))
	fc.Check(saveSyncFeatureStrategy(oldConfig.Config.Value.SyncFeatureStrategy, userInput.normalConfig.SyncFeatureStrategy, oldConfig))
	fc.Check(saveSyncPerennialStrategy(oldConfig.Config.Value.SyncPerennialStrategy, userInput.normalConfig.SyncPerennialStrategy, oldConfig))
	fc.Check(saveSyncUpstream(oldConfig.Config.Value.SyncUpstream, userInput.normalConfig.SyncUpstream, oldConfig))
	fc.Check(saveSyncTags(oldConfig.Config.Value.SyncTags, userInput.normalConfig.SyncTags, oldConfig))
	return fc.Err
}

func saveAliases(oldAliases, newAliases configdomain.Aliases, gitCommands git.Commands, frontend gitdomain.Runner) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := oldAliases[aliasableCommand]
		newAlias, hasNew := newAliases[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err = gitCommands.RemoveGitAlias(frontend, aliasableCommand)
		case newAlias != oldAlias:
			err = gitCommands.SetGitAlias(frontend, aliasableCommand)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveBitbucketAppPassword(oldPassword, newPassword Option[configdomain.BitbucketAppPassword], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newPassword == oldPassword {
		return nil
	}
	if value, has := newPassword.Get(); has {
		return gitCommands.SetBitbucketAppPassword(frontend, value)
	}
	return gitCommands.RemoveBitbucketAppPassword(frontend)
}

func saveBitbucketUsername(oldValue, newValue Option[configdomain.BitbucketUsername], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitCommands.SetBitbucketUsername(frontend, value)
	}
	return gitCommands.RemoveBitbucketUsername(frontend)
}

func saveCreatePrototypeBranches(oldValue, newValue configdomain.CreatePrototypeBranches, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetCreatePrototypeBranches(newValue)
}

func saveDefaultBranchType(oldValue, newValue configdomain.DefaultBranchType, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetDefaultBranchTypeLocally(newValue)
}

func saveFeatureRegex(oldValue, newValue Option[configdomain.FeatureRegex], config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return config.SetFeatureRegexLocally(value)
	}
	config.RemoveFeatureRegex()
	return nil
}

func saveGiteaToken(oldToken, newToken Option[configdomain.GiteaToken], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken == oldToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGiteaToken(frontend, value)
	}
	return gitCommands.RemoveGiteaToken(frontend)
}

func saveGitHubToken(oldToken, newToken Option[configdomain.GitHubToken], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken == oldToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGitHubToken(frontend, value)
	}
	return gitCommands.RemoveGitHubToken(frontend)
}

func saveGitLabToken(oldToken, newToken Option[configdomain.GitLabToken], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken == oldToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGitLabToken(frontend, value)
	}
	return gitCommands.RemoveGitLabToken(frontend)
}

func saveHostingPlatform(oldHostingPlatform, newHostingPlatform Option[configdomain.HostingPlatform], gitCommands git.Commands, frontend gitdomain.Runner) (err error) {
	oldValue, oldHas := oldHostingPlatform.Get()
	newValue, newHas := newHostingPlatform.Get()
	if !oldHas && !newHas {
		return nil
	}
	if oldValue == newValue {
		return nil
	}
	if newHas {
		return gitCommands.SetHostingPlatform(frontend, newValue)
	}
	return gitCommands.DeleteHostingPlatform(frontend)
}

func saveMainBranch(oldValue Option[gitdomain.LocalBranchName], newValue gitdomain.LocalBranchName, config config.UnvalidatedConfig) error {
	if Some(newValue) == oldValue {
		return nil
	}
	return config.SetMainBranch(newValue)
}

func saveOriginHostname(oldValue, newValue Option[configdomain.HostingOriginHostname], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitCommands.SetOriginHostname(frontend, value)
	}
	return gitCommands.DeleteOriginHostname(frontend)
}

func savePerennialBranches(oldValue, newValue gitdomain.LocalBranchNames, config config.UnvalidatedConfig) error {
	if slices.Compare(oldValue, newValue) != 0 || config.LocalGitConfig.PerennialBranches == nil {
		return config.SetPerennialBranches(newValue)
	}
	return nil
}

func savePerennialRegex(oldValue, newValue Option[configdomain.PerennialRegex], config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return config.SetPerennialRegexLocally(value)
	}
	config.RemovePerennialRegex()
	return nil
}

func savePushHook(oldValue, newValue configdomain.PushHook, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetPushHookLocally(newValue)
}

func savePushNewBranches(oldValue, newValue configdomain.PushNewBranches, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetPushNewBranches(newValue, configdomain.ConfigScopeLocal)
}

func saveShipDeleteTrackingBranch(oldValue, newValue configdomain.ShipDeleteTrackingBranch, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetShipDeleteTrackingBranch(newValue, configdomain.ConfigScopeLocal)
}

func saveShipStrategy(oldValue, newValue configdomain.ShipStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetShipStrategy(newValue, configdomain.ConfigScopeLocal)
}

func saveSyncFeatureStrategy(oldValue, newValue configdomain.SyncFeatureStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetSyncFeatureStrategy(newValue)
}

func saveSyncPerennialStrategy(oldValue, newValue configdomain.SyncPerennialStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetSyncPerennialStrategy(newValue)
}

func saveSyncUpstream(oldValue, newValue configdomain.SyncUpstream, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetSyncUpstream(newValue, configdomain.ConfigScopeLocal)
}

func saveSyncTags(oldValue, newValue configdomain.SyncTags, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.SetSyncTags(newValue)
}

func saveToFile(userInput userInput, config config.UnvalidatedConfig) error {
	err := configfile.Save(&userInput.normalConfig)
	if err != nil {
		return err
	}
	config.RemoveCreatePrototypeBranches()
	config.RemoveMainBranch()
	config.RemovePerennialBranches()
	config.RemovePerennialRegex()
	config.RemovePushNewBranches()
	config.RemovePushHook()
	config.RemoveShipStrategy()
	config.RemoveShipDeleteTrackingBranch()
	config.RemoveSyncFeatureStrategy()
	config.RemoveSyncPerennialStrategy()
	config.RemoveSyncUpstream()
	config.RemoveSyncTags()
	err = saveDefaultBranchType(config.Config.Value.DefaultBranchType, userInput.normalConfig.DefaultBranchType, config)
	if err != nil {
		return err
	}
	return saveFeatureRegex(config.Config.Value.FeatureRegex, userInput.normalConfig.FeatureRegex, config)
}
