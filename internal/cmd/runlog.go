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
Git Town logs the SHA that all local and remote branches point to
before and after each Git Town command executes.
This is an additional safety net
to allow you to manually undo a Git Town command
in case "git town undo" isn't enough.
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
	err = showRunLog(data)
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
	filepath string // filepath of the runstate file
}

func loadRunLogData(rootDir gitdomain.RepoRootDir) (runLogData, error) {
	filepath, err := state.FilePath(rootDir, state.FileTypeRunlog)
	return runLogData{
		filepath: filepath,
	}, err
}
