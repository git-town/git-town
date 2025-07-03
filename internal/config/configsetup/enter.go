package configsetup

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func enterData(data *SetupData) (configdomain.ConfigScope, Option[forgedomain.ForgeType], dialogdomain.Exit, error) {
	tokenScope := configdomain.ConfigScopeLocal
	configFile := data.ConfigFile.GetOrDefault()
	exit, err := dialog.Welcome(data.DialogInputs.Next())
	forgeType := None[forgedomain.ForgeType]()
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	data.UserInput.Config.NormalConfig.Aliases, exit, err = dialog.Aliases(configdomain.AllAliasableCommands(), data.UnvalidatedConfig.NormalConfig.Aliases, data.DialogInputs.Next())
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	var mainBranch gitdomain.LocalBranchName
	if configFileMainBranch, configFileHasMainBranch := configFile.MainBranch.Get(); configFileHasMainBranch {
		mainBranch = configFileMainBranch
	} else {
		existingMainBranch := data.UnvalidatedConfig.UnvalidatedConfig.MainBranch
		if existingMainBranch.IsNone() {
			existingMainBranch = data.Git.DefaultBranch(data.Backend)
		}
		if existingMainBranch.IsNone() {
			existingMainBranch = data.Git.OriginHead(data.Backend)
		}
		mainBranch, exit, err = dialog.MainBranch(data.LocalBranches.Names(), existingMainBranch, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
		data.UserInput.Config.UnvalidatedConfig.MainBranch = Some(mainBranch)
	}
	if len(configFile.PerennialBranches) == 0 {
		data.UserInput.Config.NormalConfig.PerennialBranches, exit, err = dialog.PerennialBranches(data.LocalBranches.Names(), data.UnvalidatedConfig.NormalConfig.PerennialBranches, mainBranch, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.PerennialRegex.IsNone() {
		data.UserInput.Config.NormalConfig.PerennialRegex, exit, err = dialog.PerennialRegex(data.UnvalidatedConfig.NormalConfig.PerennialRegex, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.FeatureRegex.IsNone() {
		data.UserInput.Config.NormalConfig.FeatureRegex, exit, err = dialog.FeatureRegex(data.UnvalidatedConfig.NormalConfig.FeatureRegex, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.UnknownBranchType.IsNone() {
		data.UserInput.Config.NormalConfig.UnknownBranchType, exit, err = dialog.UnknownBranchType(data.UnvalidatedConfig.NormalConfig.UnknownBranchType, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.DevRemote.IsNone() {
		data.UserInput.Config.NormalConfig.DevRemote, exit, err = dialog.DevRemote(data.UnvalidatedConfig.NormalConfig.DevRemote, data.Remotes, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	for {
		if configFile.HostingOriginHostname.IsNone() {
			data.UserInput.Config.NormalConfig.HostingOriginHostname, exit, err = dialog.OriginHostname(data.UnvalidatedConfig.NormalConfig.HostingOriginHostname, data.DialogInputs.Next())
			if err != nil || exit {
				return tokenScope, forgeType, exit, err
			}
		}
		if configFile.ForgeType.IsNone() {
			data.UserInput.Config.NormalConfig.ForgeType, exit, err = dialog.ForgeType(data.UnvalidatedConfig.NormalConfig.ForgeType, data.DialogInputs.Next())
			if err != nil || exit {
				return tokenScope, forgeType, exit, err
			}
		}
		forgeType = determineForgeType(data.UnvalidatedConfig, data.UserInput.Config.NormalConfig.ForgeType)
		exit, err = enterForgeAuth(data, forgeType)
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
		repeat, exit, err := testForgeAuth(data, forgeType)
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
		if !repeat {
			break
		}
	}
	tokenScope, exit, err = enterTokenScope(forgeType, data)
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	if configFile.SyncFeatureStrategy.IsNone() {
		data.UserInput.Config.NormalConfig.SyncFeatureStrategy, exit, err = dialog.SyncFeatureStrategy(data.UnvalidatedConfig.NormalConfig.SyncFeatureStrategy, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncPerennialStrategy.IsNone() {
		data.UserInput.Config.NormalConfig.SyncPerennialStrategy, exit, err = dialog.SyncPerennialStrategy(data.UnvalidatedConfig.NormalConfig.SyncPerennialStrategy, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncPrototypeStrategy.IsNone() {
		data.UserInput.Config.NormalConfig.SyncPrototypeStrategy, exit, err = dialog.SyncPrototypeStrategy(data.UnvalidatedConfig.NormalConfig.SyncPrototypeStrategy, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncUpstream.IsNone() {
		data.UserInput.Config.NormalConfig.SyncUpstream, exit, err = dialog.SyncUpstream(data.UnvalidatedConfig.NormalConfig.SyncUpstream, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.SyncTags.IsNone() {
		data.UserInput.Config.NormalConfig.SyncTags, exit, err = dialog.SyncTags(data.UnvalidatedConfig.NormalConfig.SyncTags, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ShareNewBranches.IsNone() {
		data.UserInput.Config.NormalConfig.ShareNewBranches, exit, err = dialog.ShareNewBranches(data.UnvalidatedConfig.NormalConfig.ShareNewBranches, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.PushHook.IsNone() {
		data.UserInput.Config.NormalConfig.PushHook, exit, err = dialog.PushHook(data.UnvalidatedConfig.NormalConfig.PushHook, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.NewBranchType.IsNone() {
		data.UserInput.Config.NormalConfig.NewBranchType, exit, err = dialog.NewBranchType(data.UnvalidatedConfig.NormalConfig.NewBranchType, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ShipStrategy.IsNone() {
		data.UserInput.Config.NormalConfig.ShipStrategy, exit, err = dialog.ShipStrategy(data.UnvalidatedConfig.NormalConfig.ShipStrategy, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	if configFile.ShipDeleteTrackingBranch.IsNone() {
		data.UserInput.Config.NormalConfig.ShipDeleteTrackingBranch, exit, err = dialog.ShipDeleteTrackingBranch(data.UnvalidatedConfig.NormalConfig.ShipDeleteTrackingBranch, data.DialogInputs.Next())
		if err != nil || exit {
			return tokenScope, forgeType, exit, err
		}
	}
	data.UserInput.ConfigStorage, exit, err = dialog.ConfigStorage(data.DialogInputs.Next())
	if err != nil || exit {
		return tokenScope, forgeType, exit, err
	}
	return tokenScope, forgeType, false, nil
}

func enterForgeAuth(data *SetupData, forgeTypeOpt Option[forgedomain.ForgeType]) (exit dialogdomain.Exit, err error) {
	forgeType, hasForgeType := forgeTypeOpt.Get()
	if !hasForgeType {
		return false, nil
	}
	switch forgeType {
	case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
		return enterBitbucketToken(data)
	case forgedomain.ForgeTypeCodeberg:
		return enterCodebergToken(data)
	case forgedomain.ForgeTypeGitea:
		return enterGiteaToken(data)
	case forgedomain.ForgeTypeGitHub:
		existing := data.UserInput.Config.NormalConfig.GitHubConnectorType.Or(data.UnvalidatedConfig.NormalConfig.GitHubConnectorType)
		answer, exit, err := dialog.GitHubConnectorType(existing, data.DialogInputs.Next())
		if err != nil || exit {
			return exit, err
		}
		data.UserInput.Config.NormalConfig.GitHubConnectorType = Some(answer)
		switch answer {
		case forgedomain.GitHubConnectorTypeAPI:
			return enterGitHubToken(data)
		case forgedomain.GitHubConnectorTypeGh:
			return false, nil
		}
	case forgedomain.ForgeTypeGitLab:
		existing := data.UserInput.Config.NormalConfig.GitLabConnectorType.Or(data.UnvalidatedConfig.NormalConfig.GitLabConnectorType)
		answer, exit, err := dialog.GitLabConnectorType(existing, data.DialogInputs.Next())
		if err != nil || exit {
			return exit, err
		}
		data.UserInput.Config.NormalConfig.GitLabConnectorType = Some(answer)
		switch answer {
		case forgedomain.GitLabConnectorTypeAPI:
			return enterGitLabToken(data)
		case forgedomain.GitLabConnectorTypeGlab:
			return false, nil
		}
	}
	return false, nil
}

func enterBitbucketToken(data *SetupData) (exit dialogdomain.Exit, err error) {
	existingUsername := data.UserInput.Config.NormalConfig.BitbucketUsername.Or(data.UnvalidatedConfig.NormalConfig.BitbucketUsername)
	data.UserInput.Config.NormalConfig.BitbucketUsername, exit, err = dialog.BitbucketUsername(existingUsername, data.DialogInputs.Next())
	if err != nil || exit {
		return exit, err
	}
	existingPassword := data.UserInput.Config.NormalConfig.BitbucketAppPassword.Or(data.UnvalidatedConfig.NormalConfig.BitbucketAppPassword)
	data.UserInput.Config.NormalConfig.BitbucketAppPassword, exit, err = dialog.BitbucketAppPassword(existingPassword, data.DialogInputs.Next())
	return exit, err
}

func enterCodebergToken(data *SetupData) (exit dialogdomain.Exit, err error) {
	existingToken := data.UserInput.Config.NormalConfig.CodebergToken.Or(data.UnvalidatedConfig.NormalConfig.CodebergToken)
	data.UserInput.Config.NormalConfig.CodebergToken, exit, err = dialog.CodebergToken(existingToken, data.DialogInputs.Next())
	return exit, err
}

func enterGiteaToken(data *SetupData) (exit dialogdomain.Exit, err error) {
	existingToken := data.UserInput.Config.NormalConfig.GiteaToken.Or(data.UnvalidatedConfig.NormalConfig.GiteaToken)
	data.UserInput.Config.NormalConfig.GiteaToken, exit, err = dialog.GiteaToken(existingToken, data.DialogInputs.Next())
	return exit, err
}

func enterGitHubToken(data *SetupData) (exit dialogdomain.Exit, err error) {
	existingToken := data.UserInput.Config.NormalConfig.GitHubToken.Or(data.UnvalidatedConfig.NormalConfig.GitHubToken)
	data.UserInput.Config.NormalConfig.GitHubToken, exit, err = dialog.GitHubToken(existingToken, data.DialogInputs.Next())
	return exit, err
}

func enterGitLabToken(data *SetupData) (exit dialogdomain.Exit, err error) {
	existingToken := data.UserInput.Config.NormalConfig.GitLabToken.Or(data.UnvalidatedConfig.NormalConfig.GitLabToken)
	data.UserInput.Config.NormalConfig.GitLabToken, exit, err = dialog.GitLabToken(existingToken, data.DialogInputs.Next())
	return exit, err
}

func testForgeAuth(data *SetupData, forgeType Option[forgedomain.ForgeType]) (repeat bool, exit dialogdomain.Exit, err error) {
	if _, inTest := os.LookupEnv(subshell.TestToken); inTest {
		return false, false, nil
	}
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              data.Backend,
		BitbucketAppPassword: data.UserInput.Config.NormalConfig.BitbucketAppPassword.Or(data.UnvalidatedConfig.NormalConfig.BitbucketAppPassword),
		BitbucketUsername:    data.UserInput.Config.NormalConfig.BitbucketUsername.Or(data.UnvalidatedConfig.NormalConfig.BitbucketUsername),
		CodebergToken:        data.UserInput.Config.NormalConfig.CodebergToken.Or(data.UnvalidatedConfig.NormalConfig.CodebergToken),
		ForgeType:            forgeType,
		Frontend:             data.Backend,
		GitHubConnectorType:  data.UserInput.Config.NormalConfig.GitHubConnectorType.Or(data.UnvalidatedConfig.NormalConfig.GitHubConnectorType),
		GitHubToken:          data.UserInput.Config.NormalConfig.GitHubToken.Or(data.UnvalidatedConfig.NormalConfig.GitHubToken),
		GitLabConnectorType:  data.UserInput.Config.NormalConfig.GitLabConnectorType.Or(data.UnvalidatedConfig.NormalConfig.GitLabConnectorType),
		GitLabToken:          data.UserInput.Config.NormalConfig.GitLabToken.Or(data.UnvalidatedConfig.NormalConfig.GitLabToken),
		GiteaToken:           data.UserInput.Config.NormalConfig.GiteaToken.Or(data.UnvalidatedConfig.NormalConfig.GiteaToken),
		Log:                  print.Logger{},
		RemoteURL:            data.UserInput.Config.NormalConfig.DevURL().Or(data.UnvalidatedConfig.NormalConfig.DevURL()),
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
		return dialog.CredentialsNoAccess(verifyResult.AuthenticationError, data.DialogInputs.Next())
	}
	if user, hasUser := verifyResult.AuthenticatedUser.Get(); hasUser {
		fmt.Printf(messages.CredentialsForgeUserName, components.FormattedSelection(user, exit))
	}
	if verifyResult.AuthorizationError != nil {
		return dialog.CredentialsNoProposalAccess(verifyResult.AuthorizationError, data.DialogInputs.Next())
	}
	fmt.Println(messages.CredentialsAccess)
	return false, false, nil
}

func enterTokenScope(forgeTypeOpt Option[forgedomain.ForgeType], data *SetupData) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if shouldAskForScope(forgeTypeOpt, data) {
		return tokenScopeDialog(forgeTypeOpt, data)
	}
	return configdomain.ConfigScopeLocal, false, nil
}

func shouldAskForScope(forgeTypeOpt Option[forgedomain.ForgeType], data *SetupData) bool {
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			return existsAndChanged(data.UserInput.Config.NormalConfig.BitbucketUsername, data.UnvalidatedConfig.NormalConfig.BitbucketUsername) &&
				existsAndChanged(data.UserInput.Config.NormalConfig.BitbucketAppPassword, data.UnvalidatedConfig.NormalConfig.BitbucketAppPassword)
		case forgedomain.ForgeTypeCodeberg:
			return existsAndChanged(data.UserInput.Config.NormalConfig.CodebergToken, data.UnvalidatedConfig.NormalConfig.CodebergToken)
		case forgedomain.ForgeTypeGitea:
			return existsAndChanged(data.UserInput.Config.NormalConfig.GiteaToken, data.UnvalidatedConfig.NormalConfig.GiteaToken)
		case forgedomain.ForgeTypeGitHub:
			return existsAndChanged(data.UserInput.Config.NormalConfig.GitHubToken, data.UnvalidatedConfig.NormalConfig.GitHubToken)
		case forgedomain.ForgeTypeGitLab:
			return existsAndChanged(data.UserInput.Config.NormalConfig.GitLabToken, data.UnvalidatedConfig.NormalConfig.GitLabToken)
		}
	}
	return false
}

func tokenScopeDialog(forgeTypeOpt Option[forgedomain.ForgeType], data *SetupData) (configdomain.ConfigScope, dialogdomain.Exit, error) {
	if forgeType, hasForgeType := forgeTypeOpt.Get(); hasForgeType {
		switch forgeType {
		case forgedomain.ForgeTypeBitbucket, forgedomain.ForgeTypeBitbucketDatacenter:
			existingScope := determineScope(data.ConfigSnapshot, configdomain.KeyBitbucketUsername, data.UnvalidatedConfig.NormalConfig.BitbucketUsername)
			return dialog.TokenScope(existingScope, data.DialogInputs.Next())
		case forgedomain.ForgeTypeCodeberg:
			existingScope := determineScope(data.ConfigSnapshot, configdomain.KeyCodebergToken, data.UnvalidatedConfig.NormalConfig.CodebergToken)
			return dialog.TokenScope(existingScope, data.DialogInputs.Next())
		case forgedomain.ForgeTypeGitea:
			existingScope := determineScope(data.ConfigSnapshot, configdomain.KeyGiteaToken, data.UnvalidatedConfig.NormalConfig.GiteaToken)
			return dialog.TokenScope(existingScope, data.DialogInputs.Next())
		case forgedomain.ForgeTypeGitHub:
			existingScope := determineScope(data.ConfigSnapshot, configdomain.KeyGitHubToken, data.UnvalidatedConfig.NormalConfig.GitHubToken)
			return dialog.TokenScope(existingScope, data.DialogInputs.Next())
		case forgedomain.ForgeTypeGitLab:
			existingScope := determineScope(data.ConfigSnapshot, configdomain.KeyGitLabToken, data.UnvalidatedConfig.NormalConfig.GitLabToken)
			return dialog.TokenScope(existingScope, data.DialogInputs.Next())
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
