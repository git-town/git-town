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
	if len(data.remotes) > 1 && configFile.DevRemote.IsNone() {
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

func saveAll(userInput userInput, existingGitConfig configdomain.PartialConfig, configFile Option[configdomain.PartialConfig], data setupData, frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(userInput.normalConfig.Aliases, existingGitConfig.Aliases, frontend),
	)
	if forgeType, hasForgeType := userInput.determinedForgeType.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(userInput.normalConfig.BitbucketUsername, existingGitConfig.BitbucketUsername, userInput.scope, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(userInput.normalConfig.BitbucketAppPassword, existingGitConfig.BitbucketAppPassword, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeCodeberg:
			fc.Check(
				saveCodebergToken(userInput.normalConfig.CodebergToken, existingGitConfig.CodebergToken, userInput.scope, frontend),
			)
		case forgedomain.ForgeTypeGitHub:
			fc.Check(
				saveGitHubToken(userInput.normalConfig.GitHubToken, existingGitConfig.GitHubToken, userInput.scope, userInput.normalConfig.GitHubConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitLab:
			fc.Check(
				saveGitLabToken(userInput.normalConfig.GitLabToken, existingGitConfig.GitLabToken, userInput.scope, userInput.normalConfig.GitLabConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(userInput.normalConfig.GiteaToken, existingGitConfig.GiteaToken, userInput.scope, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.storageLocation {
	case dialog.ConfigStorageOptionFile:
		return saveToFile(userInput, existingGitConfig, frontend)
	case dialog.ConfigStorageOptionGit:
		return saveToGit(userInput, existingGitConfig, configFile, data, frontend)
	}
	return nil
}

func saveToGit(userInput userInput, existingGitConfig configdomain.PartialConfig, configFileOpt Option[configdomain.PartialConfig], data setupData, frontend subshelldomain.Runner) error {
	configFile := configFileOpt.GetOrDefault()
	fc := gohacks.ErrorCollector{}
	if configFile.NewBranchType.IsNone() {
		fc.Check(
			saveNewBranchType(userInput.normalConfig.NewBranchType, existingGitConfig.NewBranchType, frontend),
		)
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(
			saveForgeType(userInput.normalConfig.ForgeType, existingGitConfig.ForgeType, frontend),
		)
	}
	if configFile.GitHubConnectorType.IsNone() {
		fc.Check(
			saveGitHubConnectorType(userInput.normalConfig.GitHubConnectorType, existingGitConfig.GitHubConnectorType, frontend),
		)
	}
	if configFile.GitLabConnectorType.IsNone() {
		fc.Check(
			saveGitLabConnectorType(userInput.normalConfig.GitLabConnectorType, existingGitConfig.GitLabConnectorType, frontend),
		)
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(
			saveOriginHostname(userInput.normalConfig.HostingOriginHostname, existingGitConfig.HostingOriginHostname, frontend),
		)
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(
			saveMainBranch(userInput.validatedConfig.MainBranch, existingGitConfig.MainBranch, frontend),
		)
	}
	if len(configFile.PerennialBranches) == 0 {
		fc.Check(
			savePerennialBranches(userInput.normalConfig.PerennialBranches, existingGitConfig.PerennialBranches, frontend),
		)
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(
			savePerennialRegex(userInput.normalConfig.PerennialRegex, existingGitConfig.PerennialRegex, frontend),
		)
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(
			saveUnknownBranchType(userInput.normalConfig.UnknownBranchType, existingGitConfig.UnknownBranchType, frontend),
		)
	}
	if len(data.remotes) > 1 {
		if configFile.DevRemote.IsNone() {
			fc.Check(
				saveDevRemote(userInput.normalConfig.DevRemote, existingGitConfig.DevRemote, frontend),
			)
		}
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(
			saveFeatureRegex(userInput.normalConfig.FeatureRegex, existingGitConfig.FeatureRegex, frontend),
		)
	}
	if configFile.ContributionRegex.IsNone() {
		fc.Check(
			saveContributionRegex(userInput.normalConfig.ContributionRegex, existingGitConfig.ContributionRegex, frontend),
		)
	}
	if configFile.ObservedRegex.IsNone() {
		fc.Check(
			saveObservedRegex(userInput.normalConfig.ObservedRegex, existingGitConfig.ObservedRegex, frontend),
		)
	}
	if configFile.PushHook.IsNone() {
		fc.Check(
			savePushHook(userInput.normalConfig.PushHook, existingGitConfig.PushHook, frontend),
		)
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(
			saveShareNewBranches(userInput.normalConfig.ShareNewBranches, existingGitConfig.ShareNewBranches, frontend),
		)
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(
			saveShipStrategy(userInput.normalConfig.ShipStrategy, existingGitConfig.ShipStrategy, frontend),
		)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(
			saveShipDeleteTrackingBranch(userInput.normalConfig.ShipDeleteTrackingBranch, existingGitConfig.ShipDeleteTrackingBranch, frontend),
		)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(
			saveSyncFeatureStrategy(userInput.normalConfig.SyncFeatureStrategy, existingGitConfig.SyncFeatureStrategy, frontend),
		)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(
			saveSyncPerennialStrategy(userInput.normalConfig.SyncPerennialStrategy, existingGitConfig.SyncPerennialStrategy, frontend),
		)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(
			saveSyncPrototypeStrategy(userInput.normalConfig.SyncPrototypeStrategy, existingGitConfig.SyncPrototypeStrategy, frontend),
		)
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(
			saveSyncUpstream(userInput.normalConfig.SyncUpstream, existingGitConfig.SyncUpstream, frontend),
		)
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(
			saveSyncTags(userInput.normalConfig.SyncTags, existingGitConfig.SyncTags, frontend),
		)
	}
	return fc.Err
}

func saveAliases(aliasesToWriteToGit configdomain.Aliases, aliasesAlreadyInGit configdomain.Aliases, frontend subshelldomain.Runner) (err error) {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		oldAlias, hasOld := aliasesAlreadyInGit[aliasableCommand]
		newAlias, hasNew := aliasesToWriteToGit[aliasableCommand]
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

func saveUnknownBranchType(valueToWriteToGit configdomain.BranchType, valueAlreadyInGit Option[configdomain.BranchType], runner subshelldomain.Runner) error {
	if valueAlreadyInGit.EqualSome(valueToWriteToGit) {
		return nil
	}
	return gitconfig.SetUnknownBranchType(runner, valueToWriteToGit)
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

func saveMainBranch(valueToWriteToGit gitdomain.LocalBranchName, valueAlreadyInGit Option[gitdomain.LocalBranchName], runner subshelldomain.Runner) error {
	if existing, hasExisting := valueAlreadyInGit.Get(); hasExisting {
		if existing == valueToWriteToGit {
			return nil
		}
	}
	return gitconfig.SetMainBranch(runner, valueToWriteToGit)
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
