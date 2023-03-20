package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func mainbranchConfigCmd() *cobra.Command {
	debug := false
	cmd := &cobra.Command{
		Use:   "main-branch [<branch>]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Displays or sets your main development branch",
		Long: `Displays or sets your main development branch

The main branch is the Git branch from which new feature branches are cut.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfigureMainBranch(debug, args)
		},
	}
	debugFlag(cmd, &debug)
	return cmd
}

func runConfigureMainBranch(debug bool, args []string) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
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

func printMainBranch(repo *git.PublicRepo) {
	cli.Println(cli.StringSetting(repo.Config.MainBranch()))
}

func setMainBranch(branch string, repo *git.PublicRepo) error {
	hasBranch, err := repo.HasLocalBranch(branch)
	if err != nil {
		return err
	}
	if !hasBranch {
		return fmt.Errorf("there is no branch named %q", branch)
	}
	return repo.Config.SetMainBranch(branch)
}
