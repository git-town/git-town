package status

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/flags"
	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/execute"
	"github.com/spf13/cobra"
)

const statusShowDesc = "Displays the detailed information from the persisted runstate"

func showRunstateCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "show",
		Args:  cobra.NoArgs,
		Short: statusShowDesc,
		Long:  cmdhelpers.Long(statusShowDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeStatusShow(verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeStatusShow(verbose configdomain.Verbose) error {
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
	data, err := loadDisplayStatusData(repo.RootDir)
	if err != nil {
		return err
	}
	showStatus(data)
	print.Footer(verbose, *repo.CommandsCounter.Value, []string{})
	return nil
}

func showStatus(data displayStatusData) {
	state, hasState := data.state.Get()
	if !hasState {
		return
	}
	fmt.Println("Displaying runstate at", data.filepath)
	fmt.Println()
	fmt.Println(state.String())
}
