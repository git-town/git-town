package setup

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func Enter(data Data) (UserInput, dialogdomain.Exit, error) {
	var emptyResult UserInput
	exit, err := dialog.Welcome(data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	aliases, exit, err := dialog.Aliases(configdomain.AllAliasableCommands(), data.Config.NormalConfig.Aliases, data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	mainBranchSetting, actualMainBranch, exit, err := enterMainBranch(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	perennialBranches, exit, err := enterPerennialBranches(data, actualMainBranch)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	perennialRegex, exit, err := enterPerennialRegex(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	featureRegex, exit, err := enterFeatureRegex(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	contributionRegex, exit, err := enterContributionRegex(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	observedRegex, exit, err := enterObservedRegex(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	newBranchType, exit, err := enterNewBranchType(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	unknownBranchType, exit, err := enterUnknownBranchType(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	devRemote, exit, err := enterDevRemote(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
EnterForgeData:
	hostingOriginHostName, exit, err := enterOriginHostName(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	enteredForgeType, exit, err := enterForgeType(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	devURL := data.Config.NormalConfig.DevURL(data.Backend)
	actualForgeType := determineForgeType(enteredForgeType.Or(data.Config.File.ForgeType), devURL)
	bitbucketUsername := None[forgedomain.BitbucketUsername]()
	bitbucketAppPassword := None[forgedomain.BitbucketAppPassword]()
	codebergToken := None[forgedomain.CodebergToken]()
	giteaToken := None[forgedomain.GiteaToken]()
	githubConnectorTypeOpt := None[forgedomain.GitHubConnectorType]()
	githubToken := None[forgedomain.GitHubToken]()
	gitlabConnectorTypeOpt := None[forgedomain.GitLabConnectorType]()
	gitlabToken := None[forgedomain.GitLabToken]()
	if forgeType, hasForgeType := actualForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			bitbucketUsername, exit, err = enterBitbucketUserName(data)
			if err != nil || exit {
				return emptyResult, exit, err
			}
			bitbucketAppPassword, exit, err = enterBitbucketAppPassword(data)
		case forgedomain.ForgeTypeCodeberg:
			codebergToken, exit, err = enterCodebergToken(data)
		case forgedomain.ForgeTypeGitea:
			giteaToken, exit, err = enterGiteaToken(data)
		case forgedomain.ForgeTypeGitHub:
			githubConnectorTypeOpt, exit, err = enterGitHubConnectorType(data)
			if err != nil || exit {
				return emptyResult, exit, err
			}
			if githubConnectorType, has := githubConnectorTypeOpt.Get(); has {
				switch githubConnectorType {
				case forgedomain.GitHubConnectorTypeAPI:
					githubToken, exit, err = enterGitHubToken(data)
				case forgedomain.GitHubConnectorTypeGh:
				}
			}
		case forgedomain.ForgeTypeGitLab:
			gitlabConnectorTypeOpt, exit, err = enterGitLabConnectorType(data)
			if err != nil || exit {
				return emptyResult, exit, err
			}
			if gitlabConnectorType, has := gitlabConnectorTypeOpt.Get(); has {
				switch gitlabConnectorType {
				case forgedomain.GitLabConnectorTypeAPI:
					gitlabToken, exit, err = enterGitLabToken(data)
				case forgedomain.GitLabConnectorTypeGlab:
				}
			}
		}
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	repeat, exit, err := testForgeAuth(testForgeAuthArgs{
		backend:              data.Backend,
		bitbucketAppPassword: bitbucketAppPassword,
		bitbucketUsername:    bitbucketUsername,
		codebergToken:        codebergToken,
		devURL:               devURL,
		forgeTypeOpt:         actualForgeType,
		giteaToken:           giteaToken,
		githubConnectorType:  githubConnectorTypeOpt,
		githubToken:          githubToken,
		gitlabConnectorType:  gitlabConnectorTypeOpt,
		gitlabToken:          gitlabToken,
		inputs:               data.Inputs,
		remoteURL:            data.Config.NormalConfig.RemoteURL(data.Backend, devRemote.GetOrElse(config.DefaultNormalConfig().DevRemote)),
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	if repeat {
		goto EnterForgeData
	}
	tokenScope, exit, err := enterTokenScope(enterTokenScopeArgs{
		bitbucketAppPassword: bitbucketAppPassword,
		bitbucketUsername:    bitbucketUsername,
		codebergToken:        codebergToken,
		data:                 data,
		determinedForgeType:  actualForgeType,
		existingConfig:       data.Config.NormalConfig,
		giteaToken:           giteaToken,
		githubToken:          githubToken,
		gitlabToken:          gitlabToken,
		inputs:               data.Inputs,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncFeatureStrategy, exit, err := enterSyncFeatureStrategy(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncPerennialStrategy, exit, err := enterSyncPerennialStrategy(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncPrototypeStrategy, exit, err := enterSyncPrototypeStrategy(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncUpstream, exit, err := enterSyncUpstream(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncTags, exit, err := enterSyncTags(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	shareNewBranches, exit, err := enterShareNewBranches(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	pushHook, exit, err := enterPushHook(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	shipStrategy, exit, err := enterShipStrategy(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	shipDeleteTrackingBranch, exit, err := enterShipDeleteTrackingBranch(data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	configStorage, exit, err := dialog.ConfigStorage(data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	normalData := configdomain.PartialConfig{
		Aliases:                  aliases,
		AutoResolve:              None[configdomain.AutoResolve](),
		BitbucketAppPassword:     bitbucketAppPassword,
		BitbucketUsername:        bitbucketUsername,
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{}, // the setup assistant doesn't ask for this
		CodebergToken:            codebergToken,
		ContributionRegex:        contributionRegex,
		DevRemote:                devRemote,
		DryRun:                   None[configdomain.DryRun](), // the setup assistant doesn't ask for this
		FeatureRegex:             featureRegex,
		ForgeType:                enteredForgeType,
		GitHubConnectorType:      githubConnectorTypeOpt,
		GitHubToken:              githubToken,
		GitLabConnectorType:      gitlabConnectorTypeOpt,
		GitLabToken:              gitlabToken,
		GitUserEmail:             None[gitdomain.GitUserEmail](),
		GitUserName:              None[gitdomain.GitUserName](),
		GiteaToken:               giteaToken,
		HostingOriginHostname:    hostingOriginHostName,
		Lineage:                  configdomain.Lineage{}, // the setup assistant doesn't ask for this
		MainBranch:               mainBranchSetting,
		NewBranchType:            newBranchType,
		ObservedRegex:            observedRegex,
		Offline:                  None[configdomain.Offline](), // the setup assistant doesn't ask for this
		PerennialBranches:        perennialBranches,
		PerennialRegex:           perennialRegex,
		ProposalsShowLineage:     None[forgedomain.ProposalsShowLineage](), // TODO: populate this in the setup assistant once https://github.com/git-town/git-town/issues/3003 is shipped
		PushHook:                 pushHook,
		ShareNewBranches:         shareNewBranches,
		ShipDeleteTrackingBranch: shipDeleteTrackingBranch,
		ShipStrategy:             shipStrategy,
		SyncFeatureStrategy:      syncFeatureStrategy,
		SyncPerennialStrategy:    syncPerennialStrategy,
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 syncTags,
		SyncUpstream:             syncUpstream,
		UnknownBranchType:        unknownBranchType,
		Verbose:                  None[configdomain.Verbose](), // the setup assistant doesn't ask for this
	}
	validatedData := configdomain.ValidatedConfigData{
		GitUserEmail: "", // the setup assistant doesn't ask for this
		GitUserName:  "", // the setup assistant doesn't ask for this
		MainBranch:   actualMainBranch,
	}
	return UserInput{normalData, actualForgeType, tokenScope, configStorage, validatedData}, false, nil
}

// data entered by the user in the setup assistant
type UserInput struct {
	Data                configdomain.PartialConfig
	DeterminedForgeType Option[forgedomain.ForgeType] // the forge type that was determined by the setup assistant - not necessarily what the user entered (could also be "auto detect")
	Scope               configdomain.ConfigScope
	StorageLocation     dialog.ConfigStorageOption
	ValidatedConfig     configdomain.ValidatedConfigData
}

func determineExistingScope(configSnapshot undoconfig.BeginConfigSnapshot, key configdomain.Key, oldValue fmt.Stringer) configdomain.ConfigScope {
	switch {
	case oldValue.String() == "":
		return configdomain.ConfigScopeLocal
	case configSnapshot.Global[key] == oldValue.String():
		return configdomain.ConfigScopeGlobal
	case configSnapshot.Local[key] == oldValue.String():
		return configdomain.ConfigScopeLocal
	default:
		return configdomain.ConfigScopeLocal
	}
}

func determineForgeType(userChoice Option[forgedomain.ForgeType], devURL Option[giturl.Parts]) Option[forgedomain.ForgeType] {
	if userChoice.IsSome() {
		return userChoice
	}
	if devURL, hasDevURL := devURL.Get(); hasDevURL {
		return forge.Detect(devURL, userChoice)
	}
	return None[forgedomain.ForgeType]()
}

func enterBitbucketAppPassword(data Data) (Option[forgedomain.BitbucketAppPassword], dialogdomain.Exit, error) {
	if data.Config.File.BitbucketUsername.IsSome() {
		return None[forgedomain.BitbucketAppPassword](), false, nil
	}
	return dialog.BitbucketAppPassword(dialog.Args[forgedomain.BitbucketAppPassword]{
		Global: data.Config.GitLocal.BitbucketAppPassword,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.BitbucketAppPassword,
	})
}

func enterBitbucketUserName(data Data) (Option[forgedomain.BitbucketUsername], dialogdomain.Exit, error) {
	if data.Config.File.BitbucketUsername.IsSome() {
		return None[forgedomain.BitbucketUsername](), false, nil
	}
	return dialog.BitbucketUsername(dialog.Args[forgedomain.BitbucketUsername]{
		Global: data.Config.GitLocal.BitbucketUsername,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.BitbucketUsername,
	})
}

func enterCodebergToken(data Data) (Option[forgedomain.CodebergToken], dialogdomain.Exit, error) {
	if data.Config.File.CodebergToken.IsSome() {
		return None[forgedomain.CodebergToken](), false, nil
	}
	return dialog.CodebergToken(dialog.Args[forgedomain.CodebergToken]{
		Global: data.Config.GitGlobal.CodebergToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.CodebergToken,
	})
}

func enterContributionRegex(data Data) (Option[configdomain.ContributionRegex], dialogdomain.Exit, error) {
	if data.Config.File.ContributionRegex.IsSome() {
		return None[configdomain.ContributionRegex](), false, nil
	}
	return dialog.ContributionRegex(dialog.Args[configdomain.ContributionRegex]{
		Global: data.Config.GitGlobal.ContributionRegex,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ContributionRegex,
	})
}

func enterDevRemote(data Data) (Option[gitdomain.Remote], dialogdomain.Exit, error) {
	if data.Config.File.DevRemote.IsSome() {
		return None[gitdomain.Remote](), false, nil
	}
	return dialog.DevRemote(data.Remotes, dialog.Args[gitdomain.Remote]{
		Global: data.Config.GitGlobal.DevRemote,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.DevRemote,
	})
}

func enterFeatureRegex(data Data) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	if data.Config.File.FeatureRegex.IsSome() {
		return None[configdomain.FeatureRegex](), false, nil
	}
	return dialog.FeatureRegex(dialog.Args[configdomain.FeatureRegex]{
		Global: data.Config.GitGlobal.FeatureRegex,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.FeatureRegex,
	})
}

func enterForgeType(data Data) (Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	if data.Config.File.ForgeType.IsSome() {
		return None[forgedomain.ForgeType](), false, nil
	}
	return dialog.ForgeType(dialog.Args[forgedomain.ForgeType]{
		Global: data.Config.GitGlobal.ForgeType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ForgeType,
	})
}

func enterGitHubConnectorType(data Data) (Option[forgedomain.GitHubConnectorType], dialogdomain.Exit, error) {
	if data.Config.File.GitHubConnectorType.IsSome() {
		return None[forgedomain.GitHubConnectorType](), false, nil
	}
	return dialog.GitHubConnectorType(dialog.Args[forgedomain.GitHubConnectorType]{
		Global: data.Config.GitGlobal.GitHubConnectorType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GitHubConnectorType,
	})
}

func enterGitHubToken(data Data) (Option[forgedomain.GitHubToken], dialogdomain.Exit, error) {
	if data.Config.File.GitHubToken.IsSome() {
		return None[forgedomain.GitHubToken](), false, nil
	}
	return dialog.GitHubToken(dialog.Args[forgedomain.GitHubToken]{
		Global: data.Config.GitGlobal.GitHubToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GitHubToken,
	})
}

func enterGitLabConnectorType(data Data) (Option[forgedomain.GitLabConnectorType], dialogdomain.Exit, error) {
	if data.Config.File.GitLabConnectorType.IsSome() {
		return None[forgedomain.GitLabConnectorType](), false, nil
	}
	return dialog.GitLabConnectorType(dialog.Args[forgedomain.GitLabConnectorType]{
		Global: data.Config.GitGlobal.GitLabConnectorType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GitLabConnectorType,
	})
}

func enterGitLabToken(data Data) (Option[forgedomain.GitLabToken], dialogdomain.Exit, error) {
	if data.Config.File.GitLabToken.IsSome() {
		return None[forgedomain.GitLabToken](), false, nil
	}
	return dialog.GitLabToken(dialog.Args[forgedomain.GitLabToken]{
		Global: data.Config.GitGlobal.GitLabToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GitLabToken,
	})
}

func enterGiteaToken(data Data) (Option[forgedomain.GiteaToken], dialogdomain.Exit, error) {
	if data.Config.File.GiteaToken.IsSome() {
		return None[forgedomain.GiteaToken](), false, nil
	}
	return dialog.GiteaToken(dialog.Args[forgedomain.GiteaToken]{
		Global: data.Config.GitGlobal.GiteaToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GiteaToken,
	})
}

func enterMainBranch(data Data) (userChoice Option[gitdomain.LocalBranchName], actualMainBranch gitdomain.LocalBranchName, exit dialogdomain.Exit, err error) {
	if configFileMainBranch, hasMain := data.Config.File.MainBranch.Get(); hasMain {
		return Some(configFileMainBranch), configFileMainBranch, false, nil
	}
	return dialog.MainBranch(dialog.MainBranchArgs{
		Inputs:         data.Inputs,
		Local:          data.Config.GitLocal.MainBranch,
		LocalBranches:  data.LocalBranches,
		StandardBranch: data.Git.StandardBranch(data.Backend),
		Unscoped:       data.Config.GitUnscoped.MainBranch,
	})
}

func enterNewBranchType(data Data) (Option[configdomain.NewBranchType], dialogdomain.Exit, error) {
	if data.Config.File.NewBranchType.IsSome() {
		return None[configdomain.NewBranchType](), false, nil
	}
	return dialog.NewBranchType(dialog.Args[configdomain.NewBranchType]{
		Global: data.Config.GitGlobal.NewBranchType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.NewBranchType,
	})
}

func enterObservedRegex(data Data) (Option[configdomain.ObservedRegex], dialogdomain.Exit, error) {
	if data.Config.File.ObservedRegex.IsSome() {
		return None[configdomain.ObservedRegex](), false, nil
	}
	return dialog.ObservedRegex(dialog.Args[configdomain.ObservedRegex]{
		Global: data.Config.GitGlobal.ObservedRegex,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ObservedRegex,
	})
}

func enterOriginHostName(data Data) (Option[configdomain.HostingOriginHostname], dialogdomain.Exit, error) {
	if data.Config.File.HostingOriginHostname.IsSome() {
		return None[configdomain.HostingOriginHostname](), false, nil
	}
	return dialog.OriginHostname(dialog.Args[configdomain.HostingOriginHostname]{
		Global: data.Config.GitGlobal.HostingOriginHostname,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.HostingOriginHostname,
	})
}

func enterPerennialBranches(data Data, mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchNames, dialogdomain.Exit, error) {
	immutablePerennials := gitdomain.LocalBranchNames{mainBranch}.
		AppendAllMissing(data.Config.File.PerennialBranches...).
		AppendAllMissing(data.Config.GitGlobal.PerennialBranches...)
	return dialog.PerennialBranches(dialog.PerennialBranchesArgs{
		ImmutableGitPerennials: immutablePerennials,
		Inputs:                 data.Inputs,
		LocalBranches:          data.LocalBranches,
		LocalGitPerennials:     data.Config.GitLocal.PerennialBranches,
		MainBranch:             mainBranch,
	})
}

func enterPerennialRegex(data Data) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	if data.Config.File.PerennialRegex.IsSome() {
		return None[configdomain.PerennialRegex](), false, nil
	}
	return dialog.PerennialRegex(dialog.Args[configdomain.PerennialRegex]{
		Global: data.Config.GitGlobal.PerennialRegex,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.PerennialRegex,
	})
}

func enterPushHook(data Data) (Option[configdomain.PushHook], dialogdomain.Exit, error) {
	if data.Config.File.PushHook.IsSome() {
		return None[configdomain.PushHook](), false, nil
	}
	return dialog.PushHook(dialog.Args[configdomain.PushHook]{
		Global: data.Config.GitGlobal.PushHook,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.PushHook,
	})
}

func enterShareNewBranches(data Data) (Option[configdomain.ShareNewBranches], dialogdomain.Exit, error) {
	if data.Config.File.ShareNewBranches.IsSome() {
		return None[configdomain.ShareNewBranches](), false, nil
	}
	return dialog.ShareNewBranches(dialog.Args[configdomain.ShareNewBranches]{
		Global: data.Config.GitGlobal.ShareNewBranches,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ShareNewBranches,
	})
}

func enterShipDeleteTrackingBranch(data Data) (Option[configdomain.ShipDeleteTrackingBranch], dialogdomain.Exit, error) {
	if data.Config.File.ShipDeleteTrackingBranch.IsSome() {
		return None[configdomain.ShipDeleteTrackingBranch](), false, nil
	}
	return dialog.ShipDeleteTrackingBranch(dialog.Args[configdomain.ShipDeleteTrackingBranch]{
		Global: data.Config.GitGlobal.ShipDeleteTrackingBranch,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ShipDeleteTrackingBranch,
	})
}

func enterShipStrategy(data Data) (Option[configdomain.ShipStrategy], dialogdomain.Exit, error) {
	if data.Config.File.ShipStrategy.IsSome() {
		return None[configdomain.ShipStrategy](), false, nil
	}
	return dialog.ShipStrategy(dialog.Args[configdomain.ShipStrategy]{
		Global: data.Config.GitGlobal.ShipStrategy,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ShipStrategy,
	})
}

func enterSyncFeatureStrategy(data Data) (Option[configdomain.SyncFeatureStrategy], dialogdomain.Exit, error) {
	if data.Config.File.SyncFeatureStrategy.IsSome() {
		return None[configdomain.SyncFeatureStrategy](), false, nil
	}
	return dialog.SyncFeatureStrategy(dialog.Args[configdomain.SyncFeatureStrategy]{
		Global: data.Config.GitGlobal.SyncFeatureStrategy,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.SyncFeatureStrategy,
	})
}

func enterSyncPerennialStrategy(data Data) (Option[configdomain.SyncPerennialStrategy], dialogdomain.Exit, error) {
	if data.Config.File.SyncPerennialStrategy.IsSome() {
		return None[configdomain.SyncPerennialStrategy](), false, nil
	}
	return dialog.SyncPerennialStrategy(dialog.Args[configdomain.SyncPerennialStrategy]{
		Global: data.Config.GitGlobal.SyncPerennialStrategy,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.SyncPerennialStrategy,
	})
}

func enterSyncPrototypeStrategy(data Data) (Option[configdomain.SyncPrototypeStrategy], dialogdomain.Exit, error) {
	if data.Config.File.SyncPrototypeStrategy.IsSome() {
		return None[configdomain.SyncPrototypeStrategy](), false, nil
	}
	return dialog.SyncPrototypeStrategy(dialog.Args[configdomain.SyncPrototypeStrategy]{
		Global: data.Config.GitGlobal.SyncPrototypeStrategy,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.SyncPrototypeStrategy,
	})
}

func enterSyncTags(data Data) (Option[configdomain.SyncTags], dialogdomain.Exit, error) {
	if data.Config.File.SyncTags.IsSome() {
		return None[configdomain.SyncTags](), false, nil
	}
	return dialog.SyncTags(dialog.Args[configdomain.SyncTags]{
		Global: data.Config.GitGlobal.SyncTags,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.SyncTags,
	})
}

func enterSyncUpstream(data Data) (Option[configdomain.SyncUpstream], dialogdomain.Exit, error) {
	if data.Config.File.SyncUpstream.IsSome() {
		return None[configdomain.SyncUpstream](), false, nil
	}
	return dialog.SyncUpstream(dialog.Args[configdomain.SyncUpstream]{
		Global: data.Config.GitGlobal.SyncUpstream,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.SyncUpstream,
	})
}

func enterTokenScope(args enterTokenScopeArgs) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if shouldAskForScope(args) {
		return tokenScopeDialog(args)
	}
	return configdomain.ConfigScopeLocal, false, nil
}

type enterTokenScopeArgs struct {
	bitbucketAppPassword Option[forgedomain.BitbucketAppPassword]
	bitbucketUsername    Option[forgedomain.BitbucketUsername]
	codebergToken        Option[forgedomain.CodebergToken]
	data                 Data
	determinedForgeType  Option[forgedomain.ForgeType]
	existingConfig       config.NormalConfig
	giteaToken           Option[forgedomain.GiteaToken]
	githubToken          Option[forgedomain.GitHubToken]
	gitlabToken          Option[forgedomain.GitLabToken]
	inputs               dialogcomponents.Inputs
}

func enterUnknownBranchType(data Data) (Option[configdomain.UnknownBranchType], dialogdomain.Exit, error) {
	if data.Config.File.UnknownBranchType.IsSome() {
		return None[configdomain.UnknownBranchType](), false, nil
	}
	return dialog.UnknownBranchType(dialog.Args[configdomain.UnknownBranchType]{
		Global: data.Config.GitGlobal.UnknownBranchType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.UnknownBranchType,
	})
}

func existsAndChanged[T fmt.Stringer](input, existing T) bool {
	return input.String() != "" && input.String() != existing.String()
}

func shouldAskForScope(args enterTokenScopeArgs) bool {
	if forgeType, hasForgeType := args.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			return existsAndChanged(args.bitbucketUsername, args.existingConfig.BitbucketUsername) &&
				existsAndChanged(args.bitbucketAppPassword, args.existingConfig.BitbucketAppPassword)
		case forgedomain.ForgeTypeCodeberg:
			return existsAndChanged(args.codebergToken, args.existingConfig.CodebergToken)
		case forgedomain.ForgeTypeGitea:
			return existsAndChanged(args.giteaToken, args.existingConfig.GiteaToken)
		case forgedomain.ForgeTypeGitHub:
			return existsAndChanged(args.githubToken, args.existingConfig.GitHubToken)
		case forgedomain.ForgeTypeGitLab:
			return existsAndChanged(args.gitlabToken, args.existingConfig.GitLabToken)
		}
	}
	return false
}

func testForgeAuth(args testForgeAuthArgs) (repeat bool, exit dialogdomain.Exit, err error) {
	if _, inTest := os.LookupEnv(subshell.TestToken); inTest {
		return false, false, nil
	}
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              args.backend,
		BitbucketAppPassword: args.bitbucketAppPassword,
		BitbucketUsername:    args.bitbucketUsername,
		CodebergToken:        args.codebergToken,
		ForgeType:            args.forgeTypeOpt,
		Frontend:             args.backend,
		GitHubConnectorType:  args.githubConnectorType,
		GitHubToken:          args.githubToken,
		GitLabConnectorType:  args.gitlabConnectorType,
		GitLabToken:          args.gitlabToken,
		GiteaToken:           args.giteaToken,
		Log:                  print.Logger{},
		RemoteURL:            args.devURL,
	})
	if err != nil {
		return false, false, err
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return false, false, nil
	}
	verifyResult := connector.VerifyConnection()
	if verifyResult.AuthenticationError != nil {
		return dialog.CredentialsNoAccess(verifyResult.AuthenticationError, args.inputs)
	}
	if user, hasUser := verifyResult.AuthenticatedUser.Get(); hasUser {
		fmt.Printf(messages.CredentialsForgeUserName, dialogcomponents.FormattedSelection(user, exit))
	}
	if verifyResult.AuthorizationError != nil {
		return dialog.CredentialsNoProposalAccess(verifyResult.AuthorizationError, args.inputs)
	}
	fmt.Println(messages.CredentialsAccess)
	return false, false, nil
}

type testForgeAuthArgs struct {
	backend              subshelldomain.RunnerQuerier
	bitbucketAppPassword Option[forgedomain.BitbucketAppPassword]
	bitbucketUsername    Option[forgedomain.BitbucketUsername]
	codebergToken        Option[forgedomain.CodebergToken]
	devURL               Option[giturl.Parts]
	forgeTypeOpt         Option[forgedomain.ForgeType]
	giteaToken           Option[forgedomain.GiteaToken]
	githubConnectorType  Option[forgedomain.GitHubConnectorType]
	githubToken          Option[forgedomain.GitHubToken]
	gitlabConnectorType  Option[forgedomain.GitLabConnectorType]
	gitlabToken          Option[forgedomain.GitLabToken]
	inputs               dialogcomponents.Inputs
	remoteURL            Option[giturl.Parts]
}

func tokenScopeDialog(args enterTokenScopeArgs) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if forgeType, hasForgeType := args.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyBitbucketUsername, args.data.Config.NormalConfig.BitbucketUsername)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeCodeberg:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyCodebergToken, args.data.Config.NormalConfig.CodebergToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitea:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyGiteaToken, args.data.Config.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitHub:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyGitHubToken, args.data.Config.NormalConfig.GitHubToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitLab:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyGitLabToken, args.data.Config.NormalConfig.GitLabToken)
			return dialog.TokenScope(existingScope, args.inputs)
		}
	}
	return configdomain.ConfigScopeLocal, false, nil
}
