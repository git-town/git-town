package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	runLogDesc = "Displays the repo state before and after previous Git Town commands"
	runLogHelp = `
Git Town records the SHA of all local and remote branches
before and after each command runs
into an immutable, append-only log file called the runlog.

The runlog provides an extra layer of safety,
making it easier to manually roll back changes
if git town undo doesn't fully undo the last command.
`
)

func runLogCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "runlog",
		Args:    cobra.NoArgs,
		GroupID: cmdhelpers.GroupIDErrors,
		Short:   runLogDesc,
		Long:    cmdhelpers.Long(runLogDesc, runLogHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeRunLog(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRunLog(cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    true,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data := loadRunLogData(repo.ConfigDir)
	if err = showRunLog(data); err != nil {
		return err
	}
	print.Footer(repo.UnvalidatedConfig.NormalConfig.Verbose, repo.CommandsCounter.Immutable(), []string{})
	return nil
}

func showRunLog(data runLogData) error {
	fmt.Printf(messages.RunlogDisplaying, data.filepath)
	fmt.Println()
	content, err := os.ReadFile(data.filepath)
	if err != nil {
		return fmt.Errorf(messages.RunLogCannotRead, data.filepath, err)
	}
	fmt.Print(string(content))
	fmt.Printf(messages.RunlogDisplaying, data.filepath)
	return nil
}

type runLogData struct {
	filepath string // filepath of the runlog file
}

func loadRunLogData(configDir configdomain.ConfigDirRepo) runLogData {
	filepath := configDir.RunlogPath()
	return runLogData{filepath: filepath}
}
