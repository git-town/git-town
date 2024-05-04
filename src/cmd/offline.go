package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/format"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
	"github.com/spf13/cobra"
)

const offlineDesc = "Display or set offline mode"

const offlineHelp = `
Git Town avoids network operations in offline mode.`

func offlineCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "offline [(yes | no)]",
		Args:    cobra.MaximumNArgs(1),
		GroupID: "setup",
		Short:   offlineDesc,
		Long:    cmdhelpers.Long(offlineDesc, offlineHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeOffline(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeOffline(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateGitRepo:  false,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	switch len(args) {
	case 0:
		displayOfflineStatus(repo.Config.Config)
	case 1:
		err = setOfflineStatus(args[0], *repo.Config)
		if err != nil {
			return err
		}
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             repo.Backend,
		BeginConfigSnapshot: repo.ConfigSnapshot,
		Command:             "offline",
		CommandsCounter:     repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       repo.FinalMessages,
		RootDir:             repo.RootDir,
		Verbose:             verbose,
	})
}

func displayOfflineStatus(config configdomain.FullConfig) {
	fmt.Println(format.Bool(config.Offline.Bool()))
}

func setOfflineStatus(text string, config config.Config) error {
	value, err := gohacks.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.ValueInvalid, gitconfig.KeyOffline, text)
	}
	return config.SetOffline(configdomain.Offline(value))
}
