package configsetup

import (
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, configFile Option[configdomain.PartialConfig], tokenScope configdomain.ConfigScope, forgeTypeOpt Option[forgedomain.ForgeType], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(oldConfig.NormalConfig.Aliases, userInput.config.NormalConfig.Aliases, gitCommands, frontend),
	)
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(oldConfig.NormalConfig.BitbucketUsername, userInput.config.NormalConfig.BitbucketUsername, tokenScope, gitCommands, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(oldConfig.NormalConfig.BitbucketAppPassword, userInput.config.NormalConfig.BitbucketAppPassword, tokenScope, gitCommands, frontend),
			)
		case forgedomain.ForgeTypeCodeberg:
			fc.Check(
				saveCodebergToken(oldConfig.NormalConfig.CodebergToken, userInput.config.NormalConfig.CodebergToken, tokenScope, gitCommands, frontend),
			)
		case forgedomain.ForgeTypeGitHub:
			fc.Check(
				saveGitHubToken(oldConfig.NormalConfig.GitHubToken, userInput.config.NormalConfig.GitHubToken, tokenScope, userInput.config.NormalConfig.GitHubConnectorType, gitCommands, frontend),
			)
		case forgedomain.ForgeTypeGitLab:
			fc.Check(
				saveGitLabToken(oldConfig.NormalConfig.GitLabToken, userInput.config.NormalConfig.GitLabToken, tokenScope, userInput.config.NormalConfig.GitLabConnectorType, gitCommands, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(oldConfig.NormalConfig.GiteaToken, userInput.config.NormalConfig.GiteaToken, tokenScope, gitCommands, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.configStorage {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, oldConfig)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, oldConfig, configFile, gitCommands, frontend)
	}
	return nil
}

func saveToGit(userInput userInput, oldConfig config.UnvalidatedConfig, configFileOpt Option[configdomain.PartialConfig], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	configFile := configFileOpt.GetOrDefault()
	fc := gohacks.ErrorCollector{}
	if configFile.NewBranchType.IsNone() {
		fc.Check(
			saveNewBranchType(oldConfig.NormalConfig.NewBranchType, userInput.config.NormalConfig.NewBranchType, oldConfig),
		)
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(
			saveForgeType(oldConfig.NormalConfig.ForgeType, userInput.config.NormalConfig.ForgeType, gitCommands, frontend),
		)
	}
	if configFile.GitHubConnectorType.IsNone() {
		fc.Check(
			saveGitHubConnectorType(oldConfig.NormalConfig.GitHubConnectorType, userInput.config.NormalConfig.GitHubConnectorType, gitCommands, frontend),
		)
	}
	if configFile.GitLabConnectorType.IsNone() {
		fc.Check(
			saveGitLabConnectorType(oldConfig.NormalConfig.GitLabConnectorType, userInput.config.NormalConfig.GitLabConnectorType, gitCommands, frontend),
		)
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(
			saveOriginHostname(oldConfig.NormalConfig.HostingOriginHostname, userInput.config.NormalConfig.HostingOriginHostname, gitCommands, frontend),
		)
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(
			saveMainBranch(oldConfig.UnvalidatedConfig.MainBranch, userInput.config.UnvalidatedConfig.MainBranch, oldConfig),
		)
	}
	if len(configFile.PerennialBranches) == 0 {
		fc.Check(
			savePerennialBranches(oldConfig.NormalConfig.PerennialBranches, userInput.config.NormalConfig.PerennialBranches, oldConfig),
		)
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(
			savePerennialRegex(oldConfig.NormalConfig.PerennialRegex, userInput.config.NormalConfig.PerennialRegex, oldConfig),
		)
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(
			saveUnknownBranchType(oldConfig.NormalConfig.UnknownBranchType, userInput.config.NormalConfig.UnknownBranchType, oldConfig),
		)
	}
	if configFile.DevRemote.IsNone() {
		fc.Check(
			saveDevRemote(oldConfig.NormalConfig.DevRemote, userInput.config.NormalConfig.DevRemote, oldConfig),
		)
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(
			saveFeatureRegex(oldConfig.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, oldConfig),
		)
	}
	if configFile.PushHook.IsNone() {
		fc.Check(
			savePushHook(oldConfig.NormalConfig.PushHook, userInput.config.NormalConfig.PushHook, oldConfig),
		)
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(
			saveShareNewBranches(oldConfig.NormalConfig.ShareNewBranches, userInput.config.NormalConfig.ShareNewBranches, oldConfig),
		)
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(
			saveShipStrategy(oldConfig.NormalConfig.ShipStrategy, userInput.config.NormalConfig.ShipStrategy, oldConfig),
		)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(
			saveShipDeleteTrackingBranch(oldConfig.NormalConfig.ShipDeleteTrackingBranch, userInput.config.NormalConfig.ShipDeleteTrackingBranch, oldConfig),
		)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(
			saveSyncFeatureStrategy(oldConfig.NormalConfig.SyncFeatureStrategy, userInput.config.NormalConfig.SyncFeatureStrategy, oldConfig),
		)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(
			saveSyncPerennialStrategy(oldConfig.NormalConfig.SyncPerennialStrategy, userInput.config.NormalConfig.SyncPerennialStrategy, oldConfig),
		)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(
			saveSyncPrototypeStrategy(oldConfig.NormalConfig.SyncPrototypeStrategy, userInput.config.NormalConfig.SyncPrototypeStrategy, oldConfig),
		)
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(
			saveSyncUpstream(oldConfig.NormalConfig.SyncUpstream, userInput.config.NormalConfig.SyncUpstream, oldConfig),
		)
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(
			saveSyncTags(oldConfig.NormalConfig.SyncTags, userInput.config.NormalConfig.SyncTags, oldConfig),
		)
	}
	return fc.Err
}

