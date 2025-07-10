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
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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

// the config settings to be used if the user accepts all default options
func defaultUserInput(gitVersion git.Version) userInput {
	return userInput{
		config:        config.DefaultUnvalidatedConfig(gitVersion),
		configStorage: dialog.ConfigStorageOptionFile,
	}
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
	tokenScope, forgeTypeOpt, exit, err := enterData(repo, &data)
	if err != nil || exit {
		return err
	}
	if err = saveAll(data.userInput, repo.UnvalidatedConfig, data.configFile, tokenScope, forgeTypeOpt, repo.Frontend); err != nil {
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
	config        config.UnvalidatedConfig
	configFile    Option[configdomain.PartialConfig]
	dialogInputs  dialogcomponents.TestInputs
	localBranches gitdomain.BranchInfos
	remotes       gitdomain.Remotes
	userInput     userInput
}

type userInput struct {
	config        config.UnvalidatedConfig
	configStorage dialog.ConfigStorageOption
}

func determineForgeType(config config.UnvalidatedConfig, userChoice Option[forgedomain.ForgeType], querier subshelldomain.Querier) Option[forgedomain.ForgeType] {
	if userChoice.IsSome() {
		return userChoice
	}
	if devURL, hasDevURL := config.NormalConfig.DevURL(querier).Get(); hasDevURL {
		return forge.Detect(devURL, userChoice)
	}
	return None[forgedomain.ForgeType]()
}

func enterData(repo execute.OpenRepoResult, data *setupData) (configdomain.ConfigScope, Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	tokenScope := configdomain.ConfigScopeLocal
	configFile := data.configFile.GetOrDefault()
	exit, err := dialog.Welcome(data.dialogInputs.Next())
	forgeType := None[forgedomain.ForgeType]()
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	data.userInput.config.NormalConfig.Aliases, exit, err = dialog.Aliases(configdomain.AllAliasableCommands(), repo.UnvalidatedConfig.NormalConfig.Aliases, data.dialogInputs.Next())
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
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
			return tokenScope, forgeType, exit, err
		}
		data.userInput.config.UnvalidatedConfig.MainBranch = Some(mainBranch)
	}
	if len(configFile.PerennialBranches) == 0 {
		data.userInput.config.NormalConfig.PerennialBranches, exit, err = dialog.PerennialBranches(data.localBranches.Names(), repo.UnvalidatedConfig.NormalConfig.PerennialBranches, mainBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.PerennialRegex.IsNone() {
		data.userInput.config.NormalConfig.PerennialRegex, exit, err = dialog.PerennialRegex(repo.UnvalidatedConfig.NormalConfig.PerennialRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.FeatureRegex.IsNone() {
		data.userInput.config.NormalConfig.FeatureRegex, exit, err = dialog.FeatureRegex(repo.UnvalidatedConfig.NormalConfig.FeatureRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ContributionRegex.IsNone() {
		data.userInput.config.NormalConfig.ContributionRegex, exit, err = dialog.ContributionRegex(repo.UnvalidatedConfig.NormalConfig.ContributionRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ObservedRegex.IsNone() {
		data.userInput.config.NormalConfig.ObservedRegex, exit, err = dialog.ObservedRegex(repo.UnvalidatedConfig.NormalConfig.ObservedRegex, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.UnknownBranchType.IsNone() {
		data.userInput.config.NormalConfig.UnknownBranchType, exit, err = dialog.UnknownBranchType(repo.UnvalidatedConfig.NormalConfig.UnknownBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.DevRemote.IsNone() {
		data.userInput.config.NormalConfig.DevRemote, exit, err = dialog.DevRemote(repo.UnvalidatedConfig.NormalConfig.DevRemote, data.remotes, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	for {
		if configFile.HostingOriginHostname.IsNone() {
			data.userInput.config.NormalConfig.HostingOriginHostname, exit, err = dialog.OriginHostname(repo.UnvalidatedConfig.NormalConfig.HostingOriginHostname, data.dialogInputs.Next())
			if err != nil || exit {
				return tokenScope, forgeType, exit, err
			}
		}
		if configFile.ForgeType.IsNone() {
			data.userInput.config.NormalConfig.ForgeType, exit, err = dialog.ForgeType(repo.UnvalidatedConfig.NormalConfig.ForgeType, data.dialogInputs.Next())
			if err != nil || exit {
				return tokenScope, forgeType, exit, err
			}
		}
		forgeType = determineForgeType(repo.UnvalidatedConfig, data.userInput.config.NormalConfig.ForgeType, repo.Backend)
		exit, err = enterForgeAuth(repo, data, forgeType)
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
		repeat, exit, err := testForgeAuth(data, repo, forgeType)
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
		if !repeat {
			break
		}
	}
	tokenScope, exit, err = enterTokenScope(forgeType, data, repo)
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncFeatureStrategy, exit, err = dialog.SyncFeatureStrategy(repo.UnvalidatedConfig.NormalConfig.SyncFeatureStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncPerennialStrategy, exit, err = dialog.SyncPerennialStrategy(repo.UnvalidatedConfig.NormalConfig.SyncPerennialStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		data.userInput.config.NormalConfig.SyncPrototypeStrategy, exit, err = dialog.SyncPrototypeStrategy(repo.UnvalidatedConfig.NormalConfig.SyncPrototypeStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncUpstream.IsNone() {
		data.userInput.config.NormalConfig.SyncUpstream, exit, err = dialog.SyncUpstream(repo.UnvalidatedConfig.NormalConfig.SyncUpstream, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncTags.IsNone() {
		data.userInput.config.NormalConfig.SyncTags, exit, err = dialog.SyncTags(repo.UnvalidatedConfig.NormalConfig.SyncTags, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ShareNewBranches.IsNone() {
		data.userInput.config.NormalConfig.ShareNewBranches, exit, err = dialog.ShareNewBranches(repo.UnvalidatedConfig.NormalConfig.ShareNewBranches, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.PushHook.IsNone() {
		data.userInput.config.NormalConfig.PushHook, exit, err = dialog.PushHook(repo.UnvalidatedConfig.NormalConfig.PushHook, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.NewBranchType.IsNone() {
		data.userInput.config.NormalConfig.NewBranchType, exit, err = dialog.NewBranchType(repo.UnvalidatedConfig.NormalConfig.NewBranchType, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ShipStrategy.IsNone() {
		data.userInput.config.NormalConfig.ShipStrategy, exit, err = dialog.ShipStrategy(repo.UnvalidatedConfig.NormalConfig.ShipStrategy, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		data.userInput.config.NormalConfig.ShipDeleteTrackingBranch, exit, err = dialog.ShipDeleteTrackingBranch(repo.UnvalidatedConfig.NormalConfig.ShipDeleteTrackingBranch, data.dialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	data.userInput.configStorage, exit, err = dialog.ConfigStorage(data.dialogInputs.Next())
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	return tokenScope, forgeType, false, nil
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
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: data.userInput.config.NormalConfig.BitbucketAppPassword.Or(data.config.NormalConfig.BitbucketAppPassword),
		BitbucketUsername:    data.userInput.config.NormalConfig.BitbucketUsername.Or(data.config.NormalConfig.BitbucketUsername),
		CodebergToken:        data.userInput.config.NormalConfig.CodebergToken.Or(data.config.NormalConfig.CodebergToken),
		ForgeType:            forgeTypeOpt,
		Frontend:             repo.Backend,
		GitHubConnectorType:  data.userInput.config.NormalConfig.GitHubConnectorType.Or(data.config.NormalConfig.GitHubConnectorType),
		GitHubToken:          data.userInput.config.NormalConfig.GitHubToken.Or(data.config.NormalConfig.GitHubToken),
		GitLabConnectorType:  data.userInput.config.NormalConfig.GitLabConnectorType.Or(data.config.NormalConfig.GitLabConnectorType),
		GitLabToken:          data.userInput.config.NormalConfig.GitLabToken.Or(data.config.NormalConfig.GitLabToken),
		GiteaToken:           data.userInput.config.NormalConfig.GiteaToken.Or(data.config.NormalConfig.GiteaToken),
		Log:                  print.Logger{},
		RemoteURL:            data.userInput.config.NormalConfig.DevURL(repo.Backend).Or(data.config.NormalConfig.DevURL(repo.Backend)),
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
		return dialog.CredentialsNoAccess(verifyResult.AuthenticationError, data.dialogInputs.Next())
	}
	if user, hasUser := verifyResult.AuthenticatedUser.Get(); hasUser {
		fmt.Printf(messages.CredentialsForgeUserName, dialogcomponents.FormattedSelection(user, exit))
	}
	if verifyResult.AuthorizationError != nil {
		return dialog.CredentialsNoProposalAccess(verifyResult.AuthorizationError, data.dialogInputs.Next())
	}
	fmt.Println(messages.CredentialsAccess)
	return false, false, nil
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
			existingScope := determineExistingScope(repo.ConfigSnapshot, configdomain.KeyBitbucketUsername, repo.UnvalidatedConfig.NormalConfig.BitbucketUsername)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeCodeberg:
			existingScope := determineExistingScope(repo.ConfigSnapshot, configdomain.KeyCodebergToken, repo.UnvalidatedConfig.NormalConfig.CodebergToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeGitea:
			existingScope := determineExistingScope(repo.ConfigSnapshot, configdomain.KeyGiteaToken, repo.UnvalidatedConfig.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeGitHub:
			existingScope := determineExistingScope(repo.ConfigSnapshot, configdomain.KeyGitHubToken, repo.UnvalidatedConfig.NormalConfig.GitHubToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
		case forgedomain.ForgeTypeGitLab:
			existingScope := determineExistingScope(repo.ConfigSnapshot, configdomain.KeyGitLabToken, repo.UnvalidatedConfig.NormalConfig.GitLabToken)
			return dialog.TokenScope(existingScope, data.dialogInputs.Next())
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
		config:        repo.UnvalidatedConfig,
		configFile:    repo.UnvalidatedConfig.NormalConfig.File,
		dialogInputs:  dialogTestInputs,
		localBranches: branchesSnapshot.Branches,
		remotes:       remotes,
		userInput:     defaultUserInput(repo.UnvalidatedConfig.NormalConfig.GitVersion),
	}, exit, nil
}

func saveAll(userInput userInput, oldConfig config.UnvalidatedConfig, configFile Option[configdomain.PartialConfig], tokenScope configdomain.ConfigScope, forgeTypeOpt Option[forgedomain.ForgeType], frontend subshelldomain.Runner) error {
	fc := gohacks.ErrorCollector{}
	fc.Check(
		saveAliases(userInput.config.NormalConfig.Aliases, oldConfig.NormalConfig, frontend),
	)
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			fc.Check(
				saveBitbucketUsername(userInput.config.NormalConfig.BitbucketUsername, oldConfig.NormalConfig, tokenScope, frontend),
			)
			fc.Check(
				saveBitbucketAppPassword(userInput.config.NormalConfig.BitbucketAppPassword, oldConfig.NormalConfig, tokenScope, frontend),
			)
		case forgedomain.ForgeTypeCodeberg:
			fc.Check(
				saveCodebergToken(userInput.config.NormalConfig.CodebergToken, oldConfig.NormalConfig, tokenScope, frontend),
			)
		case forgedomain.ForgeTypeGitHub:
			fc.Check(
				saveGitHubToken(userInput.config.NormalConfig.GitHubToken, oldConfig.NormalConfig, tokenScope, userInput.config.NormalConfig.GitHubConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitLab:
			fc.Check(
				saveGitLabToken(userInput.config.NormalConfig.GitLabToken, oldConfig.NormalConfig, tokenScope, userInput.config.NormalConfig.GitLabConnectorType, frontend),
			)
		case forgedomain.ForgeTypeGitea:
			fc.Check(
				saveGiteaToken(userInput.config.NormalConfig.GiteaToken, oldConfig.NormalConfig, tokenScope, frontend),
			)
		}
	}
	if fc.Err != nil {
		return fc.Err
	}
	switch userInput.configStorage {
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
			saveNewBranchType(userInput.config.NormalConfig.NewBranchType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ForgeType.IsNone() {
		fc.Check(
			saveForgeType(userInput.config.NormalConfig.ForgeType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.GitHubConnectorType.IsNone() {
		fc.Check(
			saveGitHubConnectorType(userInput.config.NormalConfig.GitHubConnectorType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.GitLabConnectorType.IsNone() {
		fc.Check(
			saveGitLabConnectorType(userInput.config.NormalConfig.GitLabConnectorType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.HostingOriginHostname.IsNone() {
		fc.Check(
			saveOriginHostname(userInput.config.NormalConfig.HostingOriginHostname, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.MainBranch.IsNone() {
		fc.Check(
			saveMainBranch(userInput.config.UnvalidatedConfig.MainBranch, oldConfig, frontend),
		)
	}
	if len(configFile.PerennialBranches) == 0 {
		fc.Check(
			savePerennialBranches(userInput.config.NormalConfig.PerennialBranches, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.PerennialRegex.IsNone() {
		fc.Check(
			savePerennialRegex(userInput.config.NormalConfig.PerennialRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.UnknownBranchType.IsNone() {
		fc.Check(
			saveUnknownBranchType(userInput.config.NormalConfig.UnknownBranchType, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.DevRemote.IsNone() {
		fc.Check(
			saveDevRemote(userInput.config.NormalConfig.DevRemote, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.FeatureRegex.IsNone() {
		fc.Check(
			saveFeatureRegex(userInput.config.NormalConfig.FeatureRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ContributionRegex.IsNone() {
		fc.Check(
			saveContributionRegex(userInput.config.NormalConfig.ContributionRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ObservedRegex.IsNone() {
		fc.Check(
			saveObservedRegex(userInput.config.NormalConfig.ObservedRegex, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.PushHook.IsNone() {
		fc.Check(
			savePushHook(userInput.config.NormalConfig.PushHook, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ShareNewBranches.IsNone() {
		fc.Check(
			saveShareNewBranches(userInput.config.NormalConfig.ShareNewBranches, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ShipStrategy.IsNone() {
		fc.Check(
			saveShipStrategy(userInput.config.NormalConfig.ShipStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		fc.Check(
			saveShipDeleteTrackingBranch(userInput.config.NormalConfig.ShipDeleteTrackingBranch, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		fc.Check(
			saveSyncFeatureStrategy(userInput.config.NormalConfig.SyncFeatureStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		fc.Check(
			saveSyncPerennialStrategy(userInput.config.NormalConfig.SyncPerennialStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		fc.Check(
			saveSyncPrototypeStrategy(userInput.config.NormalConfig.SyncPrototypeStrategy, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncUpstream.IsNone() {
		fc.Check(
			saveSyncUpstream(userInput.config.NormalConfig.SyncUpstream, oldConfig.NormalConfig, frontend),
		)
	}
	if configFile.SyncTags.IsNone() {
		fc.Check(
			saveSyncTags(userInput.config.NormalConfig.SyncTags, oldConfig.NormalConfig, frontend),
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

func saveMainBranch(newValue Option[gitdomain.LocalBranchName], config config.UnvalidatedConfig, runner subshelldomain.Runner) error {
	if newValue.Equal(config.UnvalidatedConfig.MainBranch) {
		return nil
	}
	if mainBranch, hasNewValue := newValue.Get(); hasNewValue {
		return config.SetMainBranch(mainBranch, runner)
	}
	return nil
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
	if err := configfile.Save(&userInput.config); err != nil {
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
	if err := saveUnknownBranchType(userInput.config.NormalConfig.UnknownBranchType, config, runner); err != nil {
		return err
	}
	return saveFeatureRegex(userInput.config.NormalConfig.FeatureRegex, config, runner)
}
