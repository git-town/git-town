package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/spf13/cobra"
)

const mainbranchDesc = "Displays or sets your main development branch"

const mainbranchHelp = `
The main branch is the Git branch from which new feature branches are cut.`

func mainbranchConfigCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "main-branch [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: mainbranchDesc,
		Long:  long(mainbranchDesc, mainbranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return configureMainBranch(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func configureMainBranch(args []string, debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		OmitBranchNames: true,
		Debug:           debug,
		DryRun:          false,
	})
	if err != nil {
		return err
	}
	_, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: false,
		ValidateIsConfigured:  false,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setMainBranch(args[0], &run)
	}
	printMainBranch(&run)
	return nil
}

func printMainBranch(run *git.ProdRunner) {
	cli.Println(cli.StringSetting(run.Config.MainBranch()))
}

func setMainBranch(branch string, run *git.ProdRunner) error {
	hasBranch, err := run.Backend.HasLocalBranch(branch)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf(messages.BranchDoesntExist, branch)
	}
	return run.Config.SetMainBranch(branch)
}