func saveAliases(oldAliases, newAliases configdomain.Aliases, gitCommands git.Commands, frontend subshelldomain.Runner) (err error) {
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

func saveBitbucketAppPassword(oldPassword, newPassword Option[forgedomain.BitbucketAppPassword], scope configdomain.ConfigScope, gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if newPassword.Equal(oldPassword) {
		return nil
	}
	if value, has := newPassword.Get(); has {
		return gitCommands.SetBitbucketAppPassword(frontend, value, scope)
	}
	return gitCommands.RemoveBitbucketAppPassword(frontend)
}

func saveBitbucketUsername(oldValue, newValue Option[forgedomain.BitbucketUsername], scope configdomain.ConfigScope, gitCommands git.Commands, frontend subshelldomain.Runner) error {
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

func saveForgeType(oldForgeType, newForgeType Option[forgedomain.ForgeType], gitCommands git.Commands, frontend subshelldomain.Runner) (err error) {
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

func saveCodebergToken(oldToken, newToken Option[forgedomain.CodebergToken], scope configdomain.ConfigScope, gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetCodebergToken(frontend, value, scope)
	}
	return gitCommands.RemoveCodebergToken(frontend)
}

func saveGiteaToken(oldToken, newToken Option[forgedomain.GiteaToken], scope configdomain.ConfigScope, gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGiteaToken(frontend, value, scope)
	}
	return gitCommands.RemoveGiteaToken(frontend)
}

func saveGitHubConnectorType(oldType, newType Option[forgedomain.GitHubConnectorType], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if newType.Equal(oldType) {
		return nil
	}
	if value, has := newType.Get(); has {
		return gitCommands.SetGitHubConnectorType(frontend, value)
	}
	return gitCommands.RemoveGitHubConnectorType(frontend)
}

func saveGitHubToken(oldToken, newToken Option[forgedomain.GitHubToken], scope configdomain.ConfigScope, githubConnectorType Option[forgedomain.GitHubConnectorType], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if connectorType, has := githubConnectorType.Get(); has {
		if connectorType == forgedomain.GitHubConnectorTypeGh {
			return nil
		}
	}
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGitHubToken(frontend, value, scope)
	}
	return gitCommands.RemoveGitHubToken(frontend)
}

func saveGitLabConnectorType(oldType, newType Option[forgedomain.GitLabConnectorType], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if newType.Equal(oldType) {
		return nil
	}
	if value, has := newType.Get(); has {
		return gitCommands.SetGitLabConnectorType(frontend, value)
	}
	return gitCommands.RemoveGitLabConnectorType(frontend)
}

func saveGitLabToken(oldToken, newToken Option[forgedomain.GitLabToken], scope configdomain.ConfigScope, gitlabConnectorType Option[forgedomain.GitLabConnectorType], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	if connectorType, has := gitlabConnectorType.Get(); has {
		if connectorType == forgedomain.GitLabConnectorTypeGlab {
			return nil
		}
	}
	if newToken.Equal(oldToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitCommands.SetGitLabToken(frontend, value, scope)
	}
	return gitCommands.RemoveGitLabToken(frontend)
}

func saveMainBranch(oldValue Option[gitdomain.LocalBranchName], newValue Option[gitdomain.LocalBranchName], config config.UnvalidatedConfig) error {
	if newValue.Equal(oldValue) {
		return nil
	}
	if mainBranch, hasNewValue := newValue.Get(); hasNewValue {
		return config.SetMainBranch(mainBranch)
	}
	return nil
}

func saveOriginHostname(oldValue, newValue Option[configdomain.HostingOriginHostname], gitCommands git.Commands, frontend subshelldomain.Runner) error {
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
	if err := configfile.Save(&userInput.config); err != nil {
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
	if err := saveUnknownBranchType(config.NormalConfig.UnknownBranchType, userInput.config.NormalConfig.UnknownBranchType, config); err != nil {
		return err
	}
	return saveFeatureRegex(config.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, config)
}
