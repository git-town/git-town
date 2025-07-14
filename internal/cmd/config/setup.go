package config

import (
	"cmp"
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
	if err = saveAll(enterDataResult, repo.UnvalidatedConfig.NormalConfig.Git, data.configFile, data, repo.Frontend); err != nil {
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
	config        config.UnvalidatedConfig
	configFile    Option[configdomain.PartialConfig]
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
	configFile := data.configFile.GetOrDefault()
	exit, err := dialog.Welcome(data.dialogInputs.Next())
	if err != nil || exit {
		return emptyResult, exit, err
	}
	aliases, exit, err := dialog.Aliases(configdomain.AllAliasableCommands(), repo.UnvalidatedConfig.NormalConfig.Aliases, data.dialogInputs.Next())
	if err != nil || exit {
		return emptyResult, exit, err
	}
	mainBranchOpt, actualMainBranch, exit, err := enterMainBranch(repo, data)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	perennialBranches, exit, err := enterPerennialBranches(repo, data, actualMainBranch)
	if err != nil || exit {
		return emptyResult, exit, err
	}
	perennialRegex, exit, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.PerennialRegex]{
		ConfigFileValue: configFile.PerennialRegex,
		HelpText:        dialog.PerennialBranchesHelp,
		Inputs:          data.dialogInputs,
		LocalValue:      repo.UnvalidatedConfig.GitLocal.PerennialRegex,
		ParseFunc:       configdomain.ParsePerennialRegex,
		Prompt:          "Perennial Regex: ",
		ResultMessage:   messages.PerennialRegex,
		Title:           dialog.PerennialRegexTitle,
		UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.PerennialRegex,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	featureRegex, exit, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.FeatureRegex]{
		ConfigFileValue: configFile.FeatureRegex,
		HelpText:        dialog.FeatureRegexHelp,
		Inputs:          data.dialogInputs,
		LocalValue:      repo.UnvalidatedConfig.GitLocal.FeatureRegex,
		ParseFunc:       configdomain.ParseFeatureRegex,
		Prompt:          "Feature Regex: ",
		ResultMessage:   messages.FeatureRegex,
		Title:           dialog.FeatureRegexTitle,
		UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.FeatureRegex,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	contributionRegex, exit, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.ContributionRegex]{
		ConfigFileValue: configFile.ContributionRegex,
		HelpText:        dialog.ContributionRegexHelp,
		Inputs:          data.dialogInputs,
		LocalValue:      repo.UnvalidatedConfig.GitLocal.ContributionRegex,
		ParseFunc:       configdomain.ParseContributionRegex,
		Prompt:          "Contribution Regex: ",
		ResultMessage:   messages.ContributionRegex,
		Title:           dialog.ContributionRegexTitle,
		UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.ContributionRegex,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	observedRegex, exit, err := dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.ObservedRegex]{
		ConfigFileValue: configFile.ObservedRegex,
		HelpText:        dialog.ObservedRegexHelp,
		Inputs:          data.dialogInputs,
		LocalValue:      repo.UnvalidatedConfig.GitLocal.ObservedRegex,
		ParseFunc:       configdomain.ParseObservedRegex,
		Prompt:          "Observed Regex: ",
		ResultMessage:   messages.ObservedRegex,
		Title:           dialog.ObservedRegexTitle,
		UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.ObservedRegex,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	unknownBranchType := repo.UnvalidatedConfig.NormalConfig.UnknownBranchType
	if configFile.UnknownBranchType.IsNone() {
		unknownBranchType, exit, err = dialog.UnknownBranchType(unknownBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	devRemote := None[gitdomain.Remote]()
	if configFile.DevRemote.IsNone() && len(data.remotes) > 1 {
		devRemote, exit, err = dialog.DevRemote(repo.UnvalidatedConfig.NormalConfig.DevRemote, data.remotes, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	hostingOriginHostName := None[configdomain.HostingOriginHostname]()
	enteredForgeType := repo.UnvalidatedConfig.NormalConfig.ForgeType.Or(repo.UnvalidatedConfig.File.GetOrDefault().ForgeType)
	var actualForgeType Option[forgedomain.ForgeType]
	bitbucketUsername := None[forgedomain.BitbucketUsername]()
	bitbucketAppPassword := None[forgedomain.BitbucketAppPassword]()
	codebergToken := None[forgedomain.CodebergToken]()
	devURL := data.config.NormalConfig.DevURL(data.backend)
	giteaToken := None[forgedomain.GiteaToken]()
	githubConnectorTypeOpt := None[forgedomain.GitHubConnectorType]()
	githubToken := None[forgedomain.GitHubToken]()
	gitlabConnectorTypeOpt := None[forgedomain.GitLabConnectorType]()
	gitlabToken := None[forgedomain.GitLabToken]()
	for {
		hostingOriginHostName, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[configdomain.HostingOriginHostname]{
			ConfigFileValue: configFile.HostingOriginHostname,
			HelpText:        dialog.OriginHostnameHelp,
			Inputs:          data.dialogInputs,
			LocalValue:      repo.UnvalidatedConfig.GitLocal.HostingOriginHostname,
			ParseFunc:       wrapParseFunc(configdomain.ParseHostingOriginHostname),
			Prompt:          "Origin hostname: ",
			ResultMessage:   messages.OriginHostname,
			Title:           dialog.OriginHostnameTitle,
			UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.HostingOriginHostname,
		})
		if err != nil || exit {
			return emptyResult, exit, err
		}
		if configFile.ForgeType.IsNone() {
			enteredForgeType, exit, err = dialog.ForgeType(enteredForgeType, data.dialogInputs.Next())
			if err != nil || exit {
				return emptyResult, exit, err
			}
		}
		actualForgeType = determineForgeType(enteredForgeType, devURL)
		if forgeType, hasForgeType := actualForgeType.Get(); hasForgeType {
			switch forgeType {
			case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
				bitbucketUsername, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.BitbucketUsername]{
					ConfigFileValue: configFile.BitbucketUsername,
					HelpText:        dialog.BitbucketUsernameHelp,
					Inputs:          data.dialogInputs,
					LocalValue:      repo.UnvalidatedConfig.GitLocal.BitbucketUsername,
					ParseFunc:       wrapParseFunc(forgedomain.ParseBitbucketUsername),
					Prompt:          "Your Bitbucket Username: ",
					ResultMessage:   messages.DialogResultBitbucketUsername,
					Title:           dialog.BitbucketUsernameTitle,
					UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.BitbucketUsername,
				})
				if err != nil || exit {
					return emptyResult, exit, err
				}
				bitbucketAppPassword, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.BitbucketAppPassword]{
					ConfigFileValue: configFile.BitbucketAppPassword,
					HelpText:        dialog.BitbucketAppPasswordHelp,
					Inputs:          data.dialogInputs,
					LocalValue:      repo.UnvalidatedConfig.GitLocal.BitbucketAppPassword,
					ParseFunc:       wrapParseFunc(forgedomain.ParseBitbucketAppPassword),
					Prompt:          "Your Bitbucket App Password: ",
					ResultMessage:   messages.DialogResultBitbucketAppPassword,
					Title:           dialog.BitbucketAppPasswordTitle,
					UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.BitbucketAppPassword,
				})
			case forgedomain.ForgeTypeCodeberg:
				codebergToken, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.CodebergToken]{
					ConfigFileValue: configFile.CodebergToken,
					HelpText:        dialog.CodebergTokenHelp,
					Inputs:          data.dialogInputs,
					LocalValue:      repo.UnvalidatedConfig.GitLocal.CodebergToken,
					ParseFunc:       wrapParseFunc(forgedomain.ParseCodebergToken),
					Prompt:          "Your Codeberg API token: ",
					ResultMessage:   messages.DialogResultCodebergToken,
					Title:           dialog.CodebergTokenTitle,
					UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.CodebergToken,
				})
			case forgedomain.ForgeTypeGitea:
				giteaToken, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.GiteaToken]{
					ConfigFileValue: configFile.GiteaToken,
					HelpText:        dialog.GiteaTokenHelp,
					Inputs:          data.dialogInputs,
					LocalValue:      repo.UnvalidatedConfig.GitLocal.GiteaToken,
					ParseFunc:       wrapParseFunc(forgedomain.ParseGiteaToken),
					Prompt:          "Your Gitea API token: ",
					ResultMessage:   messages.DialogResultGiteaToken,
					Title:           dialog.GiteaTokenTitle,
					UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.GiteaToken,
				})
			case forgedomain.ForgeTypeGitHub:
				githubConnectorTypeOpt, exit, err = dialog.GitHubConnectorType(repo.UnvalidatedConfig.NormalConfig.GitHubConnectorType, data.dialogInputs.Next())
				if err != nil || exit {
					return emptyResult, exit, err
				}
				if githubConnectorType, has := githubConnectorTypeOpt.Get(); has {
					switch githubConnectorType {
					case forgedomain.GitHubConnectorTypeAPI:
						githubToken, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.GitHubToken]{
							ConfigFileValue: configFile.GitHubToken,
							HelpText:        dialog.GitHubTokenHelp,
							Inputs:          data.dialogInputs,
							LocalValue:      repo.UnvalidatedConfig.GitLocal.GitHubToken,
							ParseFunc:       wrapParseFunc(forgedomain.ParseGitHubToken),
							Prompt:          "Your GitHub API token: ",
							ResultMessage:   messages.DialogResultGiteaToken,
							Title:           dialog.GitHubTokenTitle,
							UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.GitHubToken,
						})
					case forgedomain.GitHubConnectorTypeGh:
					}
				}
			case forgedomain.ForgeTypeGitLab:
				gitlabConnectorTypeOpt, exit, err = dialog.GitLabConnectorType(repo.UnvalidatedConfig.NormalConfig.GitLabConnectorType, data.dialogInputs.Next())
				if err != nil || exit {
					return emptyResult, exit, err
				}
				if gitlabConnectorType, has := gitlabConnectorTypeOpt.Get(); has {
					switch gitlabConnectorType {
					case forgedomain.GitLabConnectorTypeAPI:
						gitlabToken, exit, err = dialog.ConfigStringDialog(dialog.ConfigStringDialogArgs[forgedomain.GitLabToken]{
							ConfigFileValue: configFile.GitLabToken,
							HelpText:        dialog.GitLabTokenHelp,
							Inputs:          data.dialogInputs,
							LocalValue:      repo.UnvalidatedConfig.GitLocal.GitLabToken,
							ParseFunc:       wrapParseFunc(forgedomain.ParseGitLabToken),
							Prompt:          "Your GitLab API token: ",
							ResultMessage:   messages.DialogResultGiteaToken,
							Title:           dialog.GitLabTokenTitle,
							UnscopedValue:   repo.UnvalidatedConfig.NormalConfig.Git.GitLabToken,
						})
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
			remoteURL:            data.config.NormalConfig.RemoteURL(data.backend, devRemote.GetOrElse(data.remotes[0])),
		})
		if err != nil || exit {
			return emptyResult, exit, err
		}
		if !repeat {
			break
		}
	}
	tokenScope, exit, err := enterTokenScope(enterTokenScopeArgs{
		bitbucketAppPassword: bitbucketAppPassword,
		bitbucketUsername:    bitbucketUsername,
		codebergToken:        codebergToken,
		determinedForgeType:  actualForgeType,
		existingConfig:       data.config.NormalConfig.NormalConfigData,
		giteaToken:           giteaToken,
		githubToken:          githubToken,
		gitlabToken:          gitlabToken,
		inputs:               data.dialogInputs,
		repo:                 repo,
	})
	if err != nil || exit {
		return emptyResult, exit, err
	}
	syncFeatureStrategy := repo.UnvalidatedConfig.NormalConfig.SyncFeatureStrategy
	if configFile.SyncFeatureStrategy.IsNone() {
		syncFeatureStrategy, exit, err = dialog.SyncFeatureStrategy(syncFeatureStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	syncPerennialStrategy := repo.UnvalidatedConfig.NormalConfig.SyncPerennialStrategy
	if configFile.SyncPerennialStrategy.IsNone() {
		syncPerennialStrategy, exit, err = dialog.SyncPerennialStrategy(syncPerennialStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	syncPrototypeStrategy := None[configdomain.SyncPrototypeStrategy]()
	if configFile.SyncPrototypeStrategy.IsNone() {
		syncPrototypeStrategy, exit, err = dialog.SyncPrototypeStrategy(repo.UnvalidatedConfig.NormalConfig.SyncPrototypeStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	syncUpstream := repo.UnvalidatedConfig.NormalConfig.SyncUpstream
	if configFile.SyncUpstream.IsNone() {
		syncUpstream, exit, err = dialog.SyncUpstream(syncUpstream, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	syncTags := repo.UnvalidatedConfig.NormalConfig.SyncTags
	if configFile.SyncTags.IsNone() {
		syncTags, exit, err = dialog.SyncTags(syncTags, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	shareNewBranches := repo.UnvalidatedConfig.NormalConfig.ShareNewBranches
	if configFile.ShareNewBranches.IsNone() {
		shareNewBranches, exit, err = dialog.ShareNewBranches(shareNewBranches, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	pushHook := repo.UnvalidatedConfig.NormalConfig.PushHook
	if configFile.PushHook.IsNone() {
		pushHook, exit, err = dialog.PushHook(pushHook, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	newBranchType := repo.UnvalidatedConfig.NormalConfig.NewBranchType
	if configFile.NewBranchType.IsNone() {
		newBranchType, exit, err = dialog.NewBranchType(newBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	shipStrategy := repo.UnvalidatedConfig.NormalConfig.ShipStrategy
	if configFile.ShipStrategy.IsNone() {
		shipStrategy, exit, err = dialog.ShipStrategy(shipStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	shipDeleteTrackingBranch := repo.UnvalidatedConfig.NormalConfig.ShipDeleteTrackingBranch
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		shipDeleteTrackingBranch, exit, err = dialog.ShipDeleteTrackingBranch(shipDeleteTrackingBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	configStorage, exit, err := dialog.ConfigStorage(data.dialogInputs.Next())
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
		GiteaToken:               giteaToken,
		HostingOriginHostname:    hostingOriginHostName,
		Lineage:                  configdomain.Lineage{}, // the setup assistant doesn't ask for this
		MainBranch:               mainBranchOpt,
		NewBranchType:            newBranchType,
		ObservedRegex:            observedRegex,
		Offline:                  None[configdomain.Offline](), // the setup assistant doesn't ask for this
		PerennialBranches:        perennialBranches,
		PerennialRegex:           perennialRegex,
		PushHook:                 Some(pushHook),
		ShareNewBranches:         Some(shareNewBranches),
		ShipDeleteTrackingBranch: Some(shipDeleteTrackingBranch),
		ShipStrategy:             Some(shipStrategy),
		SyncFeatureStrategy:      Some(syncFeatureStrategy),
		SyncPerennialStrategy:    Some(syncPerennialStrategy),
		SyncPrototypeStrategy:    syncPrototypeStrategy,
		SyncTags:                 Some(syncTags),
		SyncUpstream:             Some(syncUpstream),
		UnknownBranchType:        Some(unknownBranchType),
		Verbose:                  None[configdomain.Verbose](), // the setup assistant doesn't ask for this
	}
	return userInput{actualForgeType, normalData, tokenScope, configStorage}, false, nil
}

// data entered by the user in the setup assistant
type userInput struct {
	determinedForgeType Option[forgedomain.ForgeType] // the forge type that was determined by the setup assistant - not necessarily what the user entered (could also be "auto detect")
	data                configdomain.PartialConfig
	scope               configdomain.ConfigScope
	storageLocation     dialog.ConfigStorageOption
}

func enterMainBranch(repo execute.OpenRepoResult, data setupData) (userInput Option[gitdomain.LocalBranchName], actualMainBranch gitdomain.LocalBranchName, exit dialogdomain.Exit, err error) {
	if configFile, hasConfigFile := repo.UnvalidatedConfig.File.Get(); hasConfigFile {
		if configFileMainBranch, hasMain := configFile.MainBranch.Get(); hasMain {
			return Some(configFileMainBranch), configFileMainBranch, false, nil
		}
	}
	repoDefault := determineGitRepoDefaultBranch(repo)
	userInput, exit, err = dialog.MainBranch(dialog.MainBranchArgs{
		GitStandardBranch:   repoDefault,
		LocalBranches:       data.localBranches.Names(),
		LocalGitMainBranch:  data.config.GitGlobal.MainBranch,
		GlobalGitMainBranch: data.config.GitLocal.MainBranch,
		Inputs:              data.dialogInputs.Next(),
	})
	if err != nil || exit {
		return None[gitdomain.LocalBranchName](), "", exit, err
	}
	actualMainBranch = userInput.Or(data.config.GitGlobal.MainBranch).GetOrPanic()
	return userInput, actualMainBranch, false, nil
}

func enterPerennialBranches(repo execute.OpenRepoResult, data setupData, mainBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchNames, dialogdomain.Exit, error) {
	if configFile, hasConfigFile := repo.UnvalidatedConfig.File.Get(); hasConfigFile {
		if len(configFile.PerennialBranches) > 0 {
			return gitdomain.LocalBranchNames{}, false, nil
		}
	}
	return dialog.PerennialBranches(dialog.PerennialBranchesArgs{
		LocalBranches:       data.localBranches.Names(),
		MainBranch:          mainBranch,
		GlobalGitPerennials: repo.UnvalidatedConfig.GitGlobal.PerennialBranches,
		LocalGitPerennials:  repo.UnvalidatedConfig.GitLocal.PerennialBranches,
		Inputs:              data.dialogInputs.Next(),
	})
}

// determines the branch that is configured in Git as the default branch
func determineGitRepoDefaultBranch(repo execute.OpenRepoResult) Option[gitdomain.LocalBranchName] {
	if defaultBranch, has := gitconfig.DefaultBranch(repo.Backend).Get(); has {
		return Some(defaultBranch)
	}
	return repo.Git.OriginHead(repo.Backend)
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
		return dialog.CredentialsNoAccess(verifyResult.AuthenticationError, args.inputs.Next())
	}
	if user, hasUser := verifyResult.AuthenticatedUser.Get(); hasUser {
		fmt.Printf(messages.CredentialsForgeUserName, dialogcomponents.FormattedSelection(user, exit))
	}
	if verifyResult.AuthorizationError != nil {
		return dialog.CredentialsNoProposalAccess(verifyResult.AuthorizationError, args.inputs.Next())
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
	existingConfig       configdomain.NormalConfigData
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
			return dialog.TokenScope(existingScope, args.inputs.Next())
		case forgedomain.ForgeTypeCodeberg:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyCodebergToken, args.repo.UnvalidatedConfig.NormalConfig.CodebergToken)
			return dialog.TokenScope(existingScope, args.inputs.Next())
		case forgedomain.ForgeTypeGitea:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyGiteaToken, args.repo.UnvalidatedConfig.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, args.inputs.Next())
		case forgedomain.ForgeTypeGitHub:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyGitHubToken, args.repo.UnvalidatedConfig.NormalConfig.GitHubToken)
			return dialog.TokenScope(existingScope, args.inputs.Next())
		case forgedomain.ForgeTypeGitLab:
			existingScope := determineExistingScope(args.repo.ConfigSnapshot, configdomain.KeyGitLabToken, args.repo.UnvalidatedConfig.NormalConfig.GitLabToken)
			return dialog.TokenScope(existingScope, args.inputs.Next())
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
		config:        repo.UnvalidatedConfig,
		configFile:    repo.UnvalidatedConfig.File,
		dialogInputs:  dialogTestInputs,
		localBranches: branchesSnapshot.Branches,
		remotes:       remotes,
	}, exit, nil
}

func saveAll(userInput userInput, existingGitConfig configdomain.PartialConfig, configFile Option[configdomain.PartialConfig], data setupData, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(userInput.data.Aliases, existingGitConfig.Aliases, frontend),
	)
	if forgeType, hasForgeType := userInput.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(userInput.data.BitbucketUsername, existingGitConfig.BitbucketUsername, userInput.scope, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(userInput.data.BitbucketAppPassword, existingGitConfig.BitbucketAppPassword, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeCodeberg:
			fc.Check(
				saveCodebergToken(userInput.data.CodebergToken, existingGitConfig.CodebergToken, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeGitHub:
			fc.Check(
				saveGitHubToken(userInput.data.GitHubToken, existingGitConfig.GitHubToken, userInput.scope, userInput.data.GitHubConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitLab:
			fc.Check(
				saveGitLabToken(userInput.data.GitLabToken, existingGitConfig.GitLabToken, userInput.scope, userInput.data.GitLabConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(userInput.data.GiteaToken, existingGitConfig.GiteaToken, userInput.scope, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.storageLocation {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, existingGitConfig, frontend)
	case dialog.ConfigStorageOptionGit: //
		return saveToGit(userInput, existingGitConfig, configFile, data, frontend)
	}
	return nil
}

func saveToGit(userInput userInput, existingGitConfig configdomain.PartialConfig, configFileOpt Option[configdomain.PartialConfig], data setupData, frontend subshelldomain.Runner) error {
	configFile := configFileOpt.GetOrDefault()
	ec := NewMutable(&gohacks.ErrorCollector{})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.BranchType]{
		configFileValue:   configFile.NewBranchType,
		saveFunc:          gitconfig.SetNewBranchType,
		removeFunc:        gitconfig.RemoveNewBranchType,
		valueToWrite:      userInput.data.NewBranchType,
		valueAlreadyInGit: existingGitConfig.NewBranchType,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[forgedomain.ForgeType]{
		configFileValue:   configFile.ForgeType,
		saveFunc:          gitconfig.SetForgeType,
		removeFunc:        gitconfig.RemoveForgeType,
		valueToWrite:      userInput.data.ForgeType,
		valueAlreadyInGit: existingGitConfig.ForgeType,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[forgedomain.GitHubConnectorType]{
		configFileValue:   configFile.GitHubConnectorType,
		saveFunc:          gitconfig.SetGitHubConnectorType,
		removeFunc:        gitconfig.RemoveGitHubConnectorType,
		valueToWrite:      userInput.data.GitHubConnectorType,
		valueAlreadyInGit: existingGitConfig.GitHubConnectorType,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[forgedomain.GitLabConnectorType]{
		configFileValue:   configFile.GitLabConnectorType,
		saveFunc:          gitconfig.SetGitLabConnectorType,
		removeFunc:        gitconfig.RemoveGitLabConnectorType,
		valueToWrite:      userInput.data.GitLabConnectorType,
		valueAlreadyInGit: existingGitConfig.GitLabConnectorType,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.HostingOriginHostname]{
		configFileValue:   configFile.HostingOriginHostname,
		saveFunc:          gitconfig.SetOriginHostname,
		removeFunc:        gitconfig.RemoveOriginHostname,
		valueToWrite:      userInput.data.HostingOriginHostname,
		valueAlreadyInGit: existingGitConfig.HostingOriginHostname,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[gitdomain.LocalBranchName]{
		configFileValue:   configFile.MainBranch,
		saveFunc:          gitconfig.SetMainBranch,
		removeFunc:        gitconfig.RemoveMainBranch,
		valueToWrite:      userInput.data.MainBranch,
		valueAlreadyInGit: existingGitConfig.MainBranch,
	})
	saveCollectionToLocalGit(ec, frontend, saveCollectionArgs[gitdomain.LocalBranchNames, gitdomain.LocalBranchName]{
		configFileValue:   configFile.PerennialBranches,
		saveFunc:          gitconfig.SetPerennialBranches,
		removeFunc:        gitconfig.RemovePerennialBranches,
		valueToWrite:      userInput.data.PerennialBranches,
		valueAlreadyInGit: existingGitConfig.PerennialBranches,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.PerennialRegex]{
		configFileValue:   configFile.PerennialRegex,
		saveFunc:          gitconfig.SetPerennialRegex,
		removeFunc:        gitconfig.RemovePerennialRegex,
		valueToWrite:      userInput.data.PerennialRegex,
		valueAlreadyInGit: existingGitConfig.PerennialRegex,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.BranchType]{
		configFileValue:   configFile.UnknownBranchType,
		saveFunc:          gitconfig.SetUnknownBranchType,
		removeFunc:        gitconfig.RemoveUnknownBranchType,
		valueToWrite:      userInput.data.UnknownBranchType,
		valueAlreadyInGit: existingGitConfig.UnknownBranchType,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[gitdomain.Remote]{
		configFileValue:   configFile.DevRemote,
		saveFunc:          gitconfig.SetDevRemote,
		removeFunc:        gitconfig.RemoveDevRemote,
		valueToWrite:      userInput.data.DevRemote,
		valueAlreadyInGit: existingGitConfig.DevRemote,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.FeatureRegex]{
		configFileValue:   configFile.FeatureRegex,
		saveFunc:          gitconfig.SetFeatureRegex,
		removeFunc:        gitconfig.RemoveFeatureRegex,
		valueToWrite:      userInput.data.FeatureRegex,
		valueAlreadyInGit: existingGitConfig.FeatureRegex,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.ContributionRegex]{
		configFileValue:   configFile.ContributionRegex,
		saveFunc:          gitconfig.SetContributionRegex,
		removeFunc:        gitconfig.RemoveContributionRegex,
		valueToWrite:      userInput.data.ContributionRegex,
		valueAlreadyInGit: existingGitConfig.ContributionRegex,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.ObservedRegex]{
		configFileValue:   configFile.ObservedRegex,
		saveFunc:          gitconfig.SetObservedRegex,
		removeFunc:        gitconfig.RemoveObservedRegex,
		valueToWrite:      userInput.data.ObservedRegex,
		valueAlreadyInGit: existingGitConfig.ObservedRegex,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.PushHook]{
		configFileValue:   configFile.PushHook,
		saveFunc:          gitconfig.SetPushHook,
		removeFunc:        gitconfig.RemovePushHook,
		valueToWrite:      userInput.data.PushHook,
		valueAlreadyInGit: existingGitConfig.PushHook,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.ShareNewBranches]{
		configFileValue:   configFile.ShareNewBranches,
		saveFunc:          gitconfig.SetShareNewBranches,
		removeFunc:        gitconfig.RemoveShareNewBranches,
		valueToWrite:      userInput.data.ShareNewBranches,
		valueAlreadyInGit: existingGitConfig.ShareNewBranches,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.ShipStrategy]{
		configFileValue:   configFile.ShipStrategy,
		saveFunc:          gitconfig.SetShipStrategy,
		removeFunc:        gitconfig.RemoveShipStrategy,
		valueToWrite:      userInput.data.ShipStrategy,
		valueAlreadyInGit: existingGitConfig.ShipStrategy,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.ShipDeleteTrackingBranch]{
		configFileValue:   configFile.ShipDeleteTrackingBranch,
		saveFunc:          gitconfig.SetShipDeleteTrackingBranch,
		removeFunc:        gitconfig.RemoveShipDeleteTrackingBranch,
		valueToWrite:      userInput.data.ShipDeleteTrackingBranch,
		valueAlreadyInGit: existingGitConfig.ShipDeleteTrackingBranch,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.SyncFeatureStrategy]{
		configFileValue:   configFile.SyncFeatureStrategy,
		saveFunc:          gitconfig.SetSyncFeatureStrategy,
		removeFunc:        gitconfig.RemoveSyncFeatureStrategy,
		valueToWrite:      userInput.data.SyncFeatureStrategy,
		valueAlreadyInGit: existingGitConfig.SyncFeatureStrategy,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.SyncPerennialStrategy]{
		configFileValue:   configFile.SyncPerennialStrategy,
		saveFunc:          gitconfig.SetSyncPerennialStrategy,
		removeFunc:        gitconfig.RemoveSyncPerennialStrategy,
		valueToWrite:      userInput.data.SyncPerennialStrategy,
		valueAlreadyInGit: existingGitConfig.SyncPerennialStrategy,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.SyncPrototypeStrategy]{
		configFileValue:   configFile.SyncPrototypeStrategy,
		saveFunc:          gitconfig.SetSyncPrototypeStrategy,
		removeFunc:        gitconfig.RemoveSyncPrototypeStrategy,
		valueToWrite:      userInput.data.SyncPrototypeStrategy,
		valueAlreadyInGit: existingGitConfig.SyncPrototypeStrategy,
	})
	saveOptionToLocalGit(ec, frontend, saveToLocalGitArgs[configdomain.SyncUpstream]{
		configFileValue:   configFile.SyncUpstream,
		saveFunc:          gitconfig.SetSyncUpstream,
		removeFunc:        gitconfig.RemoveSyncUpstream,
		valueToWrite:      userInput.data.SyncUpstream,
		valueAlreadyInGit: existingGitConfig.SyncUpstream,
	})
	if configFile.SyncTags.IsNone() {
		ec.Check(
			saveSyncTags(userInput.data.SyncTags, existingGitConfig.SyncTags, frontend),
		)
	}
	return ec.Err
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

func saveOptionToLocalGit[T comparable](ec Mutable[gohacks.ErrorCollector], runner subshelldomain.Runner, args saveToLocalGitArgs[T]) {
	if ec.Value.Err != nil {
		return
	}
	if args.valueToWrite.Equal(args.configFileValue) {
		return
	}
	if args.valueToWrite.Equal(args.valueAlreadyInGit) {
		return
	}
	if value, has := args.valueToWrite.Get(); has {
		ec.Value.Check(args.saveFunc(runner, value, configdomain.ConfigScopeLocal))
	}
	ec.Value.Check(gitconfig.RemoveBitbucketAppPassword(runner))
}

type saveToLocalGitArgs[T comparable] struct {
	configFileValue   Option[T]
	valueToWrite      Option[T]
	valueAlreadyInGit Option[T]
	saveFunc          func(subshelldomain.Runner, T, configdomain.ConfigScope) error
	removeFunc        func(subshelldomain.Runner) error
}

func saveCollectionToLocalGit[T ~[]E, E cmp.Ordered](ec Mutable[gohacks.ErrorCollector], runner subshelldomain.Runner, args saveCollectionArgs[T, E]) {
	if ec.Value.Err != nil {
		return
	}
	if slices.Compare(args.valueToWrite, args.configFileValue) == 0 {
		return
	}
	if slices.Compare(args.valueToWrite, args.valueAlreadyInGit) == 0 {
		return
	}
	ec.Value.Check(args.saveFunc(runner, args.valueToWrite, configdomain.ConfigScopeLocal))
}

type saveCollectionArgs[T ~[]E, E cmp.Ordered] struct {
	configFileValue   T
	valueToWrite      T
	valueAlreadyInGit T
	saveFunc          func(subshelldomain.Runner, T, configdomain.ConfigScope) error
	removeFunc        func(subshelldomain.Runner) error
}

type saveCommonArgs struct {
	failureCounter Mutable[gohacks.ErrorCollector]
	runner         subshelldomain.Runner
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

func saveNewBranchType(valueToWriteToGit Option[configdomain.BranchType], valueAlreadyInGit Option[configdomain.BranchType], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, hasValue := valueToWriteToGit.Get(); hasValue {
		return gitconfig.SetNewBranchType(runner, value)
	}
	_ = gitconfig.RemoveNewBranchType(runner)
	return nil
}

func saveUnknownBranchType(valueToWriteToGit Option[configdomain.BranchType], valueAlreadyInGit Option[configdomain.BranchType], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetUnknownBranchType(runner, value)
	}
	return gitconfig.RemoveUnknownBranchType(runner)
}

func saveDevRemote(valueToWriteToGit gitdomain.Remote, valueAlreadyInGit Option[gitdomain.Remote], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetDevRemote(runner, valueToWriteToGit)
}

func saveFeatureRegex(valueToWriteToGit Option[configdomain.FeatureRegex], valueAlreadyInGit Option[configdomain.FeatureRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetFeatureRegex(runner, value)
	}
	_ = gitconfig.RemoveFeatureRegex(runner)
	return nil
}

func saveContributionRegex(valueToWriteToGit Option[configdomain.ContributionRegex], valueAlreadyInGit Option[configdomain.ContributionRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetContributionRegex(runner, value)
	}
	_ = gitconfig.RemoveContributionRegex(runner)
	return nil
}

func saveObservedRegex(valueToWriteToGit Option[configdomain.ObservedRegex], valueAlreadyInGit Option[configdomain.ObservedRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetObservedRegex(runner, value)
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
		return gitconfig.SetForgeType(frontend, newValue)
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
		return gitconfig.SetGitHubConnectorType(frontend, value)
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
		return gitconfig.SetGitLabConnectorType(frontend, value)
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

func saveMainBranch(valueToWriteToGit Option[gitdomain.LocalBranchName], valueAlreadyInGit Option[gitdomain.LocalBranchName], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.Equal(valueToWriteToGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetMainBranch(runner, value)
	}
	return gitconfig.RemoveMainBranch(runner)
}

func saveOriginHostname(valueToWriteToGit Option[configdomain.HostingOriginHostname], valueAlreadyInGit Option[configdomain.HostingOriginHostname], frontend subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetOriginHostname(frontend, value)
	}
	return gitconfig.RemoveOriginHostname(frontend)
}

func savePerennialBranches(valueToWriteToGit gitdomain.LocalBranchNames, valueAlreadyInGit gitdomain.LocalBranchNames, runner subshelldomain.Runner) error {
	if slices.Compare(valueAlreadyInGit, valueToWriteToGit) == 0 {
		return nil
	}
	return gitconfig.SetPerennialBranches(runner, valueToWriteToGit)
}

func savePerennialRegex(valueToWriteToGit Option[configdomain.PerennialRegex], valueAlreadyInGit Option[configdomain.PerennialRegex], runner subshelldomain.Runner) error {
	if valueToWriteToGit.Equal(valueAlreadyInGit) {
		return nil
	}
	if value, has := valueToWriteToGit.Get(); has {
		return gitconfig.SetPerennialRegex(runner, value)
	}
	_ = gitconfig.RemovePerennialRegex(runner)
	return nil
}

func savePushHook(valueToWriteToGit configdomain.PushHook, valueAlreadyInGit Option[configdomain.PushHook], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetPushHook(runner, valueToWriteToGit)
}

func saveShareNewBranches(valueToWriteToGit configdomain.ShareNewBranches, valueAlreadyInGit Option[configdomain.ShareNewBranches], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetShareNewBranches(runner, valueToWriteToGit)
}

func saveShipDeleteTrackingBranch(valueToWriteToGit configdomain.ShipDeleteTrackingBranch, valueAlreadyInGit Option[configdomain.ShipDeleteTrackingBranch], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetShipDeleteTrackingBranch(runner, valueToWriteToGit)
}

func saveShipStrategy(valueToWriteToGit configdomain.ShipStrategy, valueAlreadyInGit Option[configdomain.ShipStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetShipStrategy(runner, valueToWriteToGit)
}

func saveSyncFeatureStrategy(valueToWriteToGit configdomain.SyncFeatureStrategy, valueAlreadyInGit Option[configdomain.SyncFeatureStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetSyncFeatureStrategy(runner, valueToWriteToGit)
}

func saveSyncPerennialStrategy(valueToWriteToGit configdomain.SyncPerennialStrategy, valueAlreadyInGit Option[configdomain.SyncPerennialStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetSyncPerennialStrategy(runner, valueToWriteToGit)
}

func saveSyncPrototypeStrategy(valueToWriteToGit configdomain.SyncPrototypeStrategy, valueAlreadyInGit Option[configdomain.SyncPrototypeStrategy], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetSyncPrototypeStrategy(runner, valueToWriteToGit)
}

func saveSyncUpstream(valueToWriteToGit configdomain.SyncUpstream, valueAlreadyInGit Option[configdomain.SyncUpstream], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetSyncUpstream(runner, valueToWriteToGit)
}

func saveSyncTags(valueToWriteToGit configdomain.SyncTags, valueAlreadyInGit Option[configdomain.SyncTags], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetSyncTags(runner, valueToWriteToGit)
}

func saveToFile(userInput userInput, gitConfig configdomain.PartialConfig, runner subshelldomain.Runner) error {
	if err := configfile.Save(userInput.normalConfig, userInput.validatedConfig.MainBranch); err != nil {
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
	if err := saveUnknownBranchType(userInput.normalConfig.UnknownBranchType, gitConfig.UnknownBranchType, runner); err != nil {
		return err
	}
	// TODO: also save ObservedRegex ContributionRegex NewBranchType
	return saveFeatureRegex(userInput.normalConfig.FeatureRegex, gitConfig.FeatureRegex, runner)
}

func wrapParseFunc[T any](parseFunc func(arg string) Option[T]) func(string) (Option[T], error) {
	return func(arg string) (Option[T], error) {
		return parseFunc(arg), nil
	}
}
