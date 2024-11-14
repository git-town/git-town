package status

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	"github.com/git-town/git-town/v16/internal/vm/statefile"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const statusDesc = "Displays or resets the current suspended Git Town command"

func RootCommand() *cobra.Command {
	addPendingFlag, readPendingFlag := flags.Pending()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "status",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   statusDesc,
		Long:    cmdhelpers.Long(statusDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			pending, err := readPendingFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeStatus(pending, verbose)
		},
	}
	addPendingFlag(&cmd)
	addVerboseFlag(&cmd)
	cmd.AddCommand(resetRunstateCommand())
	return &cmd
}

func executeStatus(pending configdomain.Pending, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		if pending {
			return nil
		}
		return err
	}
	data, err := loadDisplayStatusData(repo.RootDir)
	if err != nil {
		return err
	}
	displayStatus(data, pending)
	if !pending {
		print.Footer(verbose, *repo.CommandsCounter.Value, print.NoFinalMessages)
	}
	return nil
}

type displayStatusData struct {
	filepath string                    // filepath of the runstate file
	state    Option[runstate.RunState] // content of the runstate file
}

func loadDisplayStatusData(rootDir gitdomain.RepoRootDir) (result displayStatusData, err error) {
	filepath, err := statefile.FilePath(rootDir)
	if err != nil {
		return result, err
	}
	state, err := statefile.Load(rootDir)
	if err != nil {
		return result, err
	}
	return displayStatusData{
		filepath: filepath,
		state:    state,
	}, nil
}

func displayStatus(data displayStatusData, pending configdomain.Pending) {
	state, hasState := data.state.Get()
	if !hasState {
		if !pending {
			fmt.Println(messages.StatusFileNotFound)
		}
		return
	}
	if state.IsFinished() {
		displayFinishedStatus(state, pending)
	} else {
		displayUnfinishedStatus(state, pending)
	}
}

func displayUnfinishedStatus(state runstate.RunState, pending configdomain.Pending) {
	unfinishedDetails, hasUnfinishedDetails := state.UnfinishedDetails.Get()
	if pending {
		if hasUnfinishedDetails {
			fmt.Print(state.Command)
		}
		return
	}
	timeDiff := time.Since(unfinishedDetails.EndTime)
	fmt.Printf(messages.PreviousCommandProblem, state.Command, timeDiff)
	if state.HasAbortProgram() {
		fmt.Println(messages.UndoMessage)
	}
	if state.HasRunProgram() {
		fmt.Println(messages.ContinueMessage)
	}
	if hasUnfinishedDetails {
		if unfinishedDetails.CanSkip {
			fmt.Println(messages.SkipMessage)
		}
	}
}

func displayFinishedStatus(state runstate.RunState, pending configdomain.Pending) {
	if !pending {
		fmt.Printf(messages.PreviousCommandFinished, state.Command)
		fmt.Println(messages.UndoMessage)
	}
}
