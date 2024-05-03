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
	"github.com/git-town/git-town/v14/src/validate"
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
	data, exit, err := loadSetupData(repo, verbose)
	if err != nil || exit {
		return err
	}
	aborted, err := enterData(repo.UnvalidatedConfig, repo.Backend, data)
	if err != nil || aborted {
		return err
	}
	err = saveAll(data.runner, data.userInput)
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
	config        config.ValidatedConfig
	dialogInputs  components.TestInputs
	hasConfigFile bool
	localBranches gitdomain.BranchInfos
	runner        *git.ProdRunner
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

func loadSetupData(repo *execute.OpenRepoResult, verbose bool) (*setupData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               &repo.Backend,
		Config:                repo.UnvalidatedConfig.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		Frontend:              &repo.Frontend,
		HandleUnfinishedState: false,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil {
		return nil, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, runner, aborted, err := validate.Config(validate.ConfigArgs{
		Backend:            &repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		CommandsCounter:    repo.CommandsCounter,
		ConfigSnapshot:     repo.ConfigSnapshot,
		DialogTestInputs:   dialogTestInputs,
		FinalMessages:      &repo.FinalMessages,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		RootDir:            repo.RootDir,
		StashSize:          stashSize,
		TestInputs:         &dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
		Verbose:            verbose,
	})
	if err != nil || aborted {
		return nil, aborted, err
	}
	return &setupData{
		config:        *validatedConfig,
		dialogInputs:  dialogTestInputs,
		hasConfigFile: repo.UnvalidatedConfig.ConfigFile.IsSome(),
		localBranches: branchesSnapshot.Branches,
		runner:        runner,
		userInput:     defaultUserInput(),
	}, exit, err
}

func saveAll(runner *git.ProdRunner, userInput userInput) error {
	err := saveAliases(runner, userInput.config.Aliases)
	if err != nil {
		return err
	}
	err = saveGiteaToken(runner, userInput.config.GiteaToken)
	if err != nil {
		return err
	}
	err = saveGitHubToken(runner, userInput.config.GitHubToken)
	if err != nil {
		return err
	}
	err = saveGitLabToken(runner, userInput.config.GitLabToken)
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
	err := saveHostingPlatform(runner, userInput.config.HostingPlatform)
	if err != nil {
		return err
	}
	err = saveOriginHostname(runner, userInput.config.HostingOriginHostname)
	if err != nil {
		return err
	}
	err = saveMainBranch(runner, userInput.config.MainBranch.GetOrPanic())
	if err != nil {
		return err
	}
	err = savePerennialBranches(runner, userInput.config.PerennialBranches)
	if err != nil {
		return err
	}
	err = savePerennialRegex(runner, userInput.config.PerennialRegex)
	if err != nil {
		return err
	}
	err = savePushHook(runner, userInput.config.PushHook)
	if err != nil {
		return err
	}
	err = savePushNewBranches(runner, userInput.config.PushNewBranches)
	if err != nil {
		return err
	}
	err = saveShipDeleteTrackingBranch(runner, userInput.config.ShipDeleteTrackingBranch)
	if err != nil {
		return err
	}
	err = saveSyncFeatureStrategy(runner, userInput.config.SyncFeatureStrategy)
	if err != nil {
		return err
	}
	err = saveSyncPerennialStrategy(runner, userInput.config.SyncPerennialStrategy)
	if err != nil {
		return err
	}
	err = saveSyncUpstream(runner, userInput.config.SyncUpstream)
	if err != nil {
		return err
	}
	err = saveSyncBeforeShip(runner, userInput.config.SyncBeforeShip)
	if err != nil {
		return err
	}
	return nil
}

func saveAliases(runner *git.ProdRunner, newAliases configdomain.Aliases) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := runner.Config.Config.Aliases[aliasableCommand]
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
	if newToken == runner.Config.Config.GiteaToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return runner.Frontend.SetGiteaToken(value)
	}
	return runner.Frontend.RemoveGiteaToken()
}

