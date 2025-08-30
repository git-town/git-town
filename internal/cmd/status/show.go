package status

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				Detached:    None[configdomain.Detached](),
				DryRun:      None[configdomain.DryRun](),
				Stash:       None[configdomain.Stash](),
				Verbose:     verbose,
			})
			return executeStatusShow(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeStatusShow(cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, err := loadDisplayStatusData(repo.RootDir)
	if err != nil {
		return err
	}
	showStatus(data)
	print.Footer(repo.UnvalidatedConfig.NormalConfig.Verbose, *repo.CommandsCounter.Value, []string{})
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
