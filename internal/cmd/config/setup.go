package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
	aborted, tokenScope, forgeTypeOpt, err := enterData(repo.UnvalidatedConfig, repo.Git, repo.Backend, &data)
	if err != nil || aborted {
		return err
	}
	err = saveAll(data.userInput, repo.UnvalidatedConfig, data.configFile, tokenScope, forgeTypeOpt, repo.Git, repo.Frontend)
	if err != nil {
		return err
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
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
	configFile    Option[configdomain.PartialConfig]
	dialogInputs  components.TestInputs
	localBranches gitdomain.BranchInfos
	remotes       gitdomain.Remotes
	userInput     userInput
}

type userInput struct {
	config        config.UnvalidatedConfig
	configStorage dialog.ConfigStorageOption
}

func determineHostingPlatform(config config.UnvalidatedConfig, userChoice Option[forgedomain.ForgeType]) Option[forgedomain.ForgeType] {
	if userChoice.IsSome() {
		return userChoice
	}
	if devURL, hasDevURL := config.NormalConfig.DevURL().Get(); hasDevURL {
		return forge.Detect(devURL, userChoice)
	}
	return None[forgedomain.ForgeType]()
}

func enterData(config config.UnvalidatedConfig, gitCommands git.Commands, backend gitdomain.RunnerQuerier, data *setupData) (aborted bool, tokenScope configdomain.ConfigScope, forgeTypeOpt Option[forgedomain.ForgeType], err error) {
	tokenScope = configdomain.ConfigScopeLocal
	configFile := data.configFile.GetOrDefault()
	aborted, err = dialog.Welcome(data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, tokenScope, None[forgedomain.ForgeType](), err
	}
	data.userInput.config.NormalConfig.Aliases, aborted, err = dialog.Aliases(configdomain.AllAliasableCommands(), config.NormalConfig.Aliases, data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, tokenScope, None[forgedomain.ForgeType](), err
	}
	var mainBranch gitdomain.LocalBranchName
	if configFileMainBranch, configFileHasMainBranch := configFile.MainBranch.Get(); configFileHasMainBranch {
		mainBranch = configFileMainBranch
	} else {
		existingMainBranch := config.UnvalidatedConfig.MainBranch
		if existingMainBranch.IsNone() {
			existingMainBranch = gitCommands.DefaultBranch(backend)
		}
		if existingMainBranch.IsNone() {
			existingMainBranch = gitCommands.OriginHead(backend)
		}
		mainBranch, aborted, err = dialog.MainBranch(data.localBranches.Names(), existingMainBranch, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
		data.userInput.config.UnvalidatedConfig.MainBranch = Some(mainBranch)
	}
	if len(configFile.PerennialBranches) == 0 {
		data.userInput.config.NormalConfig.PerennialBranches, aborted, err = dialog.PerennialBranches(data.localBranches.Names(), config.NormalConfig.PerennialBranches, mainBranch, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.PerennialRegex.IsNone() {
		data.userInput.config.NormalConfig.PerennialRegex, aborted, err = dialog.PerennialRegex(config.NormalConfig.PerennialRegex, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.FeatureRegex.IsNone() {
		data.userInput.config.NormalConfig.FeatureRegex, aborted, err = dialog.FeatureRegex(config.NormalConfig.FeatureRegex, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.UnknownBranchType.IsNone() {
		data.userInput.config.NormalConfig.UnknownBranchType, aborted, err = dialog.UnknownBranchType(config.NormalConfig.UnknownBranchType, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.DevRemote.IsNone() {
		data.userInput.config.NormalConfig.DevRemote, aborted, err = dialog.DevRemote(config.NormalConfig.DevRemote, data.remotes, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.ForgeType.IsNone() {
		data.userInput.config.NormalConfig.ForgeType, aborted, err = dialog.ForgeType(config.NormalConfig.ForgeType, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	forgeTypeOpt = determineHostingPlatform(config, data.userInput.config.NormalConfig.ForgeType)
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			data.userInput.config.NormalConfig.BitbucketUsername, aborted, err = dialog.BitbucketUsername(config.NormalConfig.BitbucketUsername, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, tokenScope, None[forgedomain.ForgeType](), err
			}
			data.userInput.config.NormalConfig.BitbucketAppPassword, aborted, err = dialog.BitbucketAppPassword(config.NormalConfig.BitbucketAppPassword, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, tokenScope, None[forgedomain.ForgeType](), err
			}
			if showScopeDialog(data.userInput.config.NormalConfig.BitbucketUsername, config.NormalConfig.BitbucketUsername) &&
				showScopeDialog(data.userInput.config.NormalConfig.BitbucketAppPassword, config.NormalConfig.BitbucketAppPassword) {
				scope := determineScope(config.NormalConfig.GitConfig.BitbucketAppPassword)
				tokenScope, aborted, err = dialog.TokenScope(scope, data.dialogInputs.Next())
				if err != nil || aborted {
					return aborted, tokenScope, None[forgedomain.ForgeType](), err
				}
			}
		case forgedomain.ForgeTypeCodeberg:
			data.userInput.config.NormalConfig.CodebergToken, aborted, err = dialog.CodebergToken(config.NormalConfig.CodebergToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, tokenScope, None[forgedomain.ForgeType](), err
			}
			if showScopeDialog(data.userInput.config.NormalConfig.CodebergToken, config.NormalConfig.CodebergToken) {
				scope := determineScope(config.NormalConfig.GitConfig.CodebergToken)
				tokenScope, aborted, err = dialog.TokenScope(scope, data.dialogInputs.Next())
				if err != nil || aborted {
					return aborted, tokenScope, None[forgedomain.ForgeType](), err
				}
			}
		case forgedomain.ForgeTypeGitea:
			data.userInput.config.NormalConfig.GiteaToken, aborted, err = dialog.GiteaToken(config.NormalConfig.GiteaToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, tokenScope, None[forgedomain.ForgeType](), err
			}
			if showScopeDialog(data.userInput.config.NormalConfig.GiteaToken, config.NormalConfig.GiteaToken) {
				scope := determineScope(config.NormalConfig.GitConfig.GiteaToken)
				tokenScope, aborted, err = dialog.TokenScope(scope, data.dialogInputs.Next())
				if err != nil || aborted {
					return aborted, tokenScope, None[forgedomain.ForgeType](), err
				}
			}
		case forgedomain.ForgeTypeGitHub:
			data.userInput.config.NormalConfig.GitHubToken, aborted, err = dialog.GitHubToken(config.NormalConfig.GitHubToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, tokenScope, None[forgedomain.ForgeType](), err
			}
			if showScopeDialog(data.userInput.config.NormalConfig.GitHubToken, config.NormalConfig.GitHubToken) {
				scope := determineScope(config.NormalConfig.GitConfig.GitHubToken)
				tokenScope, aborted, err = dialog.TokenScope(scope, data.dialogInputs.Next())
				if err != nil || aborted {
					return aborted, tokenScope, None[forgedomain.ForgeType](), err
				}
			}
		case forgedomain.ForgeTypeGitLab:
			data.userInput.config.NormalConfig.GitLabToken, aborted, err = dialog.GitLabToken(config.NormalConfig.GitLabToken, data.dialogInputs.Next())
			if err != nil || aborted {
				return aborted, tokenScope, None[forgedomain.ForgeType](), err
			}
			if showScopeDialog(data.userInput.config.NormalConfig.GitLabToken, config.NormalConfig.GitLabToken) {
				scope := determineScope(config.NormalConfig.GitConfig.GitLabToken)
				tokenScope, aborted, err = dialog.TokenScope(scope, data.dialogInputs.Next())
				if err != nil || aborted {
					return aborted, tokenScope, None[forgedomain.ForgeType](), err
				}
			}
		}
	}
	if configFile.HostingOriginHostname.IsNone() {
		data.userInput.config.NormalConfig.HostingOriginHostname, aborted, err = dialog.OriginHostname(config.NormalConfig.HostingOriginHostname, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncFeatureStrategy, aborted, err = dialog.SyncFeatureStrategy(config.NormalConfig.SyncFeatureStrategy, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncPerennialStrategy, aborted, err = dialog.SyncPerennialStrategy(config.NormalConfig.SyncPerennialStrategy, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncPrototypeStrategy, aborted, err = dialog.SyncPrototypeStrategy(config.NormalConfig.SyncPrototypeStrategy, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.SyncUpstream.IsNone() {
		data.userInput.config.NormalConfig.SyncUpstream, aborted, err = dialog.SyncUpstream(config.NormalConfig.SyncUpstream, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.SyncTags.IsNone() {
		data.userInput.config.NormalConfig.SyncTags, aborted, err = dialog.SyncTags(config.NormalConfig.SyncTags, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.ShareNewBranches.IsNone() {
		data.userInput.config.NormalConfig.ShareNewBranches, aborted, err = dialog.ShareNewBranches(config.NormalConfig.ShareNewBranches, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.PushHook.IsNone() {
		data.userInput.config.NormalConfig.PushHook, aborted, err = dialog.PushHook(config.NormalConfig.PushHook, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.NewBranchType.IsNone() {
		data.userInput.config.NormalConfig.NewBranchType, aborted, err = dialog.NewBranchType(config.NormalConfig.NewBranchType, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.ShipStrategy.IsNone() {
		data.userInput.config.NormalConfig.ShipStrategy, aborted, err = dialog.ShipStrategy(config.NormalConfig.ShipStrategy, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		data.userInput.config.NormalConfig.ShipDeleteTrackingBranch, aborted, err = dialog.ShipDeleteTrackingBranch(config.NormalConfig.ShipDeleteTrackingBranch, data.dialogInputs.Next())
		if err != nil || aborted {
			return aborted, tokenScope, None[forgedomain.ForgeType](), err
		}
	}
	data.userInput.configStorage, aborted, err = dialog.ConfigStorage(data.dialogInputs.Next())
	if err != nil || aborted {
		return aborted, tokenScope, None[forgedomain.ForgeType](), err
	}
	return false, tokenScope, forgeTypeOpt, nil
}

type option interface {
	IsSome() bool
}

func determineScope(global option) configdomain.ConfigScope {
	if global.IsSome() {
		return configdomain.ConfigScopeGlobal
	}
	return configdomain.ConfigScopeLocal
}

func showScopeDialog[T fmt.Stringer](input, existing T) bool {
	return input.String() != "" && input.String() != existing.String()
}

func loadSetupData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data setupData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Detached:              false,
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
	if len(remotes) == 0 {
		remotes = gitdomain.Remotes{repo.Git.DefaultRemote(repo.Backend)}
	}
	return setupData{
		config:        repo.UnvalidatedConfig,
		configFile:    repo.UnvalidatedConfig.NormalConfig.ConfigFile,
		dialogInputs:  dialogTestInputs,
		localBranches: branchesSnapshot.Branches,
		remotes:       remotes,
		userInput:     defaultUserInput(repo.UnvalidatedConfig.NormalConfig.GitConfigAccess, repo.UnvalidatedConfig.NormalConfig.GitVersion),
	}, exit, nil
}

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, configFile Option[configdomain.PartialConfig], tokenScope configdomain.ConfigScope, forgeTypeOpt Option[forgedomain.ForgeType], gitCommands git.Commands, frontend gitdomain.Runner) error {
	err := saveAliases(oldConfig.NormalConfig.Aliases, userInput.config.NormalConfig.Aliases, gitCommands, frontend)
	if err != nil {
		return err
	}
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			err = saveBitbucketUsername(oldConfig.NormalConfig.BitbucketUsername, userInput.config.NormalConfig.BitbucketUsername, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
			err = saveBitbucketAppPassword(oldConfig.NormalConfig.BitbucketAppPassword, userInput.config.NormalConfig.BitbucketAppPassword, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeCodeberg:
			err = saveCodebergToken(oldConfig.NormalConfig.CodebergToken, userInput.config.NormalConfig.CodebergToken, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeGitHub:
			err = saveGitHubToken(oldConfig.NormalConfig.GitHubToken, userInput.config.NormalConfig.GitHubToken, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeGitLab:
			err = saveGitLabToken(oldConfig.NormalConfig.GitLabToken, userInput.config.NormalConfig.GitLabToken, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeGitea:
			err = saveGiteaToken(oldConfig.NormalConfig.GiteaToken, userInput.config.NormalConfig.GiteaToken, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		}
	}
	switch userInput.configStorage {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, oldConfig)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, oldConfig, configFile, gitCommands, frontend)
	}
	panic("unknown configStorage: " + userInput.configStorage)
}

func saveToGit(userInput userInput, oldConfig config.UnvalidatedConfig, configFileOpt Option[configdomain.PartialConfig], gitCommands git.Commands, frontend gitdomain.Runner) error {
	configFile := configFileOpt.GetOrDefault()
	fc := execute.FailureCollector{}
	if configFile.NewBranchType.IsNone() {
		fc.Check(saveNewBranchType(oldConfig.NormalConfig.NewBranchType, userInput.config.NormalConfig.NewBranchType, oldConfig))
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(saveForgeType(oldConfig.NormalConfig.ForgeType, userInput.config.NormalConfig.ForgeType, gitCommands, frontend))
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(saveOriginHostname(oldConfig.NormalConfig.HostingOriginHostname, userInput.config.NormalConfig.HostingOriginHostname, gitCommands, frontend))
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(saveMainBranch(oldConfig.UnvalidatedConfig.MainBranch, userInput.config.UnvalidatedConfig.MainBranch.GetOrPanic(), oldConfig))
	}
	if len(configFile.PerennialBranches) == 0 {
		fc.Check(savePerennialBranches(oldConfig.NormalConfig.PerennialBranches, userInput.config.NormalConfig.PerennialBranches, oldConfig))
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(savePerennialRegex(oldConfig.NormalConfig.PerennialRegex, userInput.config.NormalConfig.PerennialRegex, oldConfig))
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(saveUnknownBranchType(oldConfig.NormalConfig.UnknownBranchType, userInput.config.NormalConfig.UnknownBranchType, oldConfig))
	}
	if configFile.DevRemote.IsNone() {
		fc.Check(saveDevRemote(oldConfig.NormalConfig.DevRemote, userInput.config.NormalConfig.DevRemote, oldConfig))
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(saveFeatureRegex(oldConfig.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, oldConfig))
	}
	if configFile.PushHook.IsNone() {
		fc.Check(savePushHook(oldConfig.NormalConfig.PushHook, userInput.config.NormalConfig.PushHook, oldConfig))
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(saveShareNewBranches(oldConfig.NormalConfig.ShareNewBranches, userInput.config.NormalConfig.ShareNewBranches, oldConfig))
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(saveShipStrategy(oldConfig.NormalConfig.ShipStrategy, userInput.config.NormalConfig.ShipStrategy, oldConfig))
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(saveShipDeleteTrackingBranch(oldConfig.NormalConfig.ShipDeleteTrackingBranch, userInput.config.NormalConfig.ShipDeleteTrackingBranch, oldConfig))
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(saveSyncFeatureStrategy(oldConfig.NormalConfig.SyncFeatureStrategy, userInput.config.NormalConfig.SyncFeatureStrategy, oldConfig))
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(saveSyncPerennialStrategy(oldConfig.NormalConfig.SyncPerennialStrategy, userInput.config.NormalConfig.SyncPerennialStrategy, oldConfig))
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(saveSyncPrototypeStrategy(oldConfig.NormalConfig.SyncPrototypeStrategy, userInput.config.NormalConfig.SyncPrototypeStrategy, oldConfig))
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(saveSyncUpstream(oldConfig.NormalConfig.SyncUpstream, userInput.config.NormalConfig.SyncUpstream, oldConfig))
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(saveSyncTags(oldConfig.NormalConfig.SyncTags, userInput.config.NormalConfig.SyncTags, oldConfig))
	}
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

func saveBitbucketAppPassword(oldPassword, newPassword Option[configdomain.BitbucketAppPassword], scope configdomain.ConfigScope, gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newPassword.Equal(oldPassword) {
		return nil
	}
	if value, has := newPassword.Get(); has {
		return gitCommands.SetBitbucketAppPassword(frontend, value, scope)
	}
	return gitCommands.RemoveBitbucketAppPassword(frontend)
}

func saveBitbucketUsername(oldValue, newValue Option[configdomain.BitbucketUsername], scope configdomain.ConfigScope, gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newValue.Equal(oldValue) {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitCommands.SetBitbucketUsername(frontend, value, scope)
	}
	return gitCommands.RemoveBitbucketUsername(frontend)
}

func saveNewBranchType(oldValue, newValue Option[configdomain.BranchType], config config.UnvalidatedConfig) error {
	if newValue.Equal(oldValue) {
		return nil
	}
	if value, hasValue := newValue.Get(); hasValue {
		return config.NormalConfig.SetNewBranchType(value)
	}
	config.NormalConfig.RemoveNewBranchType()
	return nil
}

func saveUnknownBranchType(oldValue, newValue configdomain.BranchType, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetUnknownBranchTypeLocally(newValue)
}

func saveDevRemote(oldValue, newValue gitdomain.Remote, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetDevRemote(newValue)
}

func saveFeatureRegex(oldValue, newValue Option[configdomain.FeatureRegex], config config.UnvalidatedConfig) error {
	if newValue.Equal(oldValue) {
		return nil
	}
	if value, has := newValue.Get(); has {
		return config.NormalConfig.SetFeatureRegexLocally(value)
	}
	config.NormalConfig.RemoveFeatureRegex()
	return nil
}

func saveForgeType(oldForgeType, newForgeType Option[forgedomain.ForgeType], gitCommands git.Commands, frontend gitdomain.Runner) (err error) {
	oldValue, oldHas := oldForgeType.Get()
	newValue, newHas := newForgeType.Get()
	if !oldHas && !newHas {
		return nil
	}
	if oldValue == newValue {
		return nil
	}
	if newHas {
		return gitCommands.SetForgeType(frontend, newValue)
	}
	return gitCommands.DeleteConfigEntryForgeType(frontend)
}

func saveCodebergToken(oldToken, newToken Option[configdomain.CodebergToken], scope configdomain.ConfigScope, gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetCodebergToken(frontend, value, scope)
	}
	return gitCommands.RemoveCodebergToken(frontend)
}

func saveGiteaToken(oldToken, newToken Option[configdomain.GiteaToken], scope configdomain.ConfigScope, gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGiteaToken(frontend, value, scope)
	}
	return gitCommands.RemoveGiteaToken(frontend)
}

func saveGitHubToken(oldToken, newToken Option[configdomain.GitHubToken], scope configdomain.ConfigScope, gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGitHubToken(frontend, value, scope)
	}
	return gitCommands.RemoveGitHubToken(frontend)
}

func saveGitLabToken(oldToken, newToken Option[configdomain.GitLabToken], scope configdomain.ConfigScope, gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGitLabToken(frontend, value, scope)
	}
	return gitCommands.RemoveGitLabToken(frontend)
}

func saveMainBranch(oldValue Option[gitdomain.LocalBranchName], newValue gitdomain.LocalBranchName, config config.UnvalidatedConfig) error {
	if Some(newValue).Equal(oldValue) {
		return nil
	}
	return config.SetMainBranch(newValue)
}

func saveOriginHostname(oldValue, newValue Option[configdomain.HostingOriginHostname], gitCommands git.Commands, frontend gitdomain.Runner) error {
	if newValue.Equal(oldValue) {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitCommands.SetOriginHostname(frontend, value)
	}
	return gitCommands.DeleteConfigEntryOriginHostname(frontend)
}

func savePerennialBranches(oldValue, newValue gitdomain.LocalBranchNames, config config.UnvalidatedConfig) error {
	if slices.Compare(oldValue, newValue) != 0 || config.NormalConfig.GitConfig.PerennialBranches == nil {
		return config.NormalConfig.SetPerennialBranches(newValue)
	}
	return nil
}

func savePerennialRegex(oldValue, newValue Option[configdomain.PerennialRegex], config config.UnvalidatedConfig) error {
	if newValue.Equal(oldValue) {
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

func saveShareNewBranches(oldValue, newValue configdomain.ShareNewBranches, config config.UnvalidatedConfig) error {
	if newValue == oldValue {
		return nil
	}
	return config.NormalConfig.SetShareNewBranches(newValue, configdomain.ConfigScopeLocal)
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
	config.NormalConfig.RemoveShareNewBranches()
	config.NormalConfig.RemovePushHook()
	config.NormalConfig.RemoveShipStrategy()
	config.NormalConfig.RemoveShipDeleteTrackingBranch()
	config.NormalConfig.RemoveSyncFeatureStrategy()
	config.NormalConfig.RemoveSyncPerennialStrategy()
	config.NormalConfig.RemoveSyncPrototypeStrategy()
	config.NormalConfig.RemoveSyncUpstream()
	config.NormalConfig.RemoveSyncTags()
	err = saveUnknownBranchType(config.NormalConfig.UnknownBranchType, userInput.config.NormalConfig.UnknownBranchType, config)
	if err != nil {
		return err
	}
	return saveFeatureRegex(config.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, config)
}
