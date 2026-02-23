package setup

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/git/giturl"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func Enter(data Data, configDir configdomain.RepoConfigDir) (UserInput, dialogdomain.Exit, bool, error) {
	var emptyResult UserInput
	exit, err := dialog.Welcome(data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	aliases, exit, err := dialog.Aliases(configdomain.AllAliasableCommands(), data.Config.NormalConfig.Aliases, data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	mainBranchResult, exit, err := enterMainBranch(data)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	perennialBranches, exit, err := enterPerennialBranches(data, mainBranchResult.ActualMainBranch)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
EnterForgeData:
	devRemote, exit, err := enterDevRemote(data)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	hostingOriginHostName, exit, err := enterOriginHostName(data)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	enteredForgeType, exit, err := enterForgeType(data)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	devURL := data.Config.NormalConfig.DevURL(data.Backend)
	actualForgeType := determineForgeType(enteredForgeType.Or(data.Config.File.ForgeType), devURL)
	bitbucketUsername := None[forgedomain.BitbucketUsername]()
	bitbucketAppPassword := None[forgedomain.BitbucketAppPassword]()
	forgejoToken := None[forgedomain.ForgejoToken]()
	giteaToken := None[forgedomain.GiteaToken]()
	githubConnectorTypeOpt := None[forgedomain.GithubConnectorType]()
	githubToken := None[forgedomain.GithubToken]()
	gitlabConnectorTypeOpt := None[forgedomain.GitlabConnectorType]()
	gitlabToken := None[forgedomain.GitlabToken]()
	if forgeType, hasForgeType := actualForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeAzuredevops:
			// the Azure DevOps connector doesn't have connectivity to the API implemented for now
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			bitbucketUsername, exit, err = enterBitbucketUserName(data)
			if err != nil || exit {
				return emptyResult, exit, false, err
			}
			bitbucketAppPassword, exit, err = enterBitbucketAppPassword(data)
			if err != nil || exit {
				return emptyResult, exit, false, err
			}
		case forgedomain.ForgeTypeForgejo:
			forgejoToken, exit, err = enterForgejoToken(data)
			if err != nil || exit {
				return emptyResult, exit, false, err
			}
		case forgedomain.ForgeTypeGitea:
			giteaToken, exit, err = enterGiteaToken(data)
			if err != nil || exit {
				return emptyResult, exit, false, err
			}
		case forgedomain.ForgeTypeGithub:
			githubConnectorTypeOpt, exit, err = enterGithubConnectorType(data)
			if err != nil || exit {
				return emptyResult, exit, false, err
			}
			if githubConnectorType, has := githubConnectorTypeOpt.Or(data.Config.File.GithubConnectorType).Get(); has {
				switch githubConnectorType {
				case forgedomain.GithubConnectorTypeAPI:
					githubToken, exit, err = enterGithubToken(data)
					if err != nil || exit {
						return emptyResult, exit, false, err
					}
				case forgedomain.GithubConnectorTypeGh:
				}
			}
		case forgedomain.ForgeTypeGitlab:
			gitlabConnectorTypeOpt, exit, err = enterGitlabConnectorType(data)
			if err != nil || exit {
				return emptyResult, exit, false, err
			}
			if gitlabConnectorType, has := gitlabConnectorTypeOpt.Or(data.Config.File.GitlabConnectorType).Get(); has {
				switch gitlabConnectorType {
				case forgedomain.GitlabConnectorTypeAPI:
					gitlabToken, exit, err = enterGitlabToken(data)
					if err != nil || exit {
						return emptyResult, exit, false, err
					}
				case forgedomain.GitlabConnectorTypeGlab:
				}
			}
		}
	}
	flow, exit, err := testForgeAuth(testForgeAuthArgs{
		backend:              data.Backend,
		bitbucketAppPassword: bitbucketAppPassword.Or(data.Config.GitGlobal.BitbucketAppPassword),
		bitbucketUsername:    bitbucketUsername.Or(data.Config.GitGlobal.BitbucketUsername),
		configDir:            configDir,
		devURL:               devURL,
		forgeTypeOpt:         actualForgeType,
		forgejoToken:         forgejoToken.Or(data.Config.GitGlobal.ForgejoToken),
		giteaToken:           giteaToken.Or(data.Config.GitGlobal.GiteaToken),
		githubConnectorType:  githubConnectorTypeOpt.Or(data.Config.GitGlobal.GithubConnectorType),
		githubToken:          githubToken.Or(data.Config.GitGlobal.GithubToken),
		gitlabConnectorType:  gitlabConnectorTypeOpt.Or(data.Config.GitGlobal.GitlabConnectorType),
		gitlabToken:          gitlabToken.Or(data.Config.GitGlobal.GitlabToken),
		inputs:               data.Inputs,
		remoteURL:            data.Config.NormalConfig.RemoteURL(data.Backend, devRemote.GetOr(config.DefaultNormalConfig().DevRemote)),
	})
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	if flow == configdomain.ProgramFlowRestart {
		goto EnterForgeData
	}
	tokenScope, exit, err := enterTokenScope(enterTokenScopeArgs{
		bitbucketAppPassword: bitbucketAppPassword,
		bitbucketUsername:    bitbucketUsername,
		data:                 data,
		determinedForgeType:  actualForgeType,
		existingConfig:       data.Config.NormalConfig,
		forgejoToken:         forgejoToken,
		giteaToken:           giteaToken,
		githubToken:          githubToken,
		gitlabToken:          gitlabToken,
		inputs:               data.Inputs,
	})
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	enterAll, exit, err := dialog.EnterAll(data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, enterAll, err
	}
	autoSync := None[configdomain.AutoSync]()
	perennialRegex := None[configdomain.PerennialRegex]()
	featureRegex := None[configdomain.FeatureRegex]()
	contributionRegex := None[configdomain.ContributionRegex]()
	observedRegex := None[configdomain.ObservedRegex]()
	branchPrefix := None[configdomain.BranchPrefix]()
	order := None[configdomain.Order]()
	newBranchType := None[configdomain.NewBranchType]()
	unknownBranchType := None[configdomain.UnknownBranchType]()
	syncFeatureStrategy := None[configdomain.SyncFeatureStrategy]()
	syncPerennialStrategy := None[configdomain.SyncPerennialStrategy]()
	syncPrototypeStrategy := None[configdomain.SyncPrototypeStrategy]()
	syncUpstream := None[configdomain.SyncUpstream]()
	syncTags := None[configdomain.SyncTags]()
	detached := None[configdomain.Detached]()
	stash := None[configdomain.Stash]()
	shareNewBranches := None[configdomain.ShareNewBranches]()
	pushBranches := None[configdomain.PushBranches]()
	pushHook := None[configdomain.PushHook]()
	shipStrategy := None[configdomain.ShipStrategy]()
	shipDeleteTrackingBranch := None[configdomain.ShipDeleteTrackingBranch]()
	ignoreUncommitted := None[configdomain.IgnoreUncommitted]()
	proposalBreadcrumb := None[configdomain.ProposalBreadcrumb]()
	if enterAll {
		perennialRegex, exit, err = enterPerennialRegex(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		featureRegex, exit, err = enterFeatureRegex(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		contributionRegex, exit, err = enterContributionRegex(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		observedRegex, exit, err = enterObservedRegex(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		branchPrefix, exit, err = enterBranchPrefix(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		newBranchType, exit, err = enterNewBranchType(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		unknownBranchType, exit, err = enterUnknownBranchType(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		syncFeatureStrategy, exit, err = enterSyncFeatureStrategy(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		syncPerennialStrategy, exit, err = enterSyncPerennialStrategy(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		syncPrototypeStrategy, exit, err = enterSyncPrototypeStrategy(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		syncUpstream, exit, err = enterSyncUpstream(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		autoSync, exit, err = enterAutoSync(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		syncTags, exit, err = enterSyncTags(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		detached, exit, err = enterDetached(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		stash, exit, err = enterStash(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		shareNewBranches, exit, err = enterShareNewBranches(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		pushBranches, exit, err = enterPushBranches(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		pushHook, exit, err = enterPushHook(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		shipStrategy, exit, err = enterShipStrategy(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		shipDeleteTrackingBranch, exit, err = enterShipDeleteTrackingBranch(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		ignoreUncommitted, exit, err = enterIgnoreUncommitted(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		order, exit, err = enterOrder(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
		proposalBreadcrumb, exit, err = enterProposalBreadcrumb(data)
		if err != nil || exit {
			return emptyResult, exit, false, err
		}
	}
	configStorage, exit, err := dialog.ConfigStorage(data.Inputs)
	if err != nil || exit {
		return emptyResult, exit, false, err
	}
	normalData := configdomain.PartialConfig{
		Aliases:                     aliases,
		AutoResolve:                 None[configdomain.AutoResolve](),
		AutoSync:                    autoSync,
		BitbucketAppPassword:        bitbucketAppPassword,
		BitbucketUsername:           bitbucketUsername,
		BranchPrefix:                branchPrefix,
		BranchTypeOverrides:         configdomain.BranchTypeOverrides{}, // the setup assistant doesn't ask for this
		Browser:                     None[configdomain.Browser](),
		ForgejoToken:                forgejoToken,
		ContributionRegex:           contributionRegex,
		Detached:                    detached,
		DevRemote:                   devRemote,
		DisplayTypes:                None[configdomain.DisplayTypes](),
		DryRun:                      None[configdomain.DryRun](), // the setup assistant doesn't ask for this
		FeatureRegex:                featureRegex,
		ForgeType:                   enteredForgeType,
		GithubConnectorType:         githubConnectorTypeOpt,
		GithubToken:                 githubToken,
		GitlabConnectorType:         gitlabConnectorTypeOpt,
		GitlabToken:                 gitlabToken,
		GitUserEmail:                None[gitdomain.GitUserEmail](),
		GitUserName:                 None[gitdomain.GitUserName](),
		GiteaToken:                  giteaToken,
		HostingOriginHostname:       hostingOriginHostName,
		IgnoreUncommitted:           ignoreUncommitted,
		Lineage:                     configdomain.NewLineage(), // the setup assistant doesn't ask for this
		MainBranch:                  mainBranchResult.UserChoice,
		NewBranchType:               newBranchType,
		ObservedRegex:               observedRegex,
		Offline:                     None[configdomain.Offline](), // the setup assistant doesn't ask for this
		Order:                       order,
		PerennialBranches:           perennialBranches,
		PerennialRegex:              perennialRegex,
		ProposalBreadcrumb:          proposalBreadcrumb,
		ProposalBreadcrumbDirection: None[configdomain.ProposalBreadcrumbDirection](), // TODO: add this to the setup assistant
		PushBranches:                pushBranches,
		PushHook:                    pushHook,
		ShareNewBranches:            shareNewBranches,
		ShipDeleteTrackingBranch:    shipDeleteTrackingBranch,
		ShipStrategy:                shipStrategy,
		Stash:                       stash,
		SyncFeatureStrategy:         syncFeatureStrategy,
		SyncPerennialStrategy:       syncPerennialStrategy,
		SyncPrototypeStrategy:       syncPrototypeStrategy,
		SyncTags:                    syncTags,
		SyncUpstream:                syncUpstream,
		UnknownBranchType:           unknownBranchType,
		Verbose:                     None[configdomain.Verbose](), // the setup assistant doesn't ask for this
	}
	validatedData := configdomain.ValidatedConfigData{
		MainBranch: mainBranchResult.ActualMainBranch,
	}
	return UserInput{normalData, actualForgeType, tokenScope, configStorage, validatedData}, false, enterAll, nil
}

// data entered by the user in the setup assistant
type UserInput struct {
	Data                configdomain.PartialConfig
	DeterminedForgeType Option[forgedomain.ForgeType] // the forge type that was determined by the setup assistant - not necessarily what the user entered (could also be "auto detect")
	Scope               configdomain.ConfigScope
	StorageLocation     dialog.ConfigStorageOption
	ValidatedConfig     configdomain.ValidatedConfigData
}

func determineExistingScope[T ~string](configSnapshot configdomain.BeginConfigSnapshot, key configdomain.Key, oldValueOpt Option[T]) configdomain.ConfigScope {
	oldValue, hasOldValue := oldValueOpt.Get()
	globalStr, hasGlobal := configSnapshot.Global[key]
	globalValue := T(globalStr)
	if hasOldValue && hasGlobal && globalValue == oldValue {
		return configdomain.ConfigScopeGlobal
	}
	return configdomain.ConfigScopeLocal
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

func enterAutoSync(data Data) (Option[configdomain.AutoSync], dialogdomain.Exit, error) {
	if data.Config.File.AutoSync.IsSome() {
		return None[configdomain.AutoSync](), false, nil
	}
	return dialog.AutoSync(dialog.Args[configdomain.AutoSync]{
		Global: data.Config.GitGlobal.AutoSync,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.AutoSync,
	})
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

func enterBranchPrefix(data Data) (Option[configdomain.BranchPrefix], dialogdomain.Exit, error) {
	if data.Config.File.BranchPrefix.IsSome() {
		return None[configdomain.BranchPrefix](), false, nil
	}
	return dialog.BranchPrefix(dialog.Args[configdomain.BranchPrefix]{
		Global: data.Config.GitGlobal.BranchPrefix,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.BranchPrefix,
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

func enterDetached(data Data) (Option[configdomain.Detached], dialogdomain.Exit, error) {
	if data.Config.File.Detached.IsSome() {
		return None[configdomain.Detached](), false, nil
	}
	return dialog.SyncDetached(dialog.Args[configdomain.Detached]{
		Global: data.Config.GitGlobal.Detached,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.Detached,
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

func enterForgejoToken(data Data) (Option[forgedomain.ForgejoToken], dialogdomain.Exit, error) {
	if data.Config.File.ForgejoToken.IsSome() {
		return None[forgedomain.ForgejoToken](), false, nil
	}
	return dialog.ForgejoToken(dialog.Args[forgedomain.ForgejoToken]{
		Global: data.Config.GitGlobal.ForgejoToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ForgejoToken,
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

func enterGithubConnectorType(data Data) (Option[forgedomain.GithubConnectorType], dialogdomain.Exit, error) {
	if data.Config.File.GithubConnectorType.IsSome() {
		return None[forgedomain.GithubConnectorType](), false, nil
	}
	return dialog.GithubConnectorType(dialog.Args[forgedomain.GithubConnectorType]{
		Global: data.Config.GitGlobal.GithubConnectorType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GithubConnectorType,
	})
}

func enterGithubToken(data Data) (Option[forgedomain.GithubToken], dialogdomain.Exit, error) {
	if data.Config.File.GithubToken.IsSome() {
		return None[forgedomain.GithubToken](), false, nil
	}
	return dialog.GithubToken(dialog.Args[forgedomain.GithubToken]{
		Global: data.Config.GitGlobal.GithubToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GithubToken,
	})
}

func enterGitlabConnectorType(data Data) (Option[forgedomain.GitlabConnectorType], dialogdomain.Exit, error) {
	if data.Config.File.GitlabConnectorType.IsSome() {
		return None[forgedomain.GitlabConnectorType](), false, nil
	}
	return dialog.GitlabConnectorType(dialog.Args[forgedomain.GitlabConnectorType]{
		Global: data.Config.GitGlobal.GitlabConnectorType,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GitlabConnectorType,
	})
}

func enterGitlabToken(data Data) (Option[forgedomain.GitlabToken], dialogdomain.Exit, error) {
	if data.Config.File.GitlabToken.IsSome() {
		return None[forgedomain.GitlabToken](), false, nil
	}
	return dialog.GitlabToken(dialog.Args[forgedomain.GitlabToken]{
		Global: data.Config.GitGlobal.GitlabToken,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.GitlabToken,
	})
}

func enterIgnoreUncommitted(data Data) (Option[configdomain.IgnoreUncommitted], dialogdomain.Exit, error) {
	if data.Config.File.IgnoreUncommitted.IsSome() {
		return None[configdomain.IgnoreUncommitted](), false, nil
	}
	return dialog.IgnoreUncommitted(dialog.Args[configdomain.IgnoreUncommitted]{
		Global: data.Config.GitGlobal.IgnoreUncommitted,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.IgnoreUncommitted,
	})
}

func enterMainBranch(data Data) (dialog.MainBranchResult, dialogdomain.Exit, error) {
	if configFileMainBranch, hasMain := data.Config.File.MainBranch.Get(); hasMain {
		return dialog.MainBranchResult{
			UserChoice:       Some(configFileMainBranch),
			ActualMainBranch: configFileMainBranch,
		}, false, nil
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

func enterOrder(data Data) (Option[configdomain.Order], dialogdomain.Exit, error) {
	if data.Config.File.Order.IsSome() {
		return None[configdomain.Order](), false, nil
	}
	return dialog.Order(dialog.Args[configdomain.Order]{
		Global: data.Config.GitGlobal.Order,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.Order,
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
		AppendAllMissing(data.Config.File.PerennialBranches).
		AppendAllMissing(data.Config.GitGlobal.PerennialBranches)
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

func enterProposalBreadcrumb(data Data) (Option[configdomain.ProposalBreadcrumb], dialogdomain.Exit, error) {
	if data.Config.File.ProposalBreadcrumb.IsSome() {
		return None[configdomain.ProposalBreadcrumb](), false, nil
	}
	return dialog.ProposalBreadcrumb(dialog.Args[configdomain.ProposalBreadcrumb]{
		Global: data.Config.GitGlobal.ProposalBreadcrumb,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.ProposalBreadcrumb,
	})
}

func enterPushBranches(data Data) (Option[configdomain.PushBranches], dialogdomain.Exit, error) {
	if data.Config.File.PushBranches.IsSome() {
		return None[configdomain.PushBranches](), false, nil
	}
	return dialog.PushBranches(dialog.Args[configdomain.PushBranches]{
		Global: data.Config.GitGlobal.PushBranches,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.PushBranches,
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

func enterStash(data Data) (Option[configdomain.Stash], dialogdomain.Exit, error) {
	if data.Config.File.Stash.IsSome() {
		return None[configdomain.Stash](), false, nil
	}
	return dialog.Stash(dialog.Args[configdomain.Stash]{
		Global: data.Config.GitGlobal.Stash,
		Inputs: data.Inputs,
		Local:  data.Config.GitLocal.Stash,
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
	data                 Data
	determinedForgeType  Option[forgedomain.ForgeType]
	existingConfig       config.NormalConfig
	forgejoToken         Option[forgedomain.ForgejoToken]
	giteaToken           Option[forgedomain.GiteaToken]
	githubToken          Option[forgedomain.GithubToken]
	gitlabToken          Option[forgedomain.GitlabToken]
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

func existsAndChanged[T any](input, existing Option[T]) bool {
	return input.IsSome() && !input.Equal(existing)
}

func shouldAskForScope(args enterTokenScopeArgs) bool {
	if forgeType, hasForgeType := args.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeAzuredevops:
			return false
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			return existsAndChanged(args.bitbucketUsername, args.existingConfig.BitbucketUsername) &&
				existsAndChanged(args.bitbucketAppPassword, args.existingConfig.BitbucketAppPassword)
		case forgedomain.ForgeTypeForgejo:
			return existsAndChanged(args.forgejoToken, args.existingConfig.ForgejoToken)
		case forgedomain.ForgeTypeGitea:
			return existsAndChanged(args.giteaToken, args.existingConfig.GiteaToken)
		case forgedomain.ForgeTypeGithub:
			return existsAndChanged(args.githubToken, args.existingConfig.GithubToken)
		case forgedomain.ForgeTypeGitlab:
			return existsAndChanged(args.gitlabToken, args.existingConfig.GitlabToken)
		}
	}
	return false
}

func testForgeAuth(args testForgeAuthArgs) (configdomain.ProgramFlow, dialogdomain.Exit, error) {
	if _, inTest := os.LookupEnv(subshell.TestToken); inTest {
		return configdomain.ProgramFlowContinue, false, nil
	}
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              args.backend,
		BitbucketAppPassword: args.bitbucketAppPassword,
		BitbucketUsername:    args.bitbucketUsername,
		Browser:              None[configdomain.Browser](),
		ConfigDir:            args.configDir,
		ForgeType:            args.forgeTypeOpt,
		ForgejoToken:         args.forgejoToken,
		Frontend:             args.backend,
		GiteaToken:           args.giteaToken,
		GithubConnectorType:  args.githubConnectorType,
		GithubToken:          args.githubToken,
		GitlabConnectorType:  args.gitlabConnectorType,
		GitlabToken:          args.gitlabToken,
		Log:                  print.Logger{},
		RemoteURL:            args.devURL,
	})
	if err != nil {
		return configdomain.ProgramFlowExit, false, err
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return configdomain.ProgramFlowContinue, false, nil
	}
	if credentialsVerifier, canVerifyCredentials := connector.(forgedomain.CredentialVerifier); canVerifyCredentials {
		verifyResult := credentialsVerifier.VerifyCredentials()
		if verifyResult.AuthenticationError != nil {
			return dialog.CredentialsNoAccess(verifyResult.AuthenticationError, args.inputs)
		}
		if user, hasUser := verifyResult.AuthenticatedUser.Get(); hasUser {
			fmt.Printf(messages.CredentialsForgeUserName, dialogcomponents.FormattedSelection(user, false))
		}
		if verifyResult.AuthorizationError != nil {
			return dialog.CredentialsNoProposalAccess(verifyResult.AuthorizationError, args.inputs)
		}
		fmt.Println(messages.CredentialsAccess)
	}
	return configdomain.ProgramFlowContinue, false, nil
}

type testForgeAuthArgs struct {
	backend              subshelldomain.RunnerQuerier
	bitbucketAppPassword Option[forgedomain.BitbucketAppPassword]
	bitbucketUsername    Option[forgedomain.BitbucketUsername]
	configDir            configdomain.RepoConfigDir
	devURL               Option[giturl.Parts]
	forgeTypeOpt         Option[forgedomain.ForgeType]
	forgejoToken         Option[forgedomain.ForgejoToken]
	giteaToken           Option[forgedomain.GiteaToken]
	githubConnectorType  Option[forgedomain.GithubConnectorType]
	githubToken          Option[forgedomain.GithubToken]
	gitlabConnectorType  Option[forgedomain.GitlabConnectorType]
	gitlabToken          Option[forgedomain.GitlabToken]
	inputs               dialogcomponents.Inputs
	remoteURL            Option[giturl.Parts]
}

func tokenScopeDialog(args enterTokenScopeArgs) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if forgeType, hasForgeType := args.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeAzuredevops:
			return configdomain.ConfigScopeLocal, false, nil
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyBitbucketUsername, args.data.Config.NormalConfig.BitbucketUsername)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeForgejo:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyForgejoToken, args.data.Config.NormalConfig.ForgejoToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitea:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyGiteaToken, args.data.Config.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGithub:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyGithubToken, args.data.Config.NormalConfig.GithubToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitlab:
			existingScope := determineExistingScope(args.data.Snapshot, configdomain.KeyGitlabToken, args.data.Config.NormalConfig.GitlabToken)
			return dialog.TokenScope(existingScope, args.inputs)
		}
	}
	return configdomain.ConfigScopeLocal, false, nil
}
