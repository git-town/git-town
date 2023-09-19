package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/domain"
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
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	if len(args) > 0 {
		newMainBranch := domain.NewLocalBranchName(args[0])
		return setMainBranch(newMainBranch, &repo.Runner)
	}
	printMainBranch(&repo.Runner)
	return nil
}

func printMainBranch(run *git.ProdRunner) {
	cli.Println(cli.StringSetting(run.Config.MainBranch().String()))
}

func setMainBranch(branch domain.LocalBranchName, run *git.ProdRunner) error {
	if !run.Backend.HasLocalBranch(branch) {
		return fmt.Errorf(messages.BranchDoesntExist, branch)
	}
	return run.Config.SetMainBranch(branch)
}
