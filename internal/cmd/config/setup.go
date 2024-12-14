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
	"github.com/git-town/git-town/v16/internal/config/gitconfig"
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
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeConfigSetup(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

// the config settings to be used if the user accepts all default options
func defaultUserInput(gitAccess gitconfig.Access, gitVersion git.Version) userInput {
	return userInput{
		config:        config.DefaultUnvalidatedConfig(gitAccess, gitVersion),
		configStorage: dialog.ConfigStorageOptionFile,
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
	aborted, err := enterData(repo.UnvalidatedConfig, repo.Git, repo.Backend, &data)
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
		TouchedBranches:       []gitdomain.BranchName{},
		Verbose:               verbose,
	})
}

type setupData struct {
	config        config.UnvalidatedConfig
	dialogInputs  components.TestInputs
	hasConfigFile bool
	localBranches gitdomain.BranchInfos
	remotes       gitdomain.Remotes
	userInput     userInput
}

type userInput struct {
	config        config.UnvalidatedConfig
	configStorage dialog.ConfigStorageOption
}

func determineHostingPlatform(config config.UnvalidatedConfig, userChoice Option[configdomain.HostingPlatform]) Option[configdomain.HostingPlatform] {
	if userChoice.IsSome() {
		return userChoice
	}
	if originURL, hasOriginURL := config.NormalConfig.OriginURL().Get(); hasOriginURL {
		return hosting.Detect(originURL, userChoice)
	}
	return None[configdomain.HostingPlatform]()
}

