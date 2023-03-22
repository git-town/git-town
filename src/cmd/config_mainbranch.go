package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const mainbranchDesc = "Displays or sets your main development branch"

const mainbranchHelp = `
The main branch is the Git branch from which new feature branches are cut.`

func mainbranchConfigCmd() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
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
	repo, exit, err := LoadProdRepo(RepoArgs{
		omitBranchNames:       true,
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	if len(args) > 0 {
		return setMainBranch(args[0], &repo)
	}
	printMainBranch(&repo)
	return nil
}

func printMainBranch(repo *git.ProdRepo) {
	cli.Println(cli.StringSetting(repo.Config.MainBranch()))
}

func setMainBranch(branch string, repo *git.ProdRepo) error {
	hasBranch, err := repo.Backend.HasLocalBranch(branch)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branch)
	}
	return repo.Config.SetMainBranch(branch)
}
