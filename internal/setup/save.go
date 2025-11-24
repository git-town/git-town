package setup

import (
	"cmp"
	"slices"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func Save(userInput UserInput, unvalidatedConfig config.UnvalidatedConfig, data Data, enterAll bool, frontend subshelldomain.Runner) error {
	errAliases := saveAliases(userInput.Data.Aliases, unvalidatedConfig.GitGlobal.Aliases, frontend)
	var (
		// keep-sorted start
		errBitbucketAppPassword error
		errBitbucketUsername    error
		errForgejoToken         error
		errGitHubToken          error
		errGitLabToken          error
		errGiteaToken           error
	// keep-sorted end
	)
	if forgeType, hasForgeType := userInput.DeterminedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeAzureDevOps:
			// no API token for now
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			errBitbucketUsername = saveBitbucketUsername(userInput.Data.BitbucketUsername, unvalidatedConfig.GitLocal.BitbucketUsername, userInput.Scope, frontend)
			errBitbucketAppPassword = saveBitbucketAppPassword(userInput.Data.BitbucketAppPassword, unvalidatedConfig.GitLocal.BitbucketAppPassword, userInput.Scope, frontend)
		case forgedomain.ForgeTypeForgejo:
			errForgejoToken = saveForgejoToken(userInput.Data.ForgejoToken, unvalidatedConfig.GitLocal.ForgejoToken, userInput.Scope, frontend)
		case forgedomain.ForgeTypeGitHub:
			errGitHubToken = saveGitHubToken(userInput.Data.GitHubToken, unvalidatedConfig.GitLocal.GitHubToken, userInput.Scope, userInput.Data.GitHubConnectorType, frontend)
		case forgedomain.ForgeTypeGitLab:
			errGitLabToken = saveGitLabToken(userInput.Data.GitLabToken, unvalidatedConfig.GitLocal.GitLabToken, userInput.Scope, userInput.Data.GitLabConnectorType, frontend)
		case forgedomain.ForgeTypeGitea:
			errGiteaToken = saveGiteaToken(userInput.Data.GiteaToken, unvalidatedConfig.GitLocal.GiteaToken, userInput.Scope, frontend)
		}
	}
	err := cmp.Or(
		// keep-sorted start
		errAliases,
		errBitbucketAppPassword,
		errBitbucketUsername,
		errForgejoToken,
		errGitHubToken,
		errGitLabToken,
		errGiteaToken,
		// keep-sorted end
	)
	if err != nil {
		return err
	}
	switch userInput.StorageLocation {
	case dialog.ConfigStorageOptionFile:
		return saveAllToFile(userInput, unvalidatedConfig.File, unvalidatedConfig.GitLocal, frontend)
	case dialog.ConfigStorageOptionGit:
		return saveAllToGit(userInput, unvalidatedConfig.GitLocal, unvalidatedConfig.File, data, enterAll, frontend)
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

func saveAllToFile(userInput UserInput, existingConfigFile configdomain.PartialConfig, gitConfig configdomain.PartialConfig, runner subshelldomain.Runner) error {
	userInput.Data.MainBranch = Some(userInput.ValidatedConfig.MainBranch)
	configData := existingConfigFile.Merge(userInput.Data)
	if err := configfile.Save(configData); err != nil {
		return err
	}
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
	if len(gitConfig.PerennialBranches) > 0 {
		_ = gitconfig.RemovePerennialBranches(runner)
	}
	if gitConfig.PerennialRegex.IsSome() {
		_ = gitconfig.RemovePerennialRegex(runner)
	}
	if gitConfig.ShareNewBranches.IsSome() {
		_ = gitconfig.RemoveShareNewBranches(runner)
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
	if gitConfig.ShipStrategy.IsSome() {
		_ = gitconfig.RemoveShipStrategy(runner)
	}
	if gitConfig.ShipDeleteTrackingBranch.IsSome() {
		_ = gitconfig.RemoveShipDeleteTrackingBranch(runner)
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
	if gitConfig.SyncUpstream.IsSome() {
		_ = gitconfig.RemoveSyncUpstream(runner)
	}
	if gitConfig.SyncTags.IsSome() {
		_ = gitconfig.RemoveSyncTags(runner)
	}
	if gitConfig.UnknownBranchType.IsSome() {
		_ = gitconfig.RemoveUnknownBranchType(runner)
	}
	if err := saveUnknownBranchType(userInput.Data.UnknownBranchType, gitConfig.UnknownBranchType, runner); err != nil {
		return err
	}
	return saveFeatureRegex(userInput.Data.FeatureRegex, gitConfig.FeatureRegex, runner)
}

func saveAllToGit(userInput UserInput, existingGitConfig configdomain.PartialConfig, configFile configdomain.PartialConfig, data Data, enterAll bool, frontend subshelldomain.Runner) error {
	var (
		// keep-sorted start
		errForgeType           error
		errGitHubConnectorType error
		errGitLabConnectorType error
		errMainBranch          error
		errOriginHostname      error
		errPerennialBranches   error
		errRemotes             error
	// keep-sorted end
	)

	// BASIC CONFIGURATION
	if configFile.MainBranch.IsNone() {
		errMainBranch = saveMainBranch(userInput.Data.MainBranch, existingGitConfig.MainBranch, frontend)
	}
	errPerennialBranches = savePerennialBranches(userInput.Data.PerennialBranches, existingGitConfig.PerennialBranches, frontend)
	if len(data.Remotes) > 1 && configFile.DevRemote.IsNone() {
		errRemotes = saveDevRemote(userInput.Data.DevRemote, existingGitConfig.DevRemote, frontend)
	}
	if configFile.HostingOriginHostname.IsNone() {
		errOriginHostname = saveOriginHostname(userInput.Data.HostingOriginHostname, existingGitConfig.HostingOriginHostname, frontend)
	}
	if configFile.ForgeType.IsNone() {
		errForgeType = saveForgeType(userInput.Data.ForgeType, existingGitConfig.ForgeType, frontend)
	}
	if configFile.GitHubConnectorType.IsNone() {
		errGitHubConnectorType = saveGitHubConnectorType(userInput.Data.GitHubConnectorType, existingGitConfig.GitHubConnectorType, frontend)
	}
	if configFile.GitLabConnectorType.IsNone() {
		errGitLabConnectorType = saveGitLabConnectorType(userInput.Data.GitLabConnectorType, existingGitConfig.GitLabConnectorType, frontend)
	}

	if !enterAll {
		return cmp.Or(
			// keep-sorted start
			errForgeType,
			errGitHubConnectorType,
			errGitLabConnectorType,
			errMainBranch,
			errOriginHostname,
			errPerennialBranches,
			errRemotes,
			// keep-sorted end
		)
	}

	// EXTENDED CONFIGURATION
	// keep-sorted start
	var (
		errAutoSync                 error
		errBranchPrefix             error
		errContributionRegex        error
		errDetached                 error
		errFeatureRegex             error
		errNewBranchType            error
		errObservedRegex            error
		errOrder                    error
		errPerennialRegex           error
		errProposalsShowLineage     error
		errPushBranches             error
		errPushHook                 error
		errShareNewBranches         error
		errShipDeleteTrackingBranch error
		errShipStrategy             error
		errStash                    error
		errSyncFeatureStrategy      error
		errSyncPerennialStrategy    error
		errSyncPrototypeStrategy    error
		errSyncTags                 error
		errSyncUpstream             error
		errUnknownBranchType        error
	)
	// keep-sorted end

	// TODO: sort this alphabetically
	if configFile.AutoSync.IsNone() {
		errAutoSync = saveAutoSync(userInput.Data.AutoSync, existingGitConfig.AutoSync, frontend)
	}
	if configFile.BranchPrefix.IsNone() {
		errBranchPrefix = saveBranchPrefix(userInput.Data.BranchPrefix, existingGitConfig.BranchPrefix, frontend)
	}
	if configFile.Detached.IsNone() {
		errDetached = saveDetached(userInput.Data.Detached, existingGitConfig.Detached, frontend)
	}
	if configFile.NewBranchType.IsNone() {
		errNewBranchType = saveNewBranchType(userInput.Data.NewBranchType, existingGitConfig.NewBranchType, frontend)
	}
	if configFile.PerennialRegex.IsNone() {
		errPerennialRegex = savePerennialRegex(userInput.Data.PerennialRegex, existingGitConfig.PerennialRegex, frontend)
	}
	if configFile.UnknownBranchType.IsNone() {
		errUnknownBranchType = saveUnknownBranchType(userInput.Data.UnknownBranchType, existingGitConfig.UnknownBranchType, frontend)
	}
	if configFile.FeatureRegex.IsNone() {
		errFeatureRegex = saveFeatureRegex(userInput.Data.FeatureRegex, existingGitConfig.FeatureRegex, frontend)
	}
	if configFile.ContributionRegex.IsNone() {
		errContributionRegex = saveContributionRegex(userInput.Data.ContributionRegex, existingGitConfig.ContributionRegex, frontend)
	}
	if configFile.ObservedRegex.IsNone() {
		errObservedRegex = saveObservedRegex(userInput.Data.ObservedRegex, existingGitConfig.ObservedRegex, frontend)
	}
	if configFile.Order.IsNone() {
		errOrder = saveOrder(userInput.Data.Order, existingGitConfig.Order, frontend)
	}
	if configFile.ProposalsShowLineage.IsNone() {
		errProposalsShowLineage = saveProposalsShowLineage(userInput.Data.ProposalsShowLineage, existingGitConfig.ProposalsShowLineage, frontend)
	}
	if configFile.PushBranches.IsNone() {
		errPushBranches = savePushBranches(userInput.Data.PushBranches, existingGitConfig.PushBranches, frontend)
	}
	if configFile.PushHook.IsNone() {
		errPushHook = savePushHook(userInput.Data.PushHook, existingGitConfig.PushHook, frontend)
	}
	if configFile.ShareNewBranches.IsNone() {
		errShareNewBranches = saveShareNewBranches(userInput.Data.ShareNewBranches, existingGitConfig.ShareNewBranches, frontend)
	}
	if configFile.ShipStrategy.IsNone() {
		errShipStrategy = saveShipStrategy(userInput.Data.ShipStrategy, existingGitConfig.ShipStrategy, frontend)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		errShipDeleteTrackingBranch = saveShipDeleteTrackingBranch(userInput.Data.ShipDeleteTrackingBranch, existingGitConfig.ShipDeleteTrackingBranch, frontend)
	}
	if configFile.Stash.IsNone() {
		errStash = saveStash(userInput.Data.Stash, existingGitConfig.Stash, frontend)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		errSyncFeatureStrategy = saveSyncFeatureStrategy(userInput.Data.SyncFeatureStrategy, existingGitConfig.SyncFeatureStrategy, frontend)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		errSyncPerennialStrategy = saveSyncPerennialStrategy(userInput.Data.SyncPerennialStrategy, existingGitConfig.SyncPerennialStrategy, frontend)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		errSyncPrototypeStrategy = saveSyncPrototypeStrategy(userInput.Data.SyncPrototypeStrategy, existingGitConfig.SyncPrototypeStrategy, frontend)
	}
	if configFile.SyncUpstream.IsNone() {
		errSyncUpstream = saveSyncUpstream(userInput.Data.SyncUpstream, existingGitConfig.SyncUpstream, frontend)
	}
	if configFile.SyncTags.IsNone() {
		errSyncTags = saveSyncTags(userInput.Data.SyncTags, existingGitConfig.SyncTags, frontend)
	}
	return cmp.Or(
		// keep-sorted start
		errAutoSync,
		errBranchPrefix,
		errContributionRegex,
		errDetached,
		errFeatureRegex,
		errNewBranchType,
		errObservedRegex,
		errOrder,
		errPerennialRegex,
		errProposalsShowLineage,
		errPushBranches,
		errPushHook,
		errShareNewBranches,
		errShipDeleteTrackingBranch,
		errShipStrategy,
		errStash,
		errSyncFeatureStrategy,
		errSyncPerennialStrategy,
		errSyncPrototypeStrategy,
		errSyncTags,
		errSyncUpstream,
		errUnknownBranchType,
		// keep-sorted end
	)
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

func saveForgejoToken(valueToWriteToGit Option[forgedomain.ForgejoToken], valueAlreadyInGit Option[forgedomain.ForgejoToken], scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetForgejoToken(frontend, value, scope)
	}
	return gitconfig.RemoveForgejoToken(frontend)
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

func saveProposalsShowLineage(valueToWriteToGit Option[forgedomain.ProposalsShowLineage], valueAlreadyInGit Option[forgedomain.ProposalsShowLineage], runner subshelldomain.Runner) error {
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
