package config

import (
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketcloud"
	"github.com/git-town/git-town/v21/internal/forge/bitbucketdatacenter"
	"github.com/git-town/git-town/v21/internal/forge/codeberg"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gh"
	"github.com/git-town/git-town/v21/internal/forge/gitea"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
	"github.com/git-town/git-town/v21/internal/forge/glab"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
			return executeConfigSetup(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

// the config settings to be used if the user accepts all default options
func defaultUserInput(gitAccess gitconfig.Access, gitVersion git.Version) userInput {
	return userInput{
		config:        config.DefaultUnvalidatedConfig(gitAccess, gitVersion),
		configStorage: dialog.ConfigStorageOptionFile,
	}
}

func executeConfigSetup(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := loadSetupData(repo, verbose)
	if err != nil || exit {
		return err
	}
	tokenScope, forgeTypeOpt, exit, err := enterData(repo, &data)
	if err != nil || exit {
		return err
	}
	err = saveAll(data.userInput, repo.UnvalidatedConfig, data.configFile, tokenScope, forgeTypeOpt, repo.Git, repo.Frontend)
	if err != nil {
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
		Verbose:               verbose,
	})
}

type setupData struct {
	config        config.UnvalidatedConfig
	configFile    Option[configdomain.PartialConfig]
	dialogInputs  components.TestInputs
	localBranches gitdomain.BranchInfos
	remotes       gitdomain.Remotes
	userInput     userInput
}

type userInput struct {
	config        config.UnvalidatedConfig
	configStorage dialog.ConfigStorageOption
}

func determineForgeType(config config.UnvalidatedConfig, userChoice Option[forgedomain.ForgeType]) Option[forgedomain.ForgeType] {
	if userChoice.IsSome() {
		return userChoice
	}
	if devURL, hasDevURL := config.NormalConfig.DevURL().Get(); hasDevURL {
		return forge.Detect(devURL, userChoice)
	}
	return None[forgedomain.ForgeType]()
}

func enterData(repo execute.OpenRepoResult, data *setupData) (configdomain.ConfigScope, Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	tokenScope := configdomain.ConfigScopeLocal
	configFile := data.configFile.GetOrDefault()
	exit, err := dialog.Welcome(data.dialogInputs.Next())
	forgeTypeOpt := None[forgedomain.ForgeType]()
	if err != nil || exit {
		return tokenScope, forgeTypeOpt, exit, err
	}
	data.userInput.config.NormalConfig.Aliases, exit, err = dialog.Aliases(configdomain.AllAliasableCommands(), repo.UnvalidatedConfig.NormalConfig.Aliases, data.dialogInputs.Next())
	if err != nil || exit {
		return tokenScope, forgeTypeOpt, exit, err
	}
	var mainBranch gitdomain.LocalBranchName
	if configFileMainBranch, configFileHasMainBranch := configFile.MainBranch.Get(); configFileHasMainBranch {
		mainBranch = configFileMainBranch
	} else {
		existingMainBranch := repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch
		if existingMainBranch.IsNone() {
			existingMainBranch = repo.Git.DefaultBranch(repo.Backend)
		}
		if existingMainBranch.IsNone() {
			existingMainBranch = repo.Git.OriginHead(repo.Backend)
		}
		mainBranch, exit, err = dialog.MainBranch(data.localBranches.Names(), existingMainBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
		data.userInput.config.UnvalidatedConfig.MainBranch = Some(mainBranch)
	}
	if len(configFile.PerennialBranches) == 0 {
		data.userInput.config.NormalConfig.PerennialBranches, exit, err = dialog.PerennialBranches(data.localBranches.Names(), repo.UnvalidatedConfig.NormalConfig.PerennialBranches, mainBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.PerennialRegex.IsNone() {
		data.userInput.config.NormalConfig.PerennialRegex, exit, err = dialog.PerennialRegex(repo.UnvalidatedConfig.NormalConfig.PerennialRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.FeatureRegex.IsNone() {
		data.userInput.config.NormalConfig.FeatureRegex, exit, err = dialog.FeatureRegex(repo.UnvalidatedConfig.NormalConfig.FeatureRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.UnknownBranchType.IsNone() {
		data.userInput.config.NormalConfig.UnknownBranchType, exit, err = dialog.UnknownBranchType(repo.UnvalidatedConfig.NormalConfig.UnknownBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.DevRemote.IsNone() {
		data.userInput.config.NormalConfig.DevRemote, exit, err = dialog.DevRemote(repo.UnvalidatedConfig.NormalConfig.DevRemote, data.remotes, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	for {
		if configFile.HostingOriginHostname.IsNone() {
			data.userInput.config.NormalConfig.HostingOriginHostname, exit, err = dialog.OriginHostname(repo.UnvalidatedConfig.NormalConfig.HostingOriginHostname, data.dialogInputs.Next())
			if err != nil || exit {
				return tokenScope, forgeTypeOpt, exit, err
			}
		}
		if configFile.ForgeType.IsNone() {
			data.userInput.config.NormalConfig.ForgeType, exit, err = dialog.ForgeType(repo.UnvalidatedConfig.NormalConfig.ForgeType, data.dialogInputs.Next())
			if err != nil || exit {
				return tokenScope, forgeTypeOpt, exit, err
			}
		}
		forgeTypeOpt = determineForgeType(repo.UnvalidatedConfig, data.userInput.config.NormalConfig.ForgeType)
		exit, err = enterForgeAuth(repo, data, forgeTypeOpt)
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
		repeat, exit, err := testForgeAuth(data, repo, forgeTypeOpt)
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
		if !repeat {
			break
		}
	}
	tokenScope, exit, err = enterTokenScope(forgeTypeOpt, data, repo)
	if err != nil || exit {
		return tokenScope, forgeTypeOpt, exit, err
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncFeatureStrategy, exit, err = dialog.SyncFeatureStrategy(repo.UnvalidatedConfig.NormalConfig.SyncFeatureStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncPerennialStrategy, exit, err = dialog.SyncPerennialStrategy(repo.UnvalidatedConfig.NormalConfig.SyncPerennialStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncPrototypeStrategy, exit, err = dialog.SyncPrototypeStrategy(repo.UnvalidatedConfig.NormalConfig.SyncPrototypeStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.SyncUpstream.IsNone() {
		data.userInput.config.NormalConfig.SyncUpstream, exit, err = dialog.SyncUpstream(repo.UnvalidatedConfig.NormalConfig.SyncUpstream, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.SyncTags.IsNone() {
		data.userInput.config.NormalConfig.SyncTags, exit, err = dialog.SyncTags(repo.UnvalidatedConfig.NormalConfig.SyncTags, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.ShareNewBranches.IsNone() {
		data.userInput.config.NormalConfig.ShareNewBranches, exit, err = dialog.ShareNewBranches(repo.UnvalidatedConfig.NormalConfig.ShareNewBranches, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.PushHook.IsNone() {
		data.userInput.config.NormalConfig.PushHook, exit, err = dialog.PushHook(repo.UnvalidatedConfig.NormalConfig.PushHook, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.NewBranchType.IsNone() {
		data.userInput.config.NormalConfig.NewBranchType, exit, err = dialog.NewBranchType(repo.UnvalidatedConfig.NormalConfig.NewBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.ShipStrategy.IsNone() {
		data.userInput.config.NormalConfig.ShipStrategy, exit, err = dialog.ShipStrategy(repo.UnvalidatedConfig.NormalConfig.ShipStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		data.userInput.config.NormalConfig.ShipDeleteTrackingBranch, exit, err = dialog.ShipDeleteTrackingBranch(repo.UnvalidatedConfig.NormalConfig.ShipDeleteTrackingBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeTypeOpt, exit, err
		}
	}
	data.userInput.configStorage, exit, err = dialog.ConfigStorage(data.dialogInputs.Next())
	if err != nil || exit {
		return tokenScope, forgeTypeOpt, exit, err
	}
	return tokenScope, forgeTypeOpt, false, nil
}

func enterForgeAuth(repo execute.OpenRepoResult, data *setupData, forgeTypeOpt Option[forgedomain.ForgeType]) (exit dialogdomain.Exit, err error) {
	forgeType, hasForgeType := forgeTypeOpt.Get()
	if !hasForgeType {
		return false, nil
	}
	switch forgeType {
	case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
		return enterBitbucketToken(data, repo)
	case forgedomain.ForgeTypeCodeberg:
		return enterCodebergToken(data, repo)
	case forgedomain.ForgeTypeGitea:
		return enterGiteaToken(data, repo)
	case forgedomain.ForgeTypeGitHub:
		existing := data.userInput.config.NormalConfig.GitHubConnectorType.Or(repo.UnvalidatedConfig.NormalConfig.GitHubConnectorType)
		answer, exit, err := dialog.GitHubConnectorType(existing, data.dialogInputs.Next())
		if err != nil || exit {
			return exit, err
		}
		data.userInput.config.NormalConfig.GitHubConnectorType = Some(answer)
		switch answer {
		case forgedomain.GitHubConnectorTypeAPI:
			return enterGitHubToken(data, repo)
		case forgedomain.GitHubConnectorTypeGh:
			return false, nil
		}
	case forgedomain.ForgeTypeGitLab:
		existing := data.userInput.config.NormalConfig.GitLabConnectorType.Or(repo.UnvalidatedConfig.NormalConfig.GitLabConnectorType)
		answer, exit, err := dialog.GitLabConnectorType(existing, data.dialogInputs.Next())
		if err != nil || exit {
			return exit, err
		}
		data.userInput.config.NormalConfig.GitLabConnectorType = Some(answer)
		switch answer {
		case forgedomain.GitLabConnectorTypeAPI:
			return enterGitLabToken(data, repo)
		case forgedomain.GitLabConnectorTypeGlab:
			return false, nil
		}
	}
	return false, nil
}

func enterBitbucketToken(data *setupData, repo execute.OpenRepoResult) (exit dialogdomain.Exit, err error) {
	existingUsername := data.userInput.config.NormalConfig.BitbucketUsername.Or(repo.UnvalidatedConfig.NormalConfig.BitbucketUsername)
	data.userInput.config.NormalConfig.BitbucketUsername, exit, err = dialog.BitbucketUsername(existingUsername, data.dialogInputs.Next())
	if err != nil || exit {
		return exit, err
	}
	existingPassword := data.userInput.config.NormalConfig.BitbucketAppPassword.Or(repo.UnvalidatedConfig.NormalConfig.BitbucketAppPassword)
	data.userInput.config.NormalConfig.BitbucketAppPassword, exit, err = dialog.BitbucketAppPassword(existingPassword, data.dialogInputs.Next())
	return exit, err
}

func enterCodebergToken(data *setupData, repo execute.OpenRepoResult) (exit dialogdomain.Exit, err error) {
	existingToken := data.userInput.config.NormalConfig.CodebergToken.Or(repo.UnvalidatedConfig.NormalConfig.CodebergToken)
	data.userInput.config.NormalConfig.CodebergToken, exit, err = dialog.CodebergToken(existingToken, data.dialogInputs.Next())
	return exit, err
}

func enterGiteaToken(data *setupData, repo execute.OpenRepoResult) (exit dialogdomain.Exit, err error) {
	existingToken := data.userInput.config.NormalConfig.GiteaToken.Or(repo.UnvalidatedConfig.NormalConfig.GiteaToken)
	data.userInput.config.NormalConfig.GiteaToken, exit, err = dialog.GiteaToken(existingToken, data.dialogInputs.Next())
	return exit, err
}

func enterGitHubToken(data *setupData, repo execute.OpenRepoResult) (exit dialogdomain.Exit, err error) {
	existingToken := data.userInput.config.NormalConfig.GitHubToken.Or(repo.UnvalidatedConfig.NormalConfig.GitHubToken)
	data.userInput.config.NormalConfig.GitHubToken, exit, err = dialog.GitHubToken(existingToken, data.dialogInputs.Next())
	return exit, err
}

func enterGitLabToken(data *setupData, repo execute.OpenRepoResult) (exit dialogdomain.Exit, err error) {
	existingToken := data.userInput.config.NormalConfig.GitLabToken.Or(repo.UnvalidatedConfig.NormalConfig.GitLabToken)
	data.userInput.config.NormalConfig.GitLabToken, exit, err = dialog.GitLabToken(existingToken, data.dialogInputs.Next())
	return exit, err
}

func testForgeAuth(data *setupData, repo execute.OpenRepoResult, forgeTypeOpt Option[forgedomain.ForgeType]) (repeat bool, exit dialogdomain.Exit, err error) {
	if _, inTest := os.LookupEnv(subshell.TestToken); inTest {
		return false, false, nil
	}
	connectorOpt, err := createConnector(data.userInput.config.NormalConfig, repo.Backend, forgeTypeOpt)
	if err != nil {
		return false, false, err
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return false, false, nil
	}
	verifyResult := connector.VerifyConnection()
	if verifyResult.AuthenticationError != nil {
		return dialog.CredentialsNoAccess(verifyResult.AuthenticationError, data.dialogInputs.Next())
	}
	if user, hasUser := verifyResult.AuthenticatedUser.Get(); hasUser {
		fmt.Printf(messages.CredentialsForgeUserName, components.FormattedSelection(user, exit))
	}
	if verifyResult.AuthorizationError != nil {
		return dialog.CredentialsNoProposalAccess(verifyResult.AuthorizationError, data.dialogInputs.Next())
	}
	fmt.Println(messages.CredentialsAccess)
	return false, false, nil
}

func createConnector(userInput config.NormalConfig, backend subshelldomain.RunnerQuerier, forgeTypeOpt Option[forgedomain.ForgeType]) (Option[forgedomain.Connector], error) {
	forgeType, hasForgeType := forgeTypeOpt.Get()
	if !hasForgeType {
		return None[forgedomain.Connector](), nil
	}
	var connector forgedomain.Connector
	var err error
	switch forgeType {
	case forgedomain.ForgeTypeBitbucket:
		connector = bitbucketcloud.NewConnector(bitbucketcloud.NewConnectorArgs{
			AppPassword: userInput.BitbucketAppPassword,
			ForgeType:   Some(forgedomain.ForgeTypeBitbucket),
			Log:         print.Logger{},
			RemoteURL:   userInput.DevURL().GetOrDefault(),
			UserName:    userInput.BitbucketUsername,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeBitbucketDatacenter:
		connector = bitbucketdatacenter.NewConnector(bitbucketdatacenter.NewConnectorArgs{
			AppPassword: userInput.BitbucketAppPassword,
			ForgeType:   Some(forgedomain.ForgeTypeBitbucketDatacenter),
			Log:         print.Logger{},
			RemoteURL:   userInput.DevURL().GetOrDefault(),
			UserName:    userInput.BitbucketUsername,
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeCodeberg:
		if subshell.IsInTest() {
			return None[forgedomain.Connector](), nil
		}
		connector, err = codeberg.NewConnector(codeberg.NewConnectorArgs{
			APIToken:  userInput.CodebergToken,
			Log:       print.Logger{},
			RemoteURL: userInput.DevURL().GetOrDefault(),
		})
		return Some(connector), err
	case forgedomain.ForgeTypeGitea:
		if subshell.IsInTest() {
			return None[forgedomain.Connector](), nil
		}
		connector = gitea.NewConnector(gitea.NewConnectorArgs{
			APIToken:  userInput.GiteaToken,
			Log:       print.Logger{},
			RemoteURL: userInput.DevURL().GetOrDefault(),
		})
		return Some(connector), nil
	case forgedomain.ForgeTypeGitHub:
		if connectorType, hasConnectorType := userInput.GitHubConnectorType.Get(); hasConnectorType {
			switch connectorType {
			case forgedomain.GitHubConnectorTypeAPI:
				connector, err = github.NewConnector(github.NewConnectorArgs{
					APIToken:  userInput.GitHubToken,
					Log:       print.Logger{},
					RemoteURL: userInput.DevURL().GetOrDefault(),
				})
				return Some(connector), err
			case forgedomain.GitHubConnectorTypeGh:
				connector = gh.Connector{
					Backend:  backend,
					Frontend: backend,
				}
				return Some(connector), nil
			}
		}
	case forgedomain.ForgeTypeGitLab:
		if connectorType, hasConnectorType := userInput.GitLabConnectorType.Get(); hasConnectorType {
			switch connectorType {
			case forgedomain.GitLabConnectorTypeAPI:
				connector, err = gitlab.NewConnector(gitlab.NewConnectorArgs{
					APIToken:  userInput.GitLabToken,
					Log:       print.Logger{},
					RemoteURL: userInput.DevURL().GetOrDefault(),
				})
				return Some(connector), err
			case forgedomain.GitLabConnectorTypeGlab:
				connector = glab.Connector{
					Backend:  backend,
					Frontend: backend,
				}
				return Some(connector), nil
			}
		}
	}
	return None[forgedomain.Connector](), nil
}

func enterTokenScope(forgeTypeOpt Option[forgedomain.ForgeType], data *setupData, repo execute.OpenRepoResult) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if shouldAskForScope(forgeTypeOpt, data, repo) {
		return tokenScopeDialog(forgeTypeOpt, data, repo)
	}
	return configdomain.ConfigScopeLocal, false, nil
}

func shouldAskForScope(forgeTypeOpt Option[forgedomain.ForgeType], data *setupData, repo execute.OpenRepoResult) bool {
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			return existsAndChanged(data.userInput.config.NormalConfig.BitbucketUsername, repo.UnvalidatedConfig.NormalConfig.BitbucketUsername) &&
				existsAndChanged(data.userInput.config.NormalConfig.BitbucketAppPassword, repo.UnvalidatedConfig.NormalConfig.BitbucketAppPassword)
		case forgedomain.ForgeTypeCodeberg:
			return existsAndChanged(data.userInput.config.NormalConfig.CodebergToken, repo.UnvalidatedConfig.NormalConfig.CodebergToken)
		case forgedomain.ForgeTypeGitea:
			return existsAndChanged(data.userInput.config.NormalConfig.GiteaToken, repo.UnvalidatedConfig.NormalConfig.GiteaToken)
		case forgedomain.ForgeTypeGitHub:
			return existsAndChanged(data.userInput.config.NormalConfig.GitHubToken, repo.UnvalidatedConfig.NormalConfig.GitHubToken)
		case forgedomain.ForgeTypeGitLab:
			return existsAndChanged(data.userInput.config.NormalConfig.GitLabToken, repo.UnvalidatedConfig.NormalConfig.GitLabToken)
		}
	}
	return false
}

func tokenScopeDialog(forgeTypeOpt Option[forgedomain.ForgeType], data *setupData, repo execute.OpenRepoResult) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			existingScope := determineScope(repo.ConfigSnapshot, configdomain.KeyBitbucketUsername, repo.UnvalidatedConfig.NormalConfig.BitbucketUsername)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeCodeberg:
			existingScope := determineScope(repo.ConfigSnapshot, configdomain.KeyCodebergToken, repo.UnvalidatedConfig.NormalConfig.CodebergToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeGitea:
			existingScope := determineScope(repo.ConfigSnapshot, configdomain.KeyGiteaToken, repo.UnvalidatedConfig.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeGitHub:
			existingScope := determineScope(repo.ConfigSnapshot, configdomain.KeyGitHubToken, repo.UnvalidatedConfig.NormalConfig.GitHubToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeGitLab:
			existingScope := determineScope(repo.ConfigSnapshot, configdomain.KeyGitLabToken, repo.UnvalidatedConfig.NormalConfig.GitLabToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		}
	}
	return configdomain.ConfigScopeLocal, false, nil
}

func determineScope(configSnapshot undoconfig.ConfigSnapshot, key configdomain.Key, oldValue fmt.Stringer) configdomain.ConfigScope {
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

func loadSetupData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data setupData, exit dialogdomain.Exit, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
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
		Verbose:               verbose,
	})
	if err != nil {
		return data, exit, err
	}
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, exit, err
	}
	if len(remotes) == 0 {
		remotes = gitdomain.Remotes{repo.Git.DefaultRemote(repo.Backend)}
	}
	return setupData{
		config:        repo.UnvalidatedConfig,
		configFile:    repo.UnvalidatedConfig.NormalConfig.ConfigFile,
		dialogInputs:  dialogTestInputs,
		localBranches: branchesSnapshot.Branches,
		remotes:       remotes,
		userInput:     defaultUserInput(repo.UnvalidatedConfig.NormalConfig.GitConfigAccess, repo.UnvalidatedConfig.NormalConfig.GitVersion),
	}, exit, nil
}

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, configFile Option[configdomain.PartialConfig], tokenScope configdomain.ConfigScope, forgeTypeOpt Option[forgedomain.ForgeType], gitCommands git.Commands, frontend subshelldomain.Runner) error {
	err := saveAliases(oldConfig.NormalConfig.Aliases, userInput.config.NormalConfig.Aliases, gitCommands, frontend)
	if err != nil {
		return err
	}
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			err = saveBitbucketUsername(oldConfig.NormalConfig.BitbucketUsername, userInput.config.NormalConfig.BitbucketUsername, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
			err = saveBitbucketAppPassword(oldConfig.NormalConfig.BitbucketAppPassword, userInput.config.NormalConfig.BitbucketAppPassword, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeCodeberg:
			err = saveCodebergToken(oldConfig.NormalConfig.CodebergToken, userInput.config.NormalConfig.CodebergToken, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeGitHub:
			err = saveGitHubToken(oldConfig.NormalConfig.GitHubToken, userInput.config.NormalConfig.GitHubToken, tokenScope, userInput.config.NormalConfig.GitHubConnectorType, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeGitLab:
			err = saveGitLabToken(oldConfig.NormalConfig.GitLabToken, userInput.config.NormalConfig.GitLabToken, tokenScope, userInput.config.NormalConfig.GitLabConnectorType, gitCommands, frontend)
			if err != nil {
				return err
			}
		case forgedomain.ForgeTypeGitea:
			err = saveGiteaToken(oldConfig.NormalConfig.GiteaToken, userInput.config.NormalConfig.GiteaToken, tokenScope, gitCommands, frontend)
			if err != nil {
				return err
			}
		}
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
	fc := execute.FailureCollector{}
	if configFile.NewBranchType.IsNone() {
		fc.Check(saveNewBranchType(oldConfig.NormalConfig.NewBranchType, userInput.config.NormalConfig.NewBranchType, oldConfig))
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(saveForgeType(oldConfig.NormalConfig.ForgeType, userInput.config.NormalConfig.ForgeType, gitCommands, frontend))
	}
	if configFile.GitHubConnectorType.IsNone() {
		fc.Check(saveGitHubConnectorType(oldConfig.NormalConfig.GitHubConnectorType, userInput.config.NormalConfig.GitHubConnectorType, gitCommands, frontend))
	}
	if configFile.GitLabConnectorType.IsNone() {
		fc.Check(saveGitLabConnectorType(oldConfig.NormalConfig.GitLabConnectorType, userInput.config.NormalConfig.GitLabConnectorType, gitCommands, frontend))
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(saveOriginHostname(oldConfig.NormalConfig.HostingOriginHostname, userInput.config.NormalConfig.HostingOriginHostname, gitCommands, frontend))
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(saveMainBranch(oldConfig.UnvalidatedConfig.MainBranch, userInput.config.UnvalidatedConfig.MainBranch, oldConfig))
	}
	if len(configFile.PerennialBranches) == 0 {
		fc.Check(savePerennialBranches(oldConfig.NormalConfig.PerennialBranches, userInput.config.NormalConfig.PerennialBranches, oldConfig))
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(savePerennialRegex(oldConfig.NormalConfig.PerennialRegex, userInput.config.NormalConfig.PerennialRegex, oldConfig))
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(saveUnknownBranchType(oldConfig.NormalConfig.UnknownBranchType, userInput.config.NormalConfig.UnknownBranchType, oldConfig))
	}
	if configFile.DevRemote.IsNone() {
		fc.Check(saveDevRemote(oldConfig.NormalConfig.DevRemote, userInput.config.NormalConfig.DevRemote, oldConfig))
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(saveFeatureRegex(oldConfig.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, oldConfig))
	}
	if configFile.PushHook.IsNone() {
		fc.Check(savePushHook(oldConfig.NormalConfig.PushHook, userInput.config.NormalConfig.PushHook, oldConfig))
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(saveShareNewBranches(oldConfig.NormalConfig.ShareNewBranches, userInput.config.NormalConfig.ShareNewBranches, oldConfig))
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(saveShipStrategy(oldConfig.NormalConfig.ShipStrategy, userInput.config.NormalConfig.ShipStrategy, oldConfig))
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(saveShipDeleteTrackingBranch(oldConfig.NormalConfig.ShipDeleteTrackingBranch, userInput.config.NormalConfig.ShipDeleteTrackingBranch, oldConfig))
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(saveSyncFeatureStrategy(oldConfig.NormalConfig.SyncFeatureStrategy, userInput.config.NormalConfig.SyncFeatureStrategy, oldConfig))
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(saveSyncPerennialStrategy(oldConfig.NormalConfig.SyncPerennialStrategy, userInput.config.NormalConfig.SyncPerennialStrategy, oldConfig))
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(saveSyncPrototypeStrategy(oldConfig.NormalConfig.SyncPrototypeStrategy, userInput.config.NormalConfig.SyncPrototypeStrategy, oldConfig))
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(saveSyncUpstream(oldConfig.NormalConfig.SyncUpstream, userInput.config.NormalConfig.SyncUpstream, oldConfig))
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(saveSyncTags(oldConfig.NormalConfig.SyncTags, userInput.config.NormalConfig.SyncTags, oldConfig))
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
	err := configfile.Save(&userInput.config)
	if err != nil {
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
	err = saveUnknownBranchType(config.NormalConfig.UnknownBranchType, userInput.config.NormalConfig.UnknownBranchType, config)
	if err != nil {
		return err
	}
	return saveFeatureRegex(config.NormalConfig.FeatureRegex, userInput.config.NormalConfig.FeatureRegex, config)
}