func enterData(config config.UnvalidatedConfig, gitCommands git.Commands, backend gitdomain.RunnerQuerier, data *setupData) (aborted bool, err error) {
	aborted, err = dialog.Welcome(data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.Aliases, aborted, err = dialog.Aliases(configdomain.AllAliasableCommands(), config.NormalConfig.Aliases, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	existingMainBranch := config.UnvalidatedConfig.MainBranch
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
	data.userInput.config.UnvalidatedConfig.MainBranch = Some(mainBranch)
	data.userInput.config.NormalConfig.PerennialBranches, aborted, err = dialog.PerennialBranches(data.localBranches.Names(), config.NormalConfig.PerennialBranches, mainBranch, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.PerennialRegex, aborted, err = dialog.PerennialRegex(config.NormalConfig.PerennialRegex, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.DefaultBranchType, aborted, err = dialog.DefaultBranchType(config.NormalConfig.DefaultBranchType, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.FeatureRegex, aborted, err = dialog.FeatureRegex(config.NormalConfig.FeatureRegex, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	// TODO: make this dialog work, then add the step to enter the dev-remote to all E2E tests for the setup assistant.
	data.userInput.config.NormalConfig.DevRemote, aborted, err = dialog.DevRemote(config.NormalConfig.DevRemote, data.remotes, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.HostingPlatform, aborted, err = dialog.HostingPlatform(config.NormalConfig.HostingPlatform, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	if platform, has := determineHostingPlatform(config, data.userInput.config.NormalConfig.HostingPlatform).Get(); has {
		switch platform {
		case configdomain.HostingPlatformBitbucket, configdomain.HostingPlatformBitbucketDatacenter:
			data.userInput.config.NormalConfig.BitbucketUsername, aborted, err = dialog.BitbucketUsername(config.NormalConfig.BitbucketUsername, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
			data.userInput.config.NormalConfig.BitbucketAppPassword, aborted, err = dialog.BitbucketAppPassword(config.NormalConfig.BitbucketAppPassword, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitea:
			data.userInput.config.NormalConfig.GiteaToken, aborted, err = dialog.GiteaToken(config.NormalConfig.GiteaToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitHub:
			data.userInput.config.NormalConfig.GitHubToken, aborted, err = dialog.GitHubToken(config.NormalConfig.GitHubToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		case configdomain.HostingPlatformGitLab:
			data.userInput.config.NormalConfig.GitLabToken, aborted, err = dialog.GitLabToken(config.NormalConfig.GitLabToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, err
			}
		}
	}
	data.userInput.config.NormalConfig.HostingOriginHostname, aborted, err = dialog.OriginHostname(config.NormalConfig.HostingOriginHostname, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.SyncFeatureStrategy, aborted, err = dialog.SyncFeatureStrategy(config.NormalConfig.SyncFeatureStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.SyncPerennialStrategy, aborted, err = dialog.SyncPerennialStrategy(config.NormalConfig.SyncPerennialStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.SyncPrototypeStrategy, aborted, err = dialog.SyncPrototypeStrategy(config.NormalConfig.SyncPrototypeStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.SyncUpstream, aborted, err = dialog.SyncUpstream(config.NormalConfig.SyncUpstream, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.SyncTags, aborted, err = dialog.SyncTags(config.NormalConfig.SyncTags, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.PushNewBranches, aborted, err = dialog.PushNewBranches(config.NormalConfig.PushNewBranches, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.PushHook, aborted, err = dialog.PushHook(config.NormalConfig.PushHook, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.NewBranchType, aborted, err = dialog.NewBranchType(config.NormalConfig.NewBranchType, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.ShipStrategy, aborted, err = dialog.ShipStrategy(config.NormalConfig.ShipStrategy, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, err
	}
	data.userInput.config.NormalConfig.ShipDeleteTrackingBranch, aborted, err = dialog.ShipDeleteTrackingBranch(config.NormalConfig.ShipDeleteTrackingBranch, data.dialogInputs.Next())
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
	if err != nil {
		return data, exit, err
	}
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, exit, err
	}
	return setupData{
		config:        repo.UnvalidatedConfig,
		dialogInputs:  dialogTestInputs,
		hasConfigFile: repo.UnvalidatedConfig.NormalConfig.ConfigFile.IsSome(),
		localBranches: branchesSnapshot.Branches,
		remotes:       remotes,
		userInput:     defaultUserInput(repo.UnvalidatedConfig.NormalConfig.GitConfigAccess, repo.UnvalidatedConfig.NormalConfig.GitVersion),
	}, exit, nil
}

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, gitCommands git.Commands, frontend gitdomain.Runner) error {
	err := saveAliases(oldConfig.NormalConfig.Aliases, userInput.config.NormalConfig.Aliases, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveBitbucketUsername(oldConfig.NormalConfig.BitbucketUsername, userInput.config.NormalConfig.BitbucketUsername, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveBitbucketAppPassword(oldConfig.NormalConfig.BitbucketAppPassword, userInput.config.NormalConfig.BitbucketAppPassword, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveGiteaToken(oldConfig.NormalConfig.GiteaToken, userInput.config.NormalConfig.GiteaToken, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveGitHubToken(oldConfig.NormalConfig.GitHubToken, userInput.config.NormalConfig.GitHubToken, gitCommands, frontend)
	if err != nil {
		return err
	}
	err = saveGitLabToken(oldConfig.NormalConfig.GitLabToken, userInput.config.NormalConfig.GitLabToken, gitCommands, frontend)
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
	fc.Check(saveNewBranchType(oldConfig.NormalConfig.NewBranchType, userInput.config.NormalConfig.NewBranchType, oldConfig))
	fc.Check(saveHostingPlatform(oldConfig.NormalConfig.HostingPlatform, userInput.config.NormalConfig.HostingPlatform, gitCommands, frontend))
	fc.Check(saveOriginHostname(oldConfig.NormalConfig.HostingOriginHostname, userInput.config.NormalConfig.HostingOriginHostname, gitCommands, frontend))
	fc.Check(saveMainBranch(oldConfig.UnvalidatedConfig.MainBranch, userInput.config.UnvalidatedConfig.MainBranch.GetOrPanic(), oldConfig))
	fc.Check(savePerennialBranches(oldConfig.NormalConfig.PerennialBranches, userInput.config.NormalConfig.PerennialBranches, oldConfig))
	fc.Check(savePerennialRegex(oldConfig.NormalConfig.PerennialRegex, userInput.config.NormalConfig.PerennialRegex, oldConfig))
	fc.Check(saveDefaultBranchType(oldConfig.NormalConfig.DefaultBranchType, userInput.config.NormalConfig.DefaultBranchType, oldConfig))
	fc.Check(saveDevRemote(oldConfig.NormalConfig.DevRemote, userInput.config.NormalConfig.DevRemote, oldConfig))
	fc.Check(saveFeatureRegex(oldConfig.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, oldConfig))
	fc.Check(savePushHook(oldConfig.NormalConfig.PushHook, userInput.config.NormalConfig.PushHook, oldConfig))
	fc.Check(savePushNewBranches(oldConfig.NormalConfig.PushNewBranches, userInput.config.NormalConfig.PushNewBranches, oldConfig))
	fc.Check(saveShipStrategy(oldConfig.NormalConfig.ShipStrategy, userInput.config.NormalConfig.ShipStrategy, oldConfig))
	fc.Check(saveShipDeleteTrackingBranch(oldConfig.NormalConfig.ShipDeleteTrackingBranch, userInput.config.NormalConfig.ShipDeleteTrackingBranch, oldConfig))
	fc.Check(saveSyncFeatureStrategy(oldConfig.NormalConfig.SyncFeatureStrategy, userInput.config.NormalConfig.SyncFeatureStrategy, oldConfig))
	fc.Check(saveSyncPerennialStrategy(oldConfig.NormalConfig.SyncPerennialStrategy, userInput.config.NormalConfig.SyncPerennialStrategy, oldConfig))
	fc.Check(saveSyncPrototypeStrategy(oldConfig.NormalConfig.SyncPrototypeStrategy, userInput.config.NormalConfig.SyncPrototypeStrategy, oldConfig))
	fc.Check(saveSyncUpstream(oldConfig.NormalConfig.SyncUpstream, userInput.config.NormalConfig.SyncUpstream, oldConfig))
	fc.Check(saveSyncTags(oldConfig.NormalConfig.SyncTags, userInput.config.NormalConfig.SyncTags, oldConfig))
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

func saveNewBranchType(oldValue, newValue configdomain.BranchType, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetNewBranchType(newValue)
}

func saveDefaultBranchType(oldValue, newValue configdomain.BranchType, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetDefaultBranchTypeLocally(newValue)
}

func saveDevRemote(oldValue, newValue gitdomain.Remote, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetDevRemote(newValue)
}

func saveFeatureRegex(oldValue, newValue Option[configdomain.FeatureRegex], config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return config.NormalConfig.SetFeatureRegexLocally(value)
	}
	config.NormalConfig.RemoveFeatureRegex()
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
	if slices.Compare(oldValue, newValue) != 0 || config.NormalConfig.LocalGitConfig.PerennialBranches == nil {
		return config.NormalConfig.SetPerennialBranches(newValue)
	}
	return nil
}

func savePerennialRegex(oldValue, newValue Option[configdomain.PerennialRegex], config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	if value, has := newValue.Get(); has {
		return config.NormalConfig.SetPerennialRegexLocally(value)
	}
	config.NormalConfig.RemovePerennialRegex()
	return nil
}

func savePushHook(oldValue, newValue configdomain.PushHook, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetPushHookLocally(newValue)
}

func savePushNewBranches(oldValue, newValue configdomain.PushNewBranches, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetPushNewBranches(newValue, configdomain.ConfigScopeLocal)
}

func saveShipDeleteTrackingBranch(oldValue, newValue configdomain.ShipDeleteTrackingBranch, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetShipDeleteTrackingBranch(newValue, configdomain.ConfigScopeLocal)
}

func saveShipStrategy(oldValue, newValue configdomain.ShipStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetShipStrategy(newValue, configdomain.ConfigScopeLocal)
}

func saveSyncFeatureStrategy(oldValue, newValue configdomain.SyncFeatureStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetSyncFeatureStrategy(newValue)
}

func saveSyncPerennialStrategy(oldValue, newValue configdomain.SyncPerennialStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetSyncPerennialStrategy(newValue)
}

func saveSyncPrototypeStrategy(oldValue, newValue configdomain.SyncPrototypeStrategy, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetSyncPrototypeStrategy(newValue)
}

func saveSyncUpstream(oldValue, newValue configdomain.SyncUpstream, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetSyncUpstream(newValue, configdomain.ConfigScopeLocal)
}

func saveSyncTags(oldValue, newValue configdomain.SyncTags, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetSyncTags(newValue)
}

func saveToFile(userInput userInput, config config.UnvalidatedConfig) error {
	err := configfile.Save(&userInput.config)
	if err != nil {
		return err
	}
	config.NormalConfig.RemoveCreatePrototypeBranches()
	config.NormalConfig.RemoveDevRemote()
	config.RemoveMainBranch()
	config.NormalConfig.RemoveNewBranchType()
	config.NormalConfig.RemovePerennialBranches()
	config.NormalConfig.RemovePerennialRegex()
	config.NormalConfig.RemovePushNewBranches()
	config.NormalConfig.RemovePushHook()
	config.NormalConfig.RemoveShipStrategy()
	config.NormalConfig.RemoveShipDeleteTrackingBranch()
	config.NormalConfig.RemoveSyncFeatureStrategy()
	config.NormalConfig.RemoveSyncPerennialStrategy()
	config.NormalConfig.RemoveSyncPrototypeStrategy()
	config.NormalConfig.RemoveSyncUpstream()
	config.NormalConfig.RemoveSyncTags()
	err = saveDefaultBranchType(config.NormalConfig.DefaultBranchType, userInput.config.NormalConfig.DefaultBranchType, config)
	if err != nil {
		return err
	}
	return saveFeatureRegex(config.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, config)
}
