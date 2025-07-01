package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state"
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
			return executeRunLog(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRunLog(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := loadRunLogData(repo.RootDir)
	if err != nil {
		return err
	}
	if err = showRunLog(data); err != nil {
		return err
	}
	print.Footer(verbose, *repo.CommandsCounter.Value, []string{})
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

func loadRunLogData(rootDir gitdomain.RepoRootDir) (runLogData, error) {
	filepath, err := state.FilePath(rootDir, state.FileTypeRunlog)
	return runLogData{
		filepath: filepath,
	}, err
}
