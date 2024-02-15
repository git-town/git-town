package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/git-town/git-town/v12/src/vm/statefile"
	"github.com/spf13/cobra"
)

const statusDesc = "Displays or resets the current suspended Git Town command"

func statusCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "status",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   statusDesc,
		Long:    cmdhelpers.Long(statusDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeStatus(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	cmd.AddCommand(resetRunstateCommand())
	return &cmd
}

func executeStatus(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, err := loadDisplayStatusConfig(repo.RootDir)
	if err != nil {
		return err
	}
	displayStatus(*config)
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}

type displayStatusConfig struct {
	filepath string             // filepath of the runstate file
	state    *runstate.RunState // content of the runstate file
}

func loadDisplayStatusConfig(rootDir gitdomain.RepoRootDir) (*displayStatusConfig, error) {
	filepath, err := statefile.FilePath(rootDir)
	if err != nil {
		return nil, err
	}
	state, err := statefile.Load(rootDir)
	if err != nil {
		return nil, err
	}
	return &displayStatusConfig{
		filepath: filepath,
		state:    state,
	}, nil
}

func displayStatus(config displayStatusConfig) {
	if config.state == nil {
		fmt.Println(messages.StatusFileNotFound)
		return
	}
	if config.state.IsFinished() {
		displayFinishedStatus(config)
	} else {
		displayUnfinishedStatus(config)
	}
}

func displayUnfinishedStatus(config displayStatusConfig) {
	timeDiff := time.Since(config.state.UnfinishedDetails.EndTime)
	fmt.Printf(messages.PreviousCommandProblem, config.state.Command, timeDiff)
	if config.state.HasAbortProgram() {
		fmt.Println(messages.UndoMessage)
	}
	if config.state.HasRunProgram() {
		fmt.Println(messages.ContinueMessage)
	}
	if config.state.UnfinishedDetails.CanSkip {
		fmt.Println(messages.SkipMessage)
	}
}

func displayFinishedStatus(config displayStatusConfig) {
	fmt.Printf(messages.PreviousCommandFinished, config.state.Command)
	fmt.Println(messages.UndoMessage)
}