func saveGitHubToken(runner *git.ProdRunner, newToken Option[configdomain.GitHubToken]) error {
	if newToken == runner.Config.Config.GitHubToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return runner.Frontend.SetGitHubToken(value)
	}
	return runner.Frontend.RemoveGitHubToken()
}

func saveGitLabToken(runner *git.ProdRunner, newToken Option[configdomain.GitLabToken]) error {
	if newToken == runner.Config.Config.GitLabToken {
		return nil
	}
	if value, has := newToken.Get(); has {
		return runner.Frontend.SetGitLabToken(value)
	}
	return runner.Frontend.RemoveGitLabToken()
}

func saveHostingPlatform(runner *git.ProdRunner, newHostingPlatform Option[configdomain.HostingPlatform]) (err error) {
	oldValue, oldHas := runner.Config.Config.HostingPlatform.Get()
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
	if newValue == runner.Config.Config.MainBranch {
		return nil
	}
	return runner.Config.SetMainBranch(newValue)
}

func saveOriginHostname(runner *git.ProdRunner, newValue Option[configdomain.HostingOriginHostname]) error {
	if newValue == runner.Config.Config.HostingOriginHostname {
		return nil
	}
	if value, has := newValue.Get(); has {
		return runner.Frontend.SetOriginHostname(value)
	}
	return runner.Frontend.DeleteOriginHostname()
}

func savePerennialBranches(runner *git.ProdRunner, newValue gitdomain.LocalBranchNames) error {
	oldValue := runner.Config.Config.PerennialBranches
	if slices.Compare(oldValue, newValue) != 0 || runner.Config.LocalGitConfig.PerennialBranches == nil {
		return runner.Config.SetPerennialBranches(newValue)
	}
	return nil
}

func savePerennialRegex(runner *git.ProdRunner, newValue Option[configdomain.PerennialRegex]) error {
	if newValue == runner.Config.Config.PerennialRegex {
		return nil
	}
	if value, has := newValue.Get(); has {
		return runner.Config.SetPerennialRegexLocally(value)
	}
	runner.Config.RemovePerennialRegex()
	return nil
}

func savePushHook(runner *git.ProdRunner, newValue configdomain.PushHook) error {
	if newValue == runner.Config.Config.PushHook {
		return nil
	}
	return runner.Config.SetPushHookLocally(newValue)
}

func savePushNewBranches(runner *git.ProdRunner, newValue configdomain.PushNewBranches) error {
	if newValue == runner.Config.Config.PushNewBranches {
		return nil
	}
	return runner.Config.SetPushNewBranches(newValue, false)
}

func saveShipDeleteTrackingBranch(runner *git.ProdRunner, newValue configdomain.ShipDeleteTrackingBranch) error {
	if newValue == runner.Config.Config.ShipDeleteTrackingBranch {
		return nil
	}
	return runner.Config.SetShipDeleteTrackingBranch(newValue, false)
}

func saveSyncFeatureStrategy(runner *git.ProdRunner, newValue configdomain.SyncFeatureStrategy) error {
	if newValue == runner.Config.Config.SyncFeatureStrategy {
		return nil
	}
	return runner.Config.SetSyncFeatureStrategy(newValue)
}

func saveSyncPerennialStrategy(runner *git.ProdRunner, newValue configdomain.SyncPerennialStrategy) error {
	if newValue == runner.Config.Config.SyncPerennialStrategy {
		return nil
	}
	return runner.Config.SetSyncPerennialStrategy(newValue)
}

func saveSyncUpstream(runner *git.ProdRunner, newValue configdomain.SyncUpstream) error {
	if newValue == runner.Config.Config.SyncUpstream {
		return nil
	}
	return runner.Config.SetSyncUpstream(newValue, false)
}

func saveSyncBeforeShip(runner *git.ProdRunner, newValue configdomain.SyncBeforeShip) error {
	if newValue == runner.Config.Config.SyncBeforeShip {
		return nil
	}
	return runner.Config.SetSyncBeforeShip(newValue, false)
}

func saveToFile(userInput userInput, runner *git.ProdRunner) error {
	err := configfile.Save(&userInput.config)
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
