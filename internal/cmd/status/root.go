package status

import (
	"cmp"
	"fmt"
	"time"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const statusDesc = "Displays or resets the current suspended Git Town command"

func RootCommand() *cobra.Command {
	addPendingFlag, readPendingFlag := flags.Pending()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "status",
		GroupID: cmdhelpers.GroupIDErrors,
		Args:    cobra.NoArgs,
		Short:   statusDesc,
		Long:    cmdhelpers.Long(statusDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			pending, errPending := readPendingFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errPending, errVerbose); err != nil {
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
			return executeStatus(cliConfig, pending)
		},
	}
	addPendingFlag(&cmd)
	addVerboseFlag(&cmd)
	cmd.AddCommand(resetRunstateCommand())
	cmd.AddCommand(showRunstateCommand())
	return &cmd
}

func executeStatus(cliConfig configdomain.PartialConfig, pending configdomain.Pending) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    true,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
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
		print.Footer(repo.UnvalidatedConfig.NormalConfig.Verbose, repo.CommandsCounter.Immutable(), []string{})
	}
	return nil
}

type displayStatusData struct {
	filepath string                    // filepath of the runstate file
	state    Option[runstate.RunState] // content of the runstate file
}

func loadDisplayStatusData(rootDir gitdomain.RepoRootDir) (displayStatusData, error) {
	filepath, err := state.FilePath(rootDir, state.FileTypeRunstate)
	if err != nil {
		return displayStatusData{}, err
	}
	state, err := runstate.Load(rootDir)
	if err != nil {
		return displayStatusData{}, err
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
