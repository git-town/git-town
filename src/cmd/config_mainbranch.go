package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const mainbranchDesc = "Displays or sets your main development branch"

const mainBranchHelp = `
The main branch is the Git branch from which new feature branches are cut.`

func mainbranchConfigCmd(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "main-branch [<branch>]",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: ensure(repo, isRepository),
		Short:   mainbranchDesc,
		Long:    long(mainbranchDesc, mainBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return configMainBranch(args, repo)
		},
	}
}

func configMainBranch(args []string, repo *git.ProdRepo) error {
	if len(args) > 0 {
		return setMainBranch(args[0], repo)
	}
	printMainBranch(repo)
	return nil
}

func printMainBranch(repo *git.ProdRepo) {
	cli.Println(cli.StringSetting(repo.Config.MainBranch()))
}

func setMainBranch(branch string, repo *git.ProdRepo) error {
	hasBranch, err := repo.Silent.HasLocalBranch(branch)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branch)
	}
	return repo.Config.SetMainBranch(branch)
}
