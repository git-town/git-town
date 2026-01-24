package setup

import (
	"slices"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func Save(userInput UserInput, unvalidatedConfig config.UnvalidatedConfig, data Data, enterAll bool, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(userInput.Data.Aliases, unvalidatedConfig.GitGlobal.Aliases, frontend),
	)
	if forgeType, hasForgeType := userInput.DeterminedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeAzuredevops:
			// no API token for now
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(userInput.Data.BitbucketUsername, unvalidatedConfig.GitLocal.BitbucketUsername, userInput.Scope, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(userInput.Data.BitbucketAppPassword, unvalidatedConfig.GitLocal.BitbucketAppPassword, userInput.Scope, frontend),
			)
		case forgedomain.ForgeTypeForgejo:
			fc.Check(
				saveForgejoToken(userInput.Data.ForgejoToken, unvalidatedConfig.GitLocal.ForgejoToken, userInput.Scope, frontend),
			)
		case forgedomain.ForgeTypeGithub:
			fc.Check(
				saveGithubToken(userInput.Data.GithubToken, unvalidatedConfig.GitLocal.GithubToken, userInput.Scope, userInput.Data.GithubConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitlab:
			fc.Check(
				saveGitlabToken(userInput.Data.GitlabToken, unvalidatedConfig.GitLocal.GitlabToken, userInput.Scope, userInput.Data.GitlabConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(userInput.Data.GiteaToken, unvalidatedConfig.GitLocal.GiteaToken, userInput.Scope, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.StorageLocation {
	case dialog.ConfigStorageOptionFile:
		return saveAllToFile(userInput, unvalidatedConfig.File, unvalidatedConfig.GitLocal, frontend)
	case dialog.ConfigStorageOptionGit:
		return saveAllToGit(userInput, unvalidatedConfig.GitLocal, unvalidatedConfig.File, data, enterAll, frontend)
	}
	return nil
}

func saveAliases(valuesToWriteToGit configdomain.Aliases, valuesAlreadyInGit configdomain.Aliases, frontend subshelldomain.Runner) error {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := valuesAlreadyInGit[aliasableCommand]
		newAlias, hasNew := valuesToWriteToGit[aliasableCommand]
		var err error
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

func saveAllToFile(userInput UserInput, existingConfigFile configdomain.PartialConfig, gitConfig configdomain.PartialConfig, runner subshelldomain.Runner) error {
	userInput.Data.MainBranch = Some(userInput.ValidatedConfig.MainBranch)
	configData := existingConfigFile.Merge(userInput.Data)
	if err := configfile.Save(configData); err != nil {
		return err
	}
	// keep-sorted start block=yes
	if gitConfig.AutoSync.IsSome() {
		_ = gitconfig.RemoveAutoSync(runner)
	}
	if gitConfig.BranchPrefix.IsSome() {
		_ = gitconfig.RemoveBranchPrefix(runner)
	}
	if gitConfig.ContributionRegex.IsSome() {
		_ = gitconfig.RemoveContributionRegex(runner)
	}
	if gitConfig.Detached.IsSome() {
		_ = gitconfig.RemoveDetached(runner)
	}
	if gitConfig.DevRemote.IsSome() {
		_ = gitconfig.RemoveDevRemote(runner)
	}
	if gitConfig.FeatureRegex.IsSome() {
		_ = gitconfig.RemoveFeatureRegex(runner)
	}
	if gitConfig.ForgeType.IsSome() {
		_ = gitconfig.RemoveForgeType(runner)
	}
	if gitConfig.GithubConnectorType.IsSome() {
		_ = gitconfig.RemoveGithubConnectorType(runner)
	}
	if gitConfig.GitlabConnectorType.IsSome() {
		_ = gitconfig.RemoveGitlabConnectorType(runner)
	}
	if gitConfig.HostingOriginHostname.IsSome() {
		_ = gitconfig.RemoveOriginHostname(runner)
	}
	if gitConfig.IgnoreUncommitted.IsSome() {
		_ = gitconfig.RemoveIgnoreUncommitted(runner)
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
	if gitConfig.Order.IsSome() {
		_ = gitconfig.RemoveOrder(runner)
	}
	if gitConfig.PerennialRegex.IsSome() {
		_ = gitconfig.RemovePerennialRegex(runner)
	}
	if gitConfig.ProposalsShowLineage.IsSome() {
		_ = gitconfig.RemoveProposalsShowLineage(runner)
	}
	if gitConfig.PushBranches.IsSome() {
		_ = gitconfig.RemovePushBranches(runner)
	}
	if gitConfig.PushHook.IsSome() {
		_ = gitconfig.RemovePushHook(runner)
	}
	if gitConfig.ShareNewBranches.IsSome() {
		_ = gitconfig.RemoveShareNewBranches(runner)
	}
	if gitConfig.ShipDeleteTrackingBranch.IsSome() {
		_ = gitconfig.RemoveShipDeleteTrackingBranch(runner)
	}
	if gitConfig.ShipStrategy.IsSome() {
		_ = gitconfig.RemoveShipStrategy(runner)
	}
	if gitConfig.Stash.IsSome() {
		_ = gitconfig.RemoveStash(runner)
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
	if gitConfig.SyncTags.IsSome() {
		_ = gitconfig.RemoveSyncTags(runner)
	}
	if gitConfig.SyncUpstream.IsSome() {
		_ = gitconfig.RemoveSyncUpstream(runner)
	}
	if gitConfig.UnknownBranchType.IsSome() {
		_ = gitconfig.RemoveUnknownBranchType(runner)
	}
	if len(gitConfig.PerennialBranches) > 0 {
		_ = gitconfig.RemovePerennialBranches(runner)
	}
	// keep-sorted end
	if err := saveUnknownBranchType(userInput.Data.UnknownBranchType, gitConfig.UnknownBranchType, runner); err != nil {
		return err
	}
	return saveFeatureRegex(userInput.Data.FeatureRegex, gitConfig.FeatureRegex, runner)
}

func saveAllToGit(userInput UserInput, existingGitConfig configdomain.PartialConfig, configFile configdomain.PartialConfig, data Data, enterAll bool, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}

	// BASIC CONFIGURATION
	if configFile.MainBranch.IsNone() {
		fc.Check(
			saveMainBranch(userInput.Data.MainBranch, existingGitConfig.MainBranch, frontend),
		)
	}
	fc.Check(
		savePerennialBranches(userInput.Data.PerennialBranches, existingGitConfig.PerennialBranches, frontend),
	)
	if len(data.Remotes) > 1 && configFile.DevRemote.IsNone() {
		fc.Check(
			saveDevRemote(userInput.Data.DevRemote, existingGitConfig.DevRemote, frontend),
		)
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(
			saveOriginHostname(userInput.Data.HostingOriginHostname, existingGitConfig.HostingOriginHostname, frontend),
		)
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(
			saveForgeType(userInput.Data.ForgeType, existingGitConfig.ForgeType, frontend),
		)
	}
	if configFile.GithubConnectorType.IsNone() {
		fc.Check(
			saveGithubConnectorType(userInput.Data.GithubConnectorType, existingGitConfig.GithubConnectorType, frontend),
		)
	}
	if configFile.GitlabConnectorType.IsNone() {
		fc.Check(
			saveGitlabConnectorType(userInput.Data.GitlabConnectorType, existingGitConfig.GitlabConnectorType, frontend),
		)
	}

	if !enterAll {
		return fc.Err
	}

	// EXTENDED CONFIGURATION
	// keep-sorted start block=yes
	if configFile.AutoSync.IsNone() {
		fc.Check(
			saveAutoSync(userInput.Data.AutoSync, existingGitConfig.AutoSync, frontend),
		)
	}
	if configFile.BranchPrefix.IsNone() {
		fc.Check(
			saveBranchPrefix(userInput.Data.BranchPrefix, existingGitConfig.BranchPrefix, frontend),
		)
	}
	if configFile.ContributionRegex.IsNone() {
		fc.Check(
			saveContributionRegex(userInput.Data.ContributionRegex, existingGitConfig.ContributionRegex, frontend),
		)
	}
	if configFile.Detached.IsNone() {
		fc.Check(
			saveDetached(userInput.Data.Detached, existingGitConfig.Detached, frontend),
		)
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(
			saveFeatureRegex(userInput.Data.FeatureRegex, existingGitConfig.FeatureRegex, frontend),
		)
	}
	if configFile.IgnoreUncommitted.IsNone() {
		fc.Check(
			saveIgnoreUncommitted(userInput.Data.IgnoreUncommitted, existingGitConfig.IgnoreUncommitted, frontend),
		)
	}
	if configFile.NewBranchType.IsNone() {
		fc.Check(
			saveNewBranchType(userInput.Data.NewBranchType, existingGitConfig.NewBranchType, frontend),
		)
	}
	if configFile.ObservedRegex.IsNone() {
		fc.Check(
			saveObservedRegex(userInput.Data.ObservedRegex, existingGitConfig.ObservedRegex, frontend),
		)
	}
	if configFile.Order.IsNone() {
		fc.Check(
			saveOrder(userInput.Data.Order, existingGitConfig.Order, frontend),
		)
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(
			savePerennialRegex(userInput.Data.PerennialRegex, existingGitConfig.PerennialRegex, frontend),
		)
	}
	if configFile.ProposalsShowLineage.IsNone() {
		fc.Check(
			saveProposalsShowLineage(userInput.Data.ProposalsShowLineage, existingGitConfig.ProposalsShowLineage, frontend),
		)
	}
	if configFile.PushBranches.IsNone() {
		fc.Check(
			savePushBranches(userInput.Data.PushBranches, existingGitConfig.PushBranches, frontend),
		)
	}
	if configFile.PushHook.IsNone() {
		fc.Check(
			savePushHook(userInput.Data.PushHook, existingGitConfig.PushHook, frontend),
		)
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(
			saveShareNewBranches(userInput.Data.ShareNewBranches, existingGitConfig.ShareNewBranches, frontend),
		)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(
			saveShipDeleteTrackingBranch(userInput.Data.ShipDeleteTrackingBranch, existingGitConfig.ShipDeleteTrackingBranch, frontend),
		)
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(
			saveShipStrategy(userInput.Data.ShipStrategy, existingGitConfig.ShipStrategy, frontend),
		)
	}
	if configFile.Stash.IsNone() {
		fc.Check(
			saveStash(userInput.Data.Stash, existingGitConfig.Stash, frontend),
		)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(
			saveSyncFeatureStrategy(userInput.Data.SyncFeatureStrategy, existingGitConfig.SyncFeatureStrategy, frontend),
		)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(
			saveSyncPerennialStrategy(userInput.Data.SyncPerennialStrategy, existingGitConfig.SyncPerennialStrategy, frontend),
		)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(
			saveSyncPrototypeStrategy(userInput.Data.SyncPrototypeStrategy, existingGitConfig.SyncPrototypeStrategy, frontend),
		)
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(
			saveSyncTags(userInput.Data.SyncTags, existingGitConfig.SyncTags, frontend),
		)
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(
			saveSyncUpstream(userInput.Data.SyncUpstream, existingGitConfig.SyncUpstream, frontend),
		)
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(
			saveUnknownBranchType(userInput.Data.UnknownBranchType, existingGitConfig.UnknownBranchType, frontend),
		)
	}
	// keep-sorted end
	return fc.Err
}

// TODO: simplify to shorter version: valueToWriteToGit, valueAlreadyInGit Option[configdomain.AutoSync]
func saveAutoSync(valueToWriteToGit Option[configdomain.AutoSync], valueAlreadyInGit Option[configdomain.AutoSync], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetAutoSync(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveAutoSync(runner)
	return nil
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

func saveBranchPrefix(valueToWriteToGit Option[configdomain.BranchPrefix], valueAlreadyInGit Option[configdomain.BranchPrefix], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetBranchPrefix(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveBranchPrefix(runner)
	return nil
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

func saveDetached(valueToWriteToGit Option[configdomain.Detached], valueAlreadyInGit Option[configdomain.Detached], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetDetached(runner, value, configdomain.ConfigScopeLocal)
	}
	_ = gitconfig.RemoveDetached(runner)
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

func saveForgeType(valueToWriteToGit Option[forgedomain.ForgeType], valueAlreadyInGit Option[forgedomain.ForgeType], frontend subshelldomain.Runner) error {
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

func saveForgejoToken(valueToWriteToGit Option[forgedomain.ForgejoToken], valueAlreadyInGit Option[forgedomain.ForgejoToken], scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetForgejoToken(frontend, value, scope)
	}
	return gitconfig.RemoveForgejoToken(frontend)
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

func saveGithubConnectorType(valueToWriteToGit Option[forgedomain.GithubConnectorType], valueAlreadyInGit Option[forgedomain.GithubConnectorType], frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGithubConnectorType(frontend, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveGithubConnectorType(frontend)
}

func saveGithubToken(valueToWriteToGit Option[forgedomain.GithubToken], valueAlreadyInGit Option[forgedomain.GithubToken], scope configdomain.ConfigScope, githubConnectorType Option[forgedomain.GithubConnectorType], frontend subshelldomain.Runner) error {
	if connectorType, has := githubConnectorType.Get(); has {
		if connectorType == forgedomain.GithubConnectorTypeGh {
			return nil
		}
	}
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGithubToken(frontend, value, scope)
	}
	return gitconfig.RemoveGithubToken(frontend)
}

func saveGitlabConnectorType(valueToWriteToGit Option[forgedomain.GitlabConnectorType], valueAlreadyInGit Option[forgedomain.GitlabConnectorType], frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGitlabConnectorType(frontend, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveGitlabConnectorType(frontend)
}

func saveGitlabToken(valueToWriteToGit Option[forgedomain.GitlabToken], valueAlreadyInGit Option[forgedomain.GitlabToken], scope configdomain.ConfigScope, gitlabConnectorType Option[forgedomain.GitlabConnectorType], frontend subshelldomain.Runner) error {
	if connectorType, has := gitlabConnectorType.Get(); has {
		if connectorType == forgedomain.GitlabConnectorTypeGlab {
			return nil
		}
	}
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetGitlabToken(frontend, value, scope)
	}
	return gitconfig.RemoveGitlabToken(frontend)
}

func saveIgnoreUncommitted(valueToWriteToGit Option[configdomain.IgnoreUncommitted], valueAlreadyInGit Option[configdomain.IgnoreUncommitted], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetIgnoreUncommitted(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveIgnoreUncommitted(runner)
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

func saveOrder(valueToWriteToGit Option[configdomain.Order], valueAlreadyInGit Option[configdomain.Order], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetOrder(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveOrder(runner)
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

func saveProposalsShowLineage(valueToWriteToGit Option[forgedomain.ProposalBreadcrumb], valueAlreadyInGit Option[forgedomain.ProposalBreadcrumb], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetProposalsShowLineage(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveProposalsShowLineage(runner)
}

func savePushBranches(valueToWriteToGit Option[configdomain.PushBranches], valueAlreadyInGit Option[configdomain.PushBranches], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetPushBranches(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemovePushBranches(runner)
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

func saveStash(valueToWriteToGit Option[configdomain.Stash], valueAlreadyInGit Option[configdomain.Stash], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetStash(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveStash(runner)
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
