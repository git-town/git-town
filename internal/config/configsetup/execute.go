package configsetup

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// the config settings to be used if the user accepts all default options
func defaultUserInput(gitAccess gitconfig.Access, gitVersion git.Version) userInput {
	return userInput{
		config:        config.DefaultUnvalidatedConfig(gitAccess, gitVersion),
		configStorage: dialog.ConfigStorageOptionFile,
	}
}

func Execute(data SetupData) error {
	tokenScope, forgeTypeOpt, exit, err := enterData(&data)
	if err != nil || exit {
		return err
	}
	if err = saveAll(data.UserInput, data.UnvalidatedConfig, data.ConfigFile, tokenScope, forgeTypeOpt, data.Git, data.Frontend); err != nil {
		return err
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               data.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   data.ConfigSnapshot,
		Command:               data.CommandName,
		CommandsCounter:       data.CommandsCounter,
		FinalMessages:         data.FinalMessages,
		Git:                   data.Git,
		RootDir:               data.RootDir,
		TouchedBranches:       []gitdomain.BranchName{},
		Verbose:               data.Verbose,
	})
}

type SetupData struct {
	Backend           subshelldomain.RunnerQuerier
	ConfigSnapshot    undoconfig.ConfigSnapshot
	CommandName       string
	CommandsCounter   Mutable[gohacks.Counter]
	ConfigFile        Option[configdomain.PartialConfig]
	DialogInputs      components.TestInputs
	FinalMessages     stringslice.Collector
	Frontend          subshelldomain.Runner
	Git               git.Commands
	LocalBranches     gitdomain.BranchInfos
	Remotes           gitdomain.Remotes
	RootDir           gitdomain.RepoRootDir
	UnvalidatedConfig config.UnvalidatedConfig
	UserInput         userInput
	Verbose           configdomain.Verbose
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
