package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
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
			cliConfig := cliconfig.CliConfig{
				DryRun:  false,
				Verbose: verbose,
			}
			return executeConfigSetup(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigSetup(cliConfig cliconfig.CliConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := loadSetupData(repo, cliConfig)
	if err != nil || exit {
		return err
	}
	enterDataResult, exit, err := enterData(repo, data)
	if err != nil || exit {
		return err
	}
	if err = saveAll(enterDataResult, repo.UnvalidatedConfig, repo.UnvalidatedConfig.File, data, repo.Frontend); err != nil {
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
		Verbose:               cliConfig.Verbose,
	})
}

type setupData struct {
	backend       subshelldomain.Querier
	dialogInputs  dialogcomponents.TestInputs
	localBranches gitdomain.BranchInfos
	remotes       gitdomain.Remotes
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

func enterData(repo execute.OpenRepoResult, data setupData) (userInput, dialogdomain.Exit, error) {
	var emptyResult userInput
	exit, err := dialog.Welcome(data.dialogInputs)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	aliases, exit, err := dialog.Aliases(configdomain.AllAliasableCommands(), repo.UnvalidatedConfig.NormalConfig.Aliases, data.dialogInputs)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	mainBranchSetting, actualMainBranch, exit, err := enterMainBranch(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	perennialBranches, exit, err := enterPerennialBranches(repo, data, actualMainBranch)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	perennialRegex, exit, err := enterPerennialRegex(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	featureRegex, exit, err := enterFeatureRegex(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	contributionRegex, exit, err := enterContributionRegex(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	observedRegex, exit, err := enterObservedRegex(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	newBranchType, exit, err := enterNewBranchType(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	unknownBranchType, exit, err := enterUnknownBranchType(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	devRemote, exit, err := enterDevRemote(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
EnterForgeData:
	hostingOriginHostName, exit, err := enterOriginHostName(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	enteredForgeType, exit, err := enterForgeType(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	devURL := repo.UnvalidatedConfig.NormalConfig.DevURL(data.backend)
	actualForgeType := determineForgeType(enteredForgeType.Or(repo.UnvalidatedConfig.File.ForgeType), devURL)

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
			bitbucketUsername, exit, err = enterBitbucketUserName(repo, data)
			if err != nil || exit {
				return emptyResult, exit, err
			}
			bitbucketAppPassword, exit, err = enterBitbucketAppPassword(repo, data)
		case forgedomain.ForgeTypeCodeberg:
			codebergToken, exit, err = enterCodebergToken(repo, data)
		case forgedomain.ForgeTypeGitea:
			giteaToken, exit, err = enterGiteaToken(repo, data)
		case forgedomain.ForgeTypeGitHub:
			githubConnectorTypeOpt, exit, err = enterGitHubConnectorType(repo, data)
			if err != nil || exit {
				return emptyResult, exit, err
			}
			if githubConnectorType, has := githubConnectorTypeOpt.Get(); has {
				switch githubConnectorType {
				case forgedomain.GitHubConnectorTypeAPI:
					githubToken, exit, err = enterGitHubToken(repo, data)
				case forgedomain.GitHubConnectorTypeGh:
				}
			}
		case forgedomain.ForgeTypeGitLab:
			gitlabConnectorTypeOpt, exit, err = enterGitLabConnectorType(repo, data)
			if err != nil || exit {
				return emptyResult, exit, err
			}
			if gitlabConnectorType, has := gitlabConnectorTypeOpt.Get(); has {
				switch gitlabConnectorType {
				case forgedomain.GitLabConnectorTypeAPI:
					gitlabToken, exit, err = enterGitLabToken(repo, data)
				case forgedomain.GitLabConnectorTypeGlab:
				}
			}
		}
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	repeat, exit, err := testForgeAuth(testForgeAuthArgs{
		backend:              repo.Backend,
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
		inputs:               data.dialogInputs,
		remoteURL:            repo.UnvalidatedConfig.NormalConfig.RemoteURL(data.backend, devRemote.GetOrElse(config.DefaultNormalConfig().DevRemote)),
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
		determinedForgeType:  actualForgeType,
		existingConfig:       repo.UnvalidatedConfig.NormalConfig,
		giteaToken:           giteaToken,
		githubToken:          githubToken,
		gitlabToken:          gitlabToken,
		inputs:               data.dialogInputs,
		repo:                 repo,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncFeatureStrategy, exit, err := enterSyncFeatureStrategy(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncPerennialStrategy, exit, err := enterSyncPerennialStrategy(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncPrototypeStrategy, exit, err := enterSyncPrototypeStrategy(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncUpstream, exit, err := enterSyncUpstream(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncTags, exit, err := enterSyncTags(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	shareNewBranches, exit, err := enterShareNewBranches(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	pushHook, exit, err := enterPushHook(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	shipStrategy, exit, err := enterShipStrategy(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	shipDeleteTrackingBranch, exit, err := enterShipDeleteTrackingBranch(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	configStorage, exit, err := dialog.ConfigStorage(data.dialogInputs)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	normalData := configdomain.PartialConfig{
		Aliases:                  aliases,
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
	if !data.dialogInputs.IsEmpty() {
		panic("unused dialog inputs")
	}
	return userInput{normalData, actualForgeType, tokenScope, configStorage, validatedData}, false, nil
}

// data entered by the user in the setup assistant
type userInput struct {
	data                configdomain.PartialConfig
	determinedForgeType Option[forgedomain.ForgeType] // the forge type that was determined by the setup assistant - not necessarily what the user entered (could also be "auto detect")
	scope               configdomain.ConfigScope
	storageLocation     dialog.ConfigStorageOption
	validatedConfig     configdomain.ValidatedConfigData
}

func enterBitbucketUserName(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.BitbucketUsername], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.BitbucketUsername.IsSome() {
		return None[forgedomain.BitbucketUsername](), false, nil
	}
	return dialog.BitbucketUsername(dialog.Args[forgedomain.BitbucketUsername]{
		Global: repo.UnvalidatedConfig.GitLocal.BitbucketUsername,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.BitbucketUsername,
	})
}

func enterBitbucketAppPassword(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.BitbucketAppPassword], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.BitbucketUsername.IsSome() {
		return None[forgedomain.BitbucketAppPassword](), false, nil
	}
	return dialog.BitbucketAppPassword(dialog.Args[forgedomain.BitbucketAppPassword]{
		Global: repo.UnvalidatedConfig.GitLocal.BitbucketAppPassword,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.BitbucketAppPassword,
	})
}

func enterCodebergToken(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.CodebergToken], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.CodebergToken.IsSome() {
		return None[forgedomain.CodebergToken](), false, nil
	}
	return dialog.CodebergToken(dialog.Args[forgedomain.CodebergToken]{
		Global: repo.UnvalidatedConfig.GitGlobal.CodebergToken,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.CodebergToken,
	})
}

func enterGiteaToken(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.GiteaToken], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.GiteaToken.IsSome() {
		return None[forgedomain.GiteaToken](), false, nil
	}
	return dialog.GiteaToken(dialog.Args[forgedomain.GiteaToken]{
		Global: repo.UnvalidatedConfig.GitGlobal.GiteaToken,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.GiteaToken,
	})
}

func enterGitHubConnectorType(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.GitHubConnectorType], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.GitHubConnectorType.IsSome() {
		return None[forgedomain.GitHubConnectorType](), false, nil
	}
	return dialog.GitHubConnectorType(dialog.Args[forgedomain.GitHubConnectorType]{
		Global: repo.UnvalidatedConfig.GitGlobal.GitHubConnectorType,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.GitHubConnectorType,
	})
}

func enterGitHubToken(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.GitHubToken], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.GitHubToken.IsSome() {
		return None[forgedomain.GitHubToken](), false, nil
	}
	return dialog.GitHubToken(dialog.Args[forgedomain.GitHubToken]{
		Global: repo.UnvalidatedConfig.GitGlobal.GitHubToken,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.GitHubToken,
	})
}

func enterGitLabConnectorType(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.GitLabConnectorType], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.GitLabConnectorType.IsSome() {
		return None[forgedomain.GitLabConnectorType](), false, nil
	}
	return dialog.GitLabConnectorType(dialog.Args[forgedomain.GitLabConnectorType]{
		Global: repo.UnvalidatedConfig.GitGlobal.GitLabConnectorType,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.GitLabConnectorType,
	})
}

func enterGitLabToken(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.GitLabToken], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.GitLabToken.IsSome() {
		return None[forgedomain.GitLabToken](), false, nil
	}
	return dialog.GitLabToken(dialog.Args[forgedomain.GitLabToken]{
		Global: repo.UnvalidatedConfig.GitGlobal.GitLabToken,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.GitLabToken,
	})
}

func enterContributionRegex(repo execute.OpenRepoResult, data setupData) (Option[configdomain.ContributionRegex], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.ContributionRegex.IsSome() {
		return None[configdomain.ContributionRegex](), false, nil
	}
	return dialog.ContributionRegex(dialog.Args[configdomain.ContributionRegex]{
		Global: repo.UnvalidatedConfig.GitGlobal.ContributionRegex,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.ContributionRegex,
	})
}

func enterDevRemote(repo execute.OpenRepoResult, data setupData) (Option[gitdomain.Remote], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.DevRemote.IsSome() {
		return None[gitdomain.Remote](), false, nil
	}
	return dialog.DevRemote(data.remotes, dialog.Args[gitdomain.Remote]{
		Global: repo.UnvalidatedConfig.GitGlobal.DevRemote,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.DevRemote,
	})
}

func enterFeatureRegex(repo execute.OpenRepoResult, data setupData) (Option[configdomain.FeatureRegex], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.FeatureRegex.IsSome() {
		return None[configdomain.FeatureRegex](), false, nil
	}
	return dialog.FeatureRegex(dialog.Args[configdomain.FeatureRegex]{
		Global: repo.UnvalidatedConfig.GitGlobal.FeatureRegex,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.FeatureRegex,
	})
}

func enterForgeType(repo execute.OpenRepoResult, data setupData) (Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.ForgeType.IsSome() {
		return None[forgedomain.ForgeType](), false, nil
	}
	return dialog.ForgeType(dialog.Args[forgedomain.ForgeType]{
		Global: repo.UnvalidatedConfig.GitGlobal.ForgeType,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.ForgeType,
	})
}

func enterMainBranch(repo execute.OpenRepoResult, data setupData) (userChoice Option[gitdomain.LocalBranchName], actualMainBranch gitdomain.LocalBranchName, exit dialogdomain.Exit, err error) {
	if configFileMainBranch, hasMain := repo.UnvalidatedConfig.File.MainBranch.Get(); hasMain {
		return Some(configFileMainBranch), configFileMainBranch, false, nil
	}
	return dialog.MainBranch(dialog.MainBranchArgs{
		GitStandardBranch:     repo.Git.StandardBranch(repo.Backend),
		Inputs:                data.dialogInputs,
		LocalBranches:         data.localBranches.Names(),
		LocalGitMainBranch:    repo.UnvalidatedConfig.GitLocal.MainBranch,
		UnscopedGitMainBranch: repo.UnvalidatedConfig.GitUnscoped.MainBranch,
	})
}

func enterNewBranchType(repo execute.OpenRepoResult, data setupData) (Option[configdomain.NewBranchType], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.NewBranchType.IsSome() {
		return None[configdomain.NewBranchType](), false, nil
	}
	return dialog.NewBranchType(dialog.Args[configdomain.NewBranchType]{
		Global: repo.UnvalidatedConfig.GitGlobal.NewBranchType,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.NewBranchType,
	})
}

func enterObservedRegex(repo execute.OpenRepoResult, data setupData) (Option[configdomain.ObservedRegex], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.ObservedRegex.IsSome() {
		return None[configdomain.ObservedRegex](), false, nil
	}
	return dialog.ObservedRegex(dialog.Args[configdomain.ObservedRegex]{
		Global: repo.UnvalidatedConfig.GitGlobal.ObservedRegex,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.ObservedRegex,
	})
}

func enterOriginHostName(repo execute.OpenRepoResult, data setupData) (Option[configdomain.HostingOriginHostname], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.HostingOriginHostname.IsSome() {
		return None[configdomain.HostingOriginHostname](), false, nil
	}
	return dialog.OriginHostname(dialog.Args[configdomain.HostingOriginHostname]{
		Global: repo.UnvalidatedConfig.GitGlobal.HostingOriginHostname,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.HostingOriginHostname,
	})
}

func enterPerennialBranches(repo execute.OpenRepoResult, data setupData, mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchNames, dialogdomain.Exit, error) {
	immutablePerennials := gitdomain.LocalBranchNames{mainBranch}.
		AppendAllMissing(repo.UnvalidatedConfig.File.PerennialBranches...).
		AppendAllMissing(repo.UnvalidatedConfig.GitGlobal.PerennialBranches...)
	return dialog.PerennialBranches(dialog.PerennialBranchesArgs{
		ImmutableGitPerennials: immutablePerennials,
		Inputs:                 data.dialogInputs,
		LocalBranches:          data.localBranches.Names(),
		LocalGitPerennials:     repo.UnvalidatedConfig.GitLocal.PerennialBranches,
		MainBranch:             mainBranch,
	})
}

func enterPerennialRegex(repo execute.OpenRepoResult, data setupData) (Option[configdomain.PerennialRegex], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.PerennialRegex.IsSome() {
		return None[configdomain.PerennialRegex](), false, nil
	}
	return dialog.PerennialRegex(dialog.Args[configdomain.PerennialRegex]{
		Global: repo.UnvalidatedConfig.GitGlobal.PerennialRegex,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.PerennialRegex,
	})
}

func enterSyncFeatureStrategy(repo execute.OpenRepoResult, data setupData) (Option[configdomain.SyncFeatureStrategy], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.SyncFeatureStrategy.IsSome() {
		return None[configdomain.SyncFeatureStrategy](), false, nil
	}
	return dialog.SyncFeatureStrategy(dialog.Args[configdomain.SyncFeatureStrategy]{
		Global: repo.UnvalidatedConfig.GitGlobal.SyncFeatureStrategy,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.SyncFeatureStrategy,
	})
}

func enterSyncPerennialStrategy(repo execute.OpenRepoResult, data setupData) (Option[configdomain.SyncPerennialStrategy], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.SyncPerennialStrategy.IsSome() {
		return None[configdomain.SyncPerennialStrategy](), false, nil
	}
	return dialog.SyncPerennialStrategy(dialog.Args[configdomain.SyncPerennialStrategy]{
		Global: repo.UnvalidatedConfig.GitGlobal.SyncPerennialStrategy,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.SyncPerennialStrategy,
	})
}

func enterSyncPrototypeStrategy(repo execute.OpenRepoResult, data setupData) (Option[configdomain.SyncPrototypeStrategy], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.SyncPrototypeStrategy.IsSome() {
		return None[configdomain.SyncPrototypeStrategy](), false, nil
	}
	return dialog.SyncPrototypeStrategy(dialog.Args[configdomain.SyncPrototypeStrategy]{
		Global: repo.UnvalidatedConfig.GitGlobal.SyncPrototypeStrategy,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.SyncPrototypeStrategy,
	})
}

func enterSyncUpstream(repo execute.OpenRepoResult, data setupData) (Option[configdomain.SyncUpstream], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.SyncUpstream.IsSome() {
		return None[configdomain.SyncUpstream](), false, nil
	}
	return dialog.SyncUpstream(dialog.Args[configdomain.SyncUpstream]{
		Global: repo.UnvalidatedConfig.GitGlobal.SyncUpstream,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.SyncUpstream,
	})
}

func enterSyncTags(repo execute.OpenRepoResult, data setupData) (Option[configdomain.SyncTags], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.SyncTags.IsSome() {
		return None[configdomain.SyncTags](), false, nil
	}
	return dialog.SyncTags(dialog.Args[configdomain.SyncTags]{
		Global: repo.UnvalidatedConfig.GitGlobal.SyncTags,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.SyncTags,
	})
}

func enterShareNewBranches(repo execute.OpenRepoResult, data setupData) (Option[configdomain.ShareNewBranches], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.ShareNewBranches.IsSome() {
		return None[configdomain.ShareNewBranches](), false, nil
	}
	return dialog.ShareNewBranches(dialog.Args[configdomain.ShareNewBranches]{
		Global: repo.UnvalidatedConfig.GitGlobal.ShareNewBranches,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.ShareNewBranches,
	})
}

func enterPushHook(repo execute.OpenRepoResult, data setupData) (Option[configdomain.PushHook], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.PushHook.IsSome() {
		return None[configdomain.PushHook](), false, nil
	}
	return dialog.PushHook(dialog.Args[configdomain.PushHook]{
		Global: repo.UnvalidatedConfig.GitGlobal.PushHook,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.PushHook,
	})
}

func enterShipStrategy(repo execute.OpenRepoResult, data setupData) (Option[configdomain.ShipStrategy], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.ShipStrategy.IsSome() {
		return None[configdomain.ShipStrategy](), false, nil
	}
	return dialog.ShipStrategy(dialog.Args[configdomain.ShipStrategy]{
		Global: repo.UnvalidatedConfig.GitGlobal.ShipStrategy,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.ShipStrategy,
	})
}

func enterShipDeleteTrackingBranch(repo execute.OpenRepoResult, data setupData) (Option[configdomain.ShipDeleteTrackingBranch], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.ShipDeleteTrackingBranch.IsSome() {
		return None[configdomain.ShipDeleteTrackingBranch](), false, nil
	}
	return dialog.ShipDeleteTrackingBranch(dialog.Args[configdomain.ShipDeleteTrackingBranch]{
		Global: repo.UnvalidatedConfig.GitGlobal.ShipDeleteTrackingBranch,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.ShipDeleteTrackingBranch,
	})
}

func enterUnknownBranchType(repo execute.OpenRepoResult, data setupData) (Option[configdomain.UnknownBranchType], dialogdomain.Exit, error) {
	if repo.UnvalidatedConfig.File.UnknownBranchType.IsSome() {
		return None[configdomain.UnknownBranchType](), false, nil
	}
	return dialog.UnknownBranchType(dialog.Args[configdomain.UnknownBranchType]{
		Global: repo.UnvalidatedConfig.GitGlobal.UnknownBranchType,
		Inputs: data.dialogInputs,
		Local:  repo.UnvalidatedConfig.GitLocal.UnknownBranchType,
	})
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
	inputs               dialogcomponents.TestInputs
	remoteURL            Option[giturl.Parts]
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
	determinedForgeType  Option[forgedomain.ForgeType]
	existingConfig       config.NormalConfig
	giteaToken           Option[forgedomain.GiteaToken]
	githubToken          Option[forgedomain.GitHubToken]
	gitlabToken          Option[forgedomain.GitLabToken]
	inputs               dialogcomponents.TestInputs
	repo                 execute.OpenRepoResult
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

func tokenScopeDialog(args enterTokenScopeArgs) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if forgeType, hasForgeType := args.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyBitbucketUsername, args.repo.UnvalidatedConfig.NormalConfig.BitbucketUsername)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeCodeberg:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyCodebergToken, args.repo.UnvalidatedConfig.NormalConfig.CodebergToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitea:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyGiteaToken, args.repo.UnvalidatedConfig.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitHub:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyGitHubToken, args.repo.UnvalidatedConfig.NormalConfig.GitHubToken)
			return dialog.TokenScope(existingScope, args.inputs)
		case forgedomain.ForgeTypeGitLab:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyGitLabToken, args.repo.UnvalidatedConfig.NormalConfig.GitLabToken)
			return dialog.TokenScope(existingScope, args.inputs)
		}
	}
	return configdomain.ConfigScopeLocal, false, nil
}

func determineExistingScope(configSnapshot undoconfig.ConfigSnapshot, key configdomain.Key, oldValue fmt.Stringer) configdomain.ConfigScope {
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

func existsAndChanged[T fmt.Stringer](input, existing T) bool {
	return input.String() != "" && input.String() != existing.String()
}

func loadSetupData(repo execute.OpenRepoResult, cliConfig cliconfig.CliConfig) (data setupData, exit dialogdomain.Exit, err error) {
	dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             None[forgedomain.Connector](),
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
		Verbose:               cliConfig.Verbose,
	})
	if err != nil {
		return data, exit, err
	}
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, exit, err
	}
	if len(remotes) == 0 {
		remotes = gitdomain.Remotes{gitconfig.DefaultRemote(repo.Backend)}
	}
	return setupData{
		backend:       repo.Backend,
		dialogInputs:  dialogTestInputs,
		localBranches: branchesSnapshot.Branches,
		remotes:       remotes,
	}, exit, nil
}

func saveAll(userInput userInput, unvalidatedConfig config.UnvalidatedConfig, configFile configdomain.PartialConfig, data setupData, frontend subshelldomain.Runner) error {
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
		return saveToFile(userInput, unvalidatedConfig.GitLocal, frontend)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, unvalidatedConfig.GitLocal, configFile, data, frontend)
	}
	return nil
}

func saveToGit(userInput userInput, existingGitConfig configdomain.PartialConfig, configFile configdomain.PartialConfig, data setupData, frontend subshelldomain.Runner) error {
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
			saveMainBranch(userInput.validatedConfig.MainBranch, existingGitConfig.MainBranch, frontend),
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
	if len(data.remotes) > 1 && configFile.DevRemote.IsNone() {
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

func saveUnknownBranchType(valueToWriteToGit Option[configdomain.UnknownBranchType], valueAlreadyInGit Option[configdomain.UnknownBranchType], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetUnknownBranchType(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveUnknownBranchType(runner)
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

func saveCodebergToken(valueToWriteToGit Option[forgedomain.CodebergToken], valueAlreadyInGit Option[forgedomain.CodebergToken], scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetCodebergToken(frontend, value, scope)
	}
	return gitconfig.RemoveCodebergToken(frontend)
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

func saveMainBranch(valueToWriteToGit gitdomain.LocalBranchName, valueAlreadyInGit Option[gitdomain.LocalBranchName], runner subshelldomain.Runner) error {
	if existing, hasExisting := valueAlreadyInGit.Get(); hasExisting {
		if existing == valueToWriteToGit {
			return nil
		}
	}
	return gitconfig.SetMainBranch(runner, valueToWriteToGit, configdomain.ConfigScopeLocal)
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

func saveSyncUpstream(valueToWriteToGit Option[configdomain.SyncUpstream], valueAlreadyInGit Option[configdomain.SyncUpstream], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetSyncUpstream(runner, value, configdomain.ConfigScopeLocal)
	}
	return gitconfig.RemoveSyncUpstream(runner)
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

func saveToFile(userInput userInput, gitConfig configdomain.PartialConfig, runner subshelldomain.Runner) error {
	userInput.data.MainBranch = Some(userInput.validatedConfig.MainBranch)
	if err := configfile.Save(userInput.data); err != nil {
		return err
	}
	if gitConfig.DevRemote.IsSome() {
		_ = gitconfig.RemoveDevRemote(runner)
	}
	if gitConfig.MainBranch.IsSome() {
		_ = gitconfig.RemoveMainBranch(runner)
	}
	if gitConfig.NewBranchType.IsSome() {
		_ = gitconfig.RemoveNewBranchType(runner)
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
	if err := saveUnknownBranchType(userInput.data.UnknownBranchType, gitConfig.UnknownBranchType, runner); err != nil {
		return err
	}
	// TODO: also save ObservedRegex ContributionRegex NewBranchType
	return saveFeatureRegex(userInput.data.FeatureRegex, gitConfig.FeatureRegex, runner)
}
