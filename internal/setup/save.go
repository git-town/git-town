package setup

import (
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func Save(userInput UserInput, unvalidatedConfig config.UnvalidatedConfig, data Data, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(userInput.data.Aliases, unvalidatedConfig.GitGlobal.Aliases, frontend),
	)
	if forgeType, hasForgeType := userInput.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(userInput.data.BitbucketUsername, unvalidatedConfig.GitLocal.BitbucketUsername, userInput.scope, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(userInput.data.BitbucketAppPassword, unvalidatedConfig.GitLocal.BitbucketAppPassword, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeCodeberg:
			fc.Check(
				saveCodebergToken(userInput.data.CodebergToken, unvalidatedConfig.GitLocal.CodebergToken, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeGitHub:
			fc.Check(
				saveGitHubToken(userInput.data.GitHubToken, unvalidatedConfig.GitLocal.GitHubToken, userInput.scope, userInput.data.GitHubConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitLab:
			fc.Check(
				saveGitLabToken(userInput.data.GitLabToken, unvalidatedConfig.GitLocal.GitLabToken, userInput.scope, userInput.data.GitLabConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(userInput.data.GiteaToken, unvalidatedConfig.GitLocal.GiteaToken, userInput.scope, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.storageLocation {
	case dialog.ConfigStorageOptionFile:
		return saveAllToFile(userInput, unvalidatedConfig.GitLocal, frontend)
	case dialog.ConfigStorageOptionGit:
		return saveAllToGit(userInput, unvalidatedConfig.GitLocal, unvalidatedConfig.File, data, frontend)
	}
	return nil
}

func saveAliases(valuesToWriteToGit configdomain.Aliases, valuesAlreadyInGit configdomain.Aliases, frontend subshelldomain.Runner) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := valuesAlreadyInGit[aliasableCommand]
		newAlias, hasNew := valuesToWriteToGit[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err = gitconfig.RemoveAlias(frontend, aliasableCommand)
		case hasNew && !hasOld, newAlias != oldAlias:
			err = gitconfig.SetAlias(frontend, aliasableCommand)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveAllToFile(userInput UserInput, gitConfig configdomain.PartialConfig, runner subshelldomain.Runner) error {
	userInput.data.MainBranch = Some(userInput.validatedConfig.MainBranch)
	if err := configfile.Save(userInput.data); err != nil {
		return err
	}
	if gitConfig.ContributionRegex.IsSome() {
		_ = gitconfig.RemoveContributionRegex(runner)
	}
	if gitConfig.DevRemote.IsSome() {
		_ = gitconfig.RemoveDevRemote(runner)
	}
	if gitConfig.FeatureRegex.IsSome() {
		_ = gitconfig.RemoveFeatureRegex(runner)
	}
	if gitConfig.MainBranch.IsSome() {
		_ = gitconfig.RemoveMainBranch(runner)
	}
	if gitConfig.NewBranchType.IsSome() {
		_ = gitconfig.RemoveNewBranchType(runner)
	}
	if gitConfig.ObservedRegex.IsSome() {
		_ = gitconfig.RemoveObservedRegex(runner)
	}
	if len(gitConfig.PerennialBranches) > 0 {
		_ = gitconfig.RemovePerennialBranches(runner)
	}
	if gitConfig.PerennialRegex.IsSome() {
		_ = gitconfig.RemovePerennialRegex(runner)
	}
	if gitConfig.ShareNewBranches.IsSome() {
		_ = gitconfig.RemoveShareNewBranches(runner)
	}
	if gitConfig.PushHook.IsSome() {
		_ = gitconfig.RemovePushHook(runner)
	}
	if gitConfig.ShipStrategy.IsSome() {
		_ = gitconfig.RemoveShipStrategy(runner)
	}
	if gitConfig.ShipDeleteTrackingBranch.IsSome() {
		_ = gitconfig.RemoveShipDeleteTrackingBranch(runner)
	}
	if gitConfig.SyncFeatureStrategy.IsSome() {
		_ = gitconfig.RemoveSyncFeatureStrategy(runner)
	}
	if gitConfig.SyncPerennialStrategy.IsSome() {
		_ = gitconfig.RemoveSyncPerennialStrategy(runner)
	}
	if gitConfig.SyncPrototypeStrategy.IsSome() {
		_ = gitconfig.RemoveSyncPrototypeStrategy(runner)
	}
	if gitConfig.SyncUpstream.IsSome() {
		_ = gitconfig.RemoveSyncUpstream(runner)
	}
	if gitConfig.SyncTags.IsSome() {
		_ = gitconfig.RemoveSyncTags(runner)
	}
	if gitConfig.UnknownBranchType.IsSome() {
		_ = gitconfig.RemoveUnknownBranchType(runner)
	}
	if err := saveUnknownBranchType(userInput.data.UnknownBranchType, gitConfig.UnknownBranchType, runner); err != nil {
		return err
	}
	return saveFeatureRegex(userInput.data.FeatureRegex, gitConfig.FeatureRegex, runner)
}

func saveAllToGit(userInput UserInput, existingGitConfig configdomain.PartialConfig, configFile configdomain.PartialConfig, data Data, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	if configFile.NewBranchType.IsNone() {
		fc.Check(
			saveNewBranchType(userInput.data.NewBranchType, existingGitConfig.NewBranchType, frontend),
		)
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(
			saveForgeType(userInput.data.ForgeType, existingGitConfig.ForgeType, frontend),
		)
	}
	if configFile.GitHubConnectorType.IsNone() {
		fc.Check(
			saveGitHubConnectorType(userInput.data.GitHubConnectorType, existingGitConfig.GitHubConnectorType, frontend),
		)
	}
	if configFile.GitLabConnectorType.IsNone() {
		fc.Check(
			saveGitLabConnectorType(userInput.data.GitLabConnectorType, existingGitConfig.GitLabConnectorType, frontend),
		)
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(
			saveOriginHostname(userInput.data.HostingOriginHostname, existingGitConfig.HostingOriginHostname, frontend),
		)
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(
			saveMainBranch(userInput.data.MainBranch, existingGitConfig.MainBranch, frontend),
		)
	}
	fc.Check(
		savePerennialBranches(userInput.data.PerennialBranches, existingGitConfig.PerennialBranches, frontend),
	)
	if configFile.PerennialRegex.IsNone() {
		fc.Check(
			savePerennialRegex(userInput.data.PerennialRegex, existingGitConfig.PerennialRegex, frontend),
		)
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(
			saveUnknownBranchType(userInput.data.UnknownBranchType, existingGitConfig.UnknownBranchType, frontend),
		)
	}
	if len(data.Remotes) > 1 && configFile.DevRemote.IsNone() {
		fc.Check(
			saveDevRemote(userInput.data.DevRemote, existingGitConfig.DevRemote, frontend),
		)
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(
			saveFeatureRegex(userInput.data.FeatureRegex, existingGitConfig.FeatureRegex, frontend),
		)
	}
	if configFile.ContributionRegex.IsNone() {
		fc.Check(
			saveContributionRegex(userInput.data.ContributionRegex, existingGitConfig.ContributionRegex, frontend),
		)
	}
	if configFile.ObservedRegex.IsNone() {
		fc.Check(
			saveObservedRegex(userInput.data.ObservedRegex, existingGitConfig.ObservedRegex, frontend),
		)
	}
	if configFile.PushHook.IsNone() {
		fc.Check(
			savePushHook(userInput.data.PushHook, existingGitConfig.PushHook, frontend),
		)
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(
			saveShareNewBranches(userInput.data.ShareNewBranches, existingGitConfig.ShareNewBranches, frontend),
		)
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(
			saveShipStrategy(userInput.data.ShipStrategy, existingGitConfig.ShipStrategy, frontend),
		)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(
			saveShipDeleteTrackingBranch(userInput.data.ShipDeleteTrackingBranch, existingGitConfig.ShipDeleteTrackingBranch, frontend),
		)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(
			saveSyncFeatureStrategy(userInput.data.SyncFeatureStrategy, existingGitConfig.SyncFeatureStrategy, frontend),
		)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(
			saveSyncPerennialStrategy(userInput.data.SyncPerennialStrategy, existingGitConfig.SyncPerennialStrategy, frontend),
		)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(
			saveSyncPrototypeStrategy(userInput.data.SyncPrototypeStrategy, existingGitConfig.SyncPrototypeStrategy, frontend),
		)
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(
			saveSyncUpstream(userInput.data.SyncUpstream, existingGitConfig.SyncUpstream, frontend),
		)
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(
			saveSyncTags(userInput.data.SyncTags, existingGitConfig.SyncTags, frontend),
		)
	}
	return fc.Err
}

func saveBitbucketAppPassword(valueToWriteToGit Option[forgedomain.BitbucketAppPassword], valueAlreadyInGit Option[forgedomain.BitbucketAppPassword], scope configdomain.ConfigScope, runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetBitbucketAppPassword(runner, value, scope)
	}
	return gitconfig.RemoveBitbucketAppPassword(runner)
}

func saveBitbucketUsername(valueToWriteToGit Option[forgedomain.BitbucketUsername], valueAlreadyInGit Option[forgedomain.BitbucketUsername], scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetBitbucketUsername(frontend, value, scope)
	}
	return gitconfig.RemoveBitbucketUsername(frontend)
}

func saveCodebergToken(valueToWriteToGit Option[forgedomain.CodebergToken], valueAlreadyInGit Option[forgedomain.CodebergToken], scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetCodebergToken(frontend, value, scope)
	}
	return gitconfig.RemoveCodebergToken(frontend)
}

func saveContributionRegex(valueToWriteToGit Option[configdomain.ContributionRegex], valueAlreadyInGit Option[configdomain.ContributionRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetContributionRegex(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveContributionRegex(runner)
	return nil
}

func saveDevRemote(valueToWriteToGit Option[gitdomain.Remote], valueAlreadyInGit Option[gitdomain.Remote], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetDevRemote(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveDevRemote(runner)
}

func saveFeatureRegex(valueToWriteToGit Option[configdomain.FeatureRegex], valueAlreadyInGit Option[configdomain.FeatureRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetFeatureRegex(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveFeatureRegex(runner)
	return nil
}

func saveForgeType(valueToWriteToGit Option[forgedomain.ForgeType], valueAlreadyInGit Option[forgedomain.ForgeType], frontend subshelldomain.Runner) (err error) {
	oldValue, oldHas := valueAlreadyInGit.Get()
	newValue, newHas := valueToWriteToGit.Get()
	if !oldHas && !newHas {
		return nil
	}
	if oldValue == newValue {
		return nil
	}
	if newHas {
		return gitconfig.SetForgeType(frontend, newValue, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveForgeType(frontend)
}

func saveGitHubConnectorType(valueToWriteToGit Option[forgedomain.GitHubConnectorType], valueAlreadyInGit Option[forgedomain.GitHubConnectorType], frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGitHubConnectorType(frontend, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveGitHubConnectorType(frontend)
}

func saveGitHubToken(valueToWriteToGit Option[forgedomain.GitHubToken], valueAlreadyInGit Option[forgedomain.GitHubToken], scope configdomain.ConfigScope, githubConnectorType Option[forgedomain.GitHubConnectorType], frontend subshelldomain.Runner) error {
	if connectorType, has := githubConnectorType.Get(); has {
		if connectorType == forgedomain.GitHubConnectorTypeGh {
			return nil
		}
	}
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGitHubToken(frontend, value, scope)
	}
	return gitconfig.RemoveGitHubToken(frontend)
}

func saveGitLabConnectorType(valueToWriteToGit Option[forgedomain.GitLabConnectorType], valueAlreadyInGit Option[forgedomain.GitLabConnectorType], frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGitLabConnectorType(frontend, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveGitLabConnectorType(frontend)
}

func saveGitLabToken(valueToWriteToGit Option[forgedomain.GitLabToken], valueAlreadyInGit Option[forgedomain.GitLabToken], scope configdomain.ConfigScope, gitlabConnectorType Option[forgedomain.GitLabConnectorType], frontend subshelldomain.Runner) error {
	if connectorType, has := gitlabConnectorType.Get(); has {
		if connectorType == forgedomain.GitLabConnectorTypeGlab {
			return nil
		}
	}
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGitLabToken(frontend, value, scope)
	}
	return gitconfig.RemoveGitLabToken(frontend)
}

func saveGiteaToken(valueToWriteToGit Option[forgedomain.GiteaToken], valueAlreadyInGit Option[forgedomain.GiteaToken], scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGiteaToken(frontend, value, scope)
	}
	return gitconfig.RemoveGiteaToken(frontend)
}

func saveMainBranch(valueToWriteToGit Option[gitdomain.LocalBranchName], valueAlreadyInGit Option[gitdomain.LocalBranchName], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetMainBranch(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveMainBranch(runner)
}

func saveNewBranchType(valueToWriteToGit Option[configdomain.NewBranchType], valueAlreadyInGit Option[configdomain.NewBranchType], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetNewBranchType(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveNewBranchType(runner)
	return nil
}

func saveObservedRegex(valueToWriteToGit Option[configdomain.ObservedRegex], valueAlreadyInGit Option[configdomain.ObservedRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetObservedRegex(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveObservedRegex(runner)
	return nil
}

func saveOriginHostname(valueToWriteToGit Option[configdomain.HostingOriginHostname], valueAlreadyInGit Option[configdomain.HostingOriginHostname], frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetOriginHostname(frontend, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveOriginHostname(frontend)
}

func savePerennialBranches(valueToWriteToGit gitdomain.LocalBranchNames, valueAlreadyInGit gitdomain.LocalBranchNames, runner subshelldomain.Runner) error {
	if slices.Compare(valueAlreadyInGit, valueToWriteToGit) == 0 {
		return nil
	}
	return gitconfig.SetPerennialBranches(runner, valueToWriteToGit, configdomain.ConfigScopeLocal)
}

func savePerennialRegex(valueToWriteToGit Option[configdomain.PerennialRegex], valueAlreadyInGit Option[configdomain.PerennialRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetPerennialRegex(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemovePerennialRegex(runner)
	return nil
}

func savePushHook(valueToWriteToGit Option[configdomain.PushHook], valueAlreadyInGit Option[configdomain.PushHook], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetPushHook(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemovePushHook(runner)
}

func saveShareNewBranches(valueToWriteToGit Option[configdomain.ShareNewBranches], valueAlreadyInGit Option[configdomain.ShareNewBranches], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetShareNewBranches(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveShareNewBranches(runner)
}

func saveShipDeleteTrackingBranch(valueToWriteToGit Option[configdomain.ShipDeleteTrackingBranch], valueAlreadyInGit Option[configdomain.ShipDeleteTrackingBranch], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetShipDeleteTrackingBranch(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveShipDeleteTrackingBranch(runner)
}

func saveShipStrategy(valueToWriteToGit Option[configdomain.ShipStrategy], valueAlreadyInGit Option[configdomain.ShipStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetShipStrategy(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveShipStrategy(runner)
}

func saveSyncFeatureStrategy(valueToWriteToGit Option[configdomain.SyncFeatureStrategy], valueAlreadyInGit Option[configdomain.SyncFeatureStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetSyncFeatureStrategy(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveSyncFeatureStrategy(runner)
}

func saveSyncPerennialStrategy(valueToWriteToGit Option[configdomain.SyncPerennialStrategy], valueAlreadyInGit Option[configdomain.SyncPerennialStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetSyncPerennialStrategy(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveSyncPerennialStrategy(runner)
}

func saveSyncPrototypeStrategy(valueToWriteToGit Option[configdomain.SyncPrototypeStrategy], valueAlreadyInGit Option[configdomain.SyncPrototypeStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetSyncPrototypeStrategy(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveSyncPrototypeStrategy(runner)
}

func saveSyncTags(valueToWriteToGit Option[configdomain.SyncTags], valueAlreadyInGit Option[configdomain.SyncTags], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetSyncTags(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveSyncTags(runner)
}

func saveSyncUpstream(valueToWriteToGit Option[configdomain.SyncUpstream], valueAlreadyInGit Option[configdomain.SyncUpstream], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetSyncUpstream(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveSyncUpstream(runner)
}

func saveUnknownBranchType(valueToWriteToGit Option[configdomain.UnknownBranchType], valueAlreadyInGit Option[configdomain.UnknownBranchType], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetUnknownBranchType(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveUnknownBranchType(runner)
}
