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
	if err = saveAll(enterDataResult, repo.UnvalidatedConfig, data.configFile, repo.Frontend); err != nil {
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
	var mainBranch gitdomain.LocalBranchName
	if configFileMainBranch, configFileHasMainBranch := configFile.MainBranch.Get(); configFileHasMainBranch {
		mainBranch = configFileMainBranch
	} else {
		existingMainBranch := repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch
		if existingMainBranch.IsNone() {
			existingMainBranch = gitconfig.DefaultBranch(repo.Backend)
		}
		if existingMainBranch.IsNone() {
			existingMainBranch = repo.Git.OriginHead(repo.Backend)
		}
		mainBranch, exit, err = dialog.MainBranch(data.localBranches.Names(), existingMainBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	perennialBranches := repo.UnvalidatedConfig.NormalConfig.PerennialBranches
	if len(configFile.PerennialBranches) == 0 {
		perennialBranches, exit, err = dialog.PerennialBranches(data.localBranches.Names(), perennialBranches, mainBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	perennialRegex := repo.UnvalidatedConfig.NormalConfig.PerennialRegex
	if configFile.PerennialRegex.IsNone() {
		perennialRegex, exit, err = dialog.PerennialRegex(perennialRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	featureRegex := repo.UnvalidatedConfig.NormalConfig.FeatureRegex
	if configFile.FeatureRegex.IsNone() {
		featureRegex, exit, err = dialog.FeatureRegex(featureRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	contributionRegex := repo.UnvalidatedConfig.NormalConfig.ContributionRegex
	if configFile.ContributionRegex.IsNone() {
		contributionRegex, exit, err = dialog.ContributionRegex(contributionRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	observedRegex := repo.UnvalidatedConfig.NormalConfig.ObservedRegex
	if configFile.ObservedRegex.IsNone() {
		observedRegex, exit, err = dialog.ObservedRegex(observedRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	unknownBranchType := repo.UnvalidatedConfig.NormalConfig.UnknownBranchType
	if configFile.UnknownBranchType.IsNone() {
		unknownBranchType, exit, err = dialog.UnknownBranchType(unknownBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	devRemote := repo.UnvalidatedConfig.NormalConfig.DevRemote
	if configFile.DevRemote.IsNone() {
		devRemote, exit, err = dialog.DevRemote(devRemote, data.remotes, data.dialogInputs.Next())
		if err != nil || exit {
			return emptyResult, exit, err
		}
	}
	hostingOriginHostName := repo.UnvalidatedConfig.NormalConfig.HostingOriginHostname
	enteredForgeType := repo.UnvalidatedConfig.NormalConfig.ForgeType.Or(repo.UnvalidatedConfig.File.GetOrDefault().ForgeType)
	var actualForgeType Option[forgedomain.ForgeType]
	bitbucketUsername := repo.UnvalidatedConfig.NormalConfig.BitbucketUsername
	bitbucketAppPassword := repo.UnvalidatedConfig.NormalConfig.BitbucketAppPassword
	codebergToken := repo.UnvalidatedConfig.NormalConfig.CodebergToken
	devURL := data.config.NormalConfig.DevURL(data.backend)
	giteaToken := repo.UnvalidatedConfig.NormalConfig.GiteaToken
	githubConnectorTypeOpt := repo.UnvalidatedConfig.NormalConfig.GitHubConnectorType
	githubToken := repo.UnvalidatedConfig.NormalConfig.GitHubToken
	gitlabConnectorTypeOpt := repo.UnvalidatedConfig.NormalConfig.GitLabConnectorType
	gitlabToken := repo.UnvalidatedConfig.NormalConfig.GitLabToken
	for {
		if configFile.HostingOriginHostname.IsNone() {
			hostingOriginHostName, exit, err = dialog.OriginHostname(hostingOriginHostName, data.dialogInputs.Next())
			if err != nil || exit {
				return emptyResult, exit, err
			}
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
				bitbucketUsername, exit, err = dialog.BitbucketUsername(bitbucketUsername, data.dialogInputs.Next())
				if err != nil || exit {
					return emptyResult, exit, err
				}
				bitbucketAppPassword, exit, err = dialog.BitbucketAppPassword(bitbucketAppPassword, data.dialogInputs.Next())
			case forgedomain.ForgeTypeCodeberg:
				codebergToken, exit, err = dialog.CodebergToken(codebergToken, data.dialogInputs.Next())
			case forgedomain.ForgeTypeGitea:
				giteaToken, exit, err = dialog.GiteaToken(giteaToken, data.dialogInputs.Next())
			case forgedomain.ForgeTypeGitHub:
				githubConnectorTypeOpt, exit, err = dialog.GitHubConnectorType(githubConnectorTypeOpt, data.dialogInputs.Next())
				if err != nil || exit {
					return emptyResult, exit, err
				}
				if githubConnectorType, has := githubConnectorTypeOpt.Get(); has {
					switch githubConnectorType {
					case forgedomain.GitHubConnectorTypeAPI:
						githubToken, exit, err = dialog.GitHubToken(githubToken, data.dialogInputs.Next())
					case forgedomain.GitHubConnectorTypeGh:
					}
				}
			case forgedomain.ForgeTypeGitLab:
				gitlabConnectorTypeOpt, exit, err = dialog.GitLabConnectorType(gitlabConnectorTypeOpt, data.dialogInputs.Next())
				if err != nil || exit {
					return emptyResult, exit, err
				}
				if gitlabConnectorType, has := gitlabConnectorTypeOpt.Get(); has {
					switch gitlabConnectorType {
					case forgedomain.GitLabConnectorTypeAPI:
						gitlabToken, exit, err = dialog.GitLabToken(gitlabToken, data.dialogInputs.Next())
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
			remoteURL:            data.config.NormalConfig.RemoteURL(data.backend, devRemote),
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
	syncPrototypeStrategy := repo.UnvalidatedConfig.NormalConfig.SyncPrototypeStrategy
	if configFile.SyncPrototypeStrategy.IsNone() {
		syncPrototypeStrategy, exit, err = dialog.SyncPrototypeStrategy(syncPrototypeStrategy, data.dialogInputs.Next())
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
	normalData := configdomain.NormalConfigData{
		Aliases:                  aliases,
		BitbucketAppPassword:     bitbucketAppPassword,
		BitbucketUsername:        bitbucketUsername,
		BranchTypeOverrides:      configdomain.BranchTypeOverrides{}, // the setup assistant doesn't ask for this
		CodebergToken:            codebergToken,
		ContributionRegex:        contributionRegex,
		DevRemote:                devRemote,
		DryRun:                   false, // the setup assistant doesn't ask for this
		FeatureRegex:             featureRegex,
		ForgeType:                enteredForgeType,
		GitHubConnectorType:      githubConnectorTypeOpt,
		GitHubToken:              githubToken,
		GitLabConnectorType:      gitlabConnectorTypeOpt,
		GitLabToken:              gitlabToken,
		GiteaToken:               giteaToken,
		HostingOriginHostname:    hostingOriginHostName,
		Lineage:                  configdomain.Lineage{}, // the setup assistant doesn't ask for this
		NewBranchType:            newBranchType,
		ObservedRegex:            observedRegex,
		Offline:                  false, // the setup assistant doesn't ask for this
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
		Verbose:                  false, // the setup assistant doesn't ask for this
	}
	validatedData := configdomain.ValidatedConfigData{
		GitUserEmail: "", // the setup assistant doesn't ask for this
		GitUserName:  "", // the setup assistant doesn't ask for this
		MainBranch:   mainBranch,
	}
	return userInput{actualForgeType, normalData, tokenScope, configStorage, validatedData}, false, nil
}

// data entered by the user in the setup assistant
type userInput struct {
	determinedForgeType Option[forgedomain.ForgeType] // the forge type that was determined by the setup assistant - not necessarily what the user entered (could also be "auto detect")
	normalConfig        configdomain.NormalConfigData
	scope               configdomain.ConfigScope
	storageLocation     dialog.ConfigStorageOption
	validatedConfig     configdomain.ValidatedConfigData
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

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, configFile Option[configdomain.PartialConfig], frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(userInput.normalConfig.Aliases, oldConfig.NormalConfig, frontend),
	)
	if forgeType, hasForgeType := userInput.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(userInput.normalConfig.BitbucketUsername, oldConfig.NormalConfig, userInput.scope, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(userInput.normalConfig.BitbucketAppPassword, oldConfig.NormalConfig, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeCodeberg:
			fc.Check(
				saveCodebergToken(userInput.normalConfig.CodebergToken, oldConfig.NormalConfig, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeGitHub:
			fc.Check(
				saveGitHubToken(userInput.normalConfig.GitHubToken, oldConfig.NormalConfig, userInput.scope, userInput.normalConfig.GitHubConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitLab:
			fc.Check(
				saveGitLabToken(userInput.normalConfig.GitLabToken, oldConfig.NormalConfig, userInput.scope, userInput.normalConfig.GitLabConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(userInput.normalConfig.GiteaToken, oldConfig.NormalConfig, userInput.scope, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.storageLocation {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, oldConfig.NormalConfig, frontend)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, oldConfig, configFile, frontend)
	}
	return nil
}

func saveToGit(userInput userInput, oldConfig config.UnvalidatedConfig, configFileOpt Option[configdomain.PartialConfig], frontend subshelldomain.Runner) error {
	configFile := configFileOpt.GetOrDefault()
	fc := gohacks.ErrorCollector{}
	if configFile.NewBranchType.IsNone() {
		fc.Check(
			saveNewBranchType(userInput.normalConfig.NewBranchType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(
			saveForgeType(userInput.normalConfig.ForgeType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.GitHubConnectorType.IsNone() {
		fc.Check(
			saveGitHubConnectorType(userInput.normalConfig.GitHubConnectorType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.GitLabConnectorType.IsNone() {
		fc.Check(
			saveGitLabConnectorType(userInput.normalConfig.GitLabConnectorType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(
			saveOriginHostname(userInput.normalConfig.HostingOriginHostname, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(
			saveMainBranch(userInput.validatedConfig.MainBranch, oldConfig, frontend),
		)
	}
	if len(configFile.PerennialBranches) == 0 {
		fc.Check(
			savePerennialBranches(userInput.normalConfig.PerennialBranches, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(
			savePerennialRegex(userInput.normalConfig.PerennialRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(
			saveUnknownBranchType(userInput.normalConfig.UnknownBranchType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.DevRemote.IsNone() {
		fc.Check(
			saveDevRemote(userInput.normalConfig.DevRemote, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(
			saveFeatureRegex(userInput.normalConfig.FeatureRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ContributionRegex.IsNone() {
		fc.Check(
			saveContributionRegex(userInput.normalConfig.ContributionRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ObservedRegex.IsNone() {
		fc.Check(
			saveObservedRegex(userInput.normalConfig.ObservedRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.PushHook.IsNone() {
		fc.Check(
			savePushHook(userInput.normalConfig.PushHook, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(
			saveShareNewBranches(userInput.normalConfig.ShareNewBranches, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(
			saveShipStrategy(userInput.normalConfig.ShipStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(
			saveShipDeleteTrackingBranch(userInput.normalConfig.ShipDeleteTrackingBranch, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(
			saveSyncFeatureStrategy(userInput.normalConfig.SyncFeatureStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(
			saveSyncPerennialStrategy(userInput.normalConfig.SyncPerennialStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(
			saveSyncPrototypeStrategy(userInput.normalConfig.SyncPrototypeStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(
			saveSyncUpstream(userInput.normalConfig.SyncUpstream, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(
			saveSyncTags(userInput.normalConfig.SyncTags, oldConfig.NormalConfig, frontend),
		)
	}
	return fc.Err
}

func saveAliases(values configdomain.Aliases, config config.NormalConfig, frontend subshelldomain.Runner) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := config.Aliases[aliasableCommand]
		newAlias, hasNew := values[aliasableCommand]
		switch {
		case hasOld && !hasNew:
			err = gitconfig.RemoveAlias(frontend, aliasableCommand)
		case newAlias != oldAlias:
			err = gitconfig.SetAlias(frontend, aliasableCommand)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func saveBitbucketAppPassword(value Option[forgedomain.BitbucketAppPassword], config config.NormalConfig, scope configdomain.ConfigScope, runner subshelldomain.Runner) error {
	if value.Equal(config.BitbucketAppPassword) {
		return nil
	}
	if value, has := value.Get(); has {
		return gitconfig.SetBitbucketAppPassword(runner, value, scope)
	}
	return gitconfig.RemoveBitbucketAppPassword(runner)
}

func saveBitbucketUsername(newValue Option[forgedomain.BitbucketUsername], config config.NormalConfig, scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if newValue.Equal(config.BitbucketUsername) {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitconfig.SetBitbucketUsername(frontend, value, scope)
	}
	return gitconfig.RemoveBitbucketUsername(frontend)
}

func saveNewBranchType(newValue Option[configdomain.BranchType], config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue.Equal(config.NewBranchType) {
		return nil
	}
	if value, hasValue := newValue.Get(); hasValue {
		return gitconfig.SetNewBranchType(runner, value)
	}
	_ = gitconfig.RemoveNewBranchType(runner)
	return nil
}

func saveUnknownBranchType(value configdomain.BranchType, config config.NormalConfig, runner subshelldomain.Runner) error {
	if value == config.UnknownBranchType {
		return nil
	}
	return gitconfig.SetUnknownBranchType(runner, value)
}

func saveDevRemote(value gitdomain.Remote, config config.NormalConfig, runner subshelldomain.Runner) error {
	if value == config.DevRemote {
		return nil
	}
	return gitconfig.SetDevRemote(runner, value)
}

func saveFeatureRegex(value Option[configdomain.FeatureRegex], config config.NormalConfig, runner subshelldomain.Runner) error {
	if value.Equal(config.FeatureRegex) {
		return nil
	}
	if value, has := value.Get(); has {
		return gitconfig.SetFeatureRegex(runner, value)
	}
	_ = gitconfig.RemoveFeatureRegex(runner)
	return nil
}

func saveContributionRegex(value Option[configdomain.ContributionRegex], config config.NormalConfig, runner subshelldomain.Runner) error {
	if value.Equal(config.ContributionRegex) {
		return nil
	}
	if value, has := value.Get(); has {
		return gitconfig.SetContributionRegex(runner, value)
	}
	_ = gitconfig.RemoveContributionRegex(runner)
	return nil
}

func saveObservedRegex(value Option[configdomain.ObservedRegex], config config.NormalConfig, runner subshelldomain.Runner) error {
	if value.Equal(config.ObservedRegex) {
		return nil
	}
	if value, has := value.Get(); has {
		return gitconfig.SetObservedRegex(runner, value)
	}
	_ = gitconfig.RemoveObservedRegex(runner)
	return nil
}

func saveForgeType(value Option[forgedomain.ForgeType], config config.NormalConfig, frontend subshelldomain.Runner) (err error) {
	oldValue, oldHas := config.ForgeType.Get()
	newValue, newHas := value.Get()
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

func saveCodebergToken(newToken Option[forgedomain.CodebergToken], config config.NormalConfig, scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if newToken.Equal(config.CodebergToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitconfig.SetCodebergToken(frontend, value, scope)
	}
	return gitconfig.RemoveCodebergToken(frontend)
}

func saveGiteaToken(newToken Option[forgedomain.GiteaToken], config config.NormalConfig, scope configdomain.ConfigScope, frontend subshelldomain.Runner) error {
	if newToken.Equal(config.GiteaToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitconfig.SetGiteaToken(frontend, value, scope)
	}
	return gitconfig.RemoveGiteaToken(frontend)
}

func saveGitHubConnectorType(newType Option[forgedomain.GitHubConnectorType], config config.NormalConfig, frontend subshelldomain.Runner) error {
	if newType.Equal(config.GitHubConnectorType) {
		return nil
	}
	if value, has := newType.Get(); has {
		return gitconfig.SetGitHubConnectorType(frontend, value)
	}
	return gitconfig.RemoveGitHubConnectorType(frontend)
}

func saveGitHubToken(newToken Option[forgedomain.GitHubToken], config config.NormalConfig, scope configdomain.ConfigScope, githubConnectorType Option[forgedomain.GitHubConnectorType], frontend subshelldomain.Runner) error {
	if connectorType, has := githubConnectorType.Get(); has {
		if connectorType == forgedomain.GitHubConnectorTypeGh {
			return nil
		}
	}
	if newToken.Equal(config.GitHubToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitconfig.SetGitHubToken(frontend, value, scope)
	}
	return gitconfig.RemoveGitHubToken(frontend)
}

func saveGitLabConnectorType(newType Option[forgedomain.GitLabConnectorType], config config.NormalConfig, frontend subshelldomain.Runner) error {
	if newType.Equal(config.GitLabConnectorType) {
		return nil
	}
	if value, has := newType.Get(); has {
		return gitconfig.SetGitLabConnectorType(frontend, value)
	}
	return gitconfig.RemoveGitLabConnectorType(frontend)
}

func saveGitLabToken(newToken Option[forgedomain.GitLabToken], config config.NormalConfig, scope configdomain.ConfigScope, gitlabConnectorType Option[forgedomain.GitLabConnectorType], frontend subshelldomain.Runner) error {
	if connectorType, has := gitlabConnectorType.Get(); has {
		if connectorType == forgedomain.GitLabConnectorTypeGlab {
			return nil
		}
	}
	if newToken.Equal(config.GitLabToken) {
		return nil
	}
	if value, has := newToken.Get(); has {
		return gitconfig.SetGitLabToken(frontend, value, scope)
	}
	return gitconfig.RemoveGitLabToken(frontend)
}

func saveMainBranch(value gitdomain.LocalBranchName, config config.UnvalidatedConfig, runner subshelldomain.Runner) error {
	if existing, hasExisting := config.UnvalidatedConfig.MainBranch.Get(); hasExisting {
		if existing == value {
			return nil
		}
	}
	return config.SetMainBranch(value, runner)
}

func saveOriginHostname(newValue Option[configdomain.HostingOriginHostname], config config.NormalConfig, frontend subshelldomain.Runner) error {
	if newValue.Equal(config.HostingOriginHostname) {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitconfig.SetOriginHostname(frontend, value)
	}
	return gitconfig.RemoveOriginHostname(frontend)
}

func savePerennialBranches(newValue gitdomain.LocalBranchNames, config config.NormalConfig, runner subshelldomain.Runner) error {
	if slices.Compare(config.PerennialBranches, newValue) != 0 || config.Git.PerennialBranches == nil {
		return gitconfig.SetPerennialBranches(runner, newValue)
	}
	return nil
}

func savePerennialRegex(newValue Option[configdomain.PerennialRegex], config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue.Equal(config.PerennialRegex) {
		return nil
	}
	if value, has := newValue.Get(); has {
		return gitconfig.SetPerennialRegex(runner, value)
	}
	_ = gitconfig.RemovePerennialRegex(runner)
	return nil
}

func savePushHook(newValue configdomain.PushHook, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.PushHook {
		return nil
	}
	return gitconfig.SetPushHook(runner, newValue)
}

func saveShareNewBranches(newValue configdomain.ShareNewBranches, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.ShareNewBranches {
		return nil
	}
	return gitconfig.SetShareNewBranches(runner, newValue)
}

func saveShipDeleteTrackingBranch(newValue configdomain.ShipDeleteTrackingBranch, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.ShipDeleteTrackingBranch {
		return nil
	}
	return gitconfig.SetShipDeleteTrackingBranch(runner, newValue)
}

func saveShipStrategy(newValue configdomain.ShipStrategy, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.ShipStrategy {
		return nil
	}
	return gitconfig.SetShipStrategy(runner, newValue)
}

func saveSyncFeatureStrategy(newValue configdomain.SyncFeatureStrategy, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.SyncFeatureStrategy {
		return nil
	}
	return gitconfig.SetSyncFeatureStrategy(runner, newValue)
}

func saveSyncPerennialStrategy(newValue configdomain.SyncPerennialStrategy, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.SyncPerennialStrategy {
		return nil
	}
	return gitconfig.SetSyncPerennialStrategy(runner, newValue)
}

func saveSyncPrototypeStrategy(newValue configdomain.SyncPrototypeStrategy, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.SyncPrototypeStrategy {
		return nil
	}
	return gitconfig.SetSyncPrototypeStrategy(runner, newValue)
}

func saveSyncUpstream(newValue configdomain.SyncUpstream, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.SyncUpstream {
		return nil
	}
	return gitconfig.SetSyncUpstream(runner, newValue)
}

func saveSyncTags(newValue configdomain.SyncTags, config config.NormalConfig, runner subshelldomain.Runner) error {
	if newValue == config.SyncTags {
		return nil
	}
	return gitconfig.SetSyncTags(runner, newValue)
}

func saveToFile(userInput userInput, config config.NormalConfig, runner subshelldomain.Runner) error {
	if err := configfile.Save(userInput.normalConfig, userInput.validatedConfig.MainBranch); err != nil {
		return err
	}
	if config.Git.DevRemote.IsSome() {
		_ = gitconfig.RemoveDevRemote(runner)
	}
	if config.Git.MainBranch.IsSome() {
		_ = gitconfig.RemoveMainBranch(runner)
	}
	if config.Git.NewBranchType.IsSome() {
		_ = gitconfig.RemoveNewBranchType(runner)
	}
	if len(config.Git.PerennialBranches) > 0 {
		_ = gitconfig.RemovePerennialBranches(runner)
	}
	if config.Git.PerennialRegex.IsSome() {
		_ = gitconfig.RemovePerennialRegex(runner)
	}
	if config.Git.ShareNewBranches.IsSome() {
		_ = gitconfig.RemoveShareNewBranches(runner)
	}
	if config.Git.PushHook.IsSome() {
		_ = gitconfig.RemovePushHook(runner)
	}
	if config.Git.ShipStrategy.IsSome() {
		_ = gitconfig.RemoveShipStrategy(runner)
	}
	if config.Git.ShipDeleteTrackingBranch.IsSome() {
		_ = gitconfig.RemoveShipDeleteTrackingBranch(runner)
	}
	if config.Git.SyncFeatureStrategy.IsSome() {
		_ = gitconfig.RemoveSyncFeatureStrategy(runner)
	}
	if config.Git.SyncPerennialStrategy.IsSome() {
		_ = gitconfig.RemoveSyncPerennialStrategy(runner)
	}
	if config.Git.SyncPrototypeStrategy.IsSome() {
		_ = gitconfig.RemoveSyncPrototypeStrategy(runner)
	}
	if config.Git.SyncUpstream.IsSome() {
		_ = gitconfig.RemoveSyncUpstream(runner)
	}
	if config.Git.SyncTags.IsSome() {
		_ = gitconfig.RemoveSyncTags(runner)
	}
	if err := saveUnknownBranchType(userInput.normalConfig.UnknownBranchType, config, runner); err != nil {
		return err
	}
	return saveFeatureRegex(userInput.normalConfig.FeatureRegex, config, runner)
}
