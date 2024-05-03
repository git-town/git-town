package status

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/git-town/git-town/v14/src/vm/statefile"
	"github.com/spf13/cobra"
)

const statusDesc = "Displays or resets the current suspended Git Town command"

func RootCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "status",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   statusDesc,
		Long:    cmdhelpers.Long(statusDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeStatus(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	cmd.AddCommand(resetRunstateCommand())
	return &cmd
}

func executeStatus(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := loadDisplayStatusData(repo.RootDir)
	if err != nil {
		return err
	}
	displayStatus(*data)
	print.Footer(verbose, repo.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}

type displayStatusData struct {
	filepath string                    // filepath of the runstate file
	state    Option[runstate.RunState] // content of the runstate file
}

func loadDisplayStatusData(rootDir gitdomain.RepoRootDir) (*displayStatusData, error) {
	filepath, err := statefile.FilePath(rootDir)
	if err != nil {
		return nil, err
	}
	state, err := statefile.Load(rootDir)
	if err != nil {
		return nil, err
	}
	return &displayStatusData{
		filepath: filepath,
		state:    state,
	}, nil
}

func displayStatus(data displayStatusData) {
	state, hasState := data.state.Get()
	if !hasState {
		fmt.Println(messages.StatusFileNotFound)
		return
	}
	if state.IsFinished() {
		displayFinishedStatus(state)
	} else {
		displayUnfinishedStatus(state)
	}
}

func displayUnfinishedStatus(state runstate.RunState) {
	timeDiff := time.Since(state.UnfinishedDetails.EndTime)
	fmt.Printf(messages.PreviousCommandProblem, state.Command, timeDiff)
	if state.HasAbortProgram() {
		fmt.Println(messages.UndoMessage)
	}
	if state.HasRunProgram() {
		fmt.Println(messages.ContinueMessage)
	}
	if state.UnfinishedDetails.CanSkip {
		fmt.Println(messages.SkipMessage)
	}
}

func displayFinishedStatus(state runstate.RunState) {
	fmt.Printf(messages.PreviousCommandFinished, state.Command)
	fmt.Println(messages.UndoMessage)
}
