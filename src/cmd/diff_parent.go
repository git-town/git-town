package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/flags"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const diffParentDesc = "Shows the changes committed to a feature branch"

const diffParentHelp = `
Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`

func diffParentCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "diff-parent [<branch>]",
		GroupID: "lineage",
		Args:    cobra.MaximumNArgs(1),
		Short:   diffParentDesc,
		Long:    long(diffParentDesc, diffParentHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return diffParent(args, readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func diffParent(args []string, debug bool) error {
	repo, exit, err := LoadProdRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: true,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineDiffParentConfig(args, &repo)
	if err != nil {
		return err
	}
	return repo.Frontend.DiffParent(config.branch, config.parentBranch)
}

type diffParentConfig struct {
	branch       string
	parentBranch string
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentConfig(args []string, repo *git.ProdRepo) (*diffParentConfig, error) {
	initialBranch, err := repo.Backend.CurrentBranch()
	if err != nil {
		return nil, err
	}
	var branch string
	if len(args) > 0 {
		branch = args[0]
	} else {
		branch = initialBranch
	}
	if initialBranch != branch {
		hasBranch, err := repo.Backend.HasLocalBranch(branch)
		if err != nil {
			return nil, err
		}
		if !hasBranch {
			return nil, fmt.Errorf("there is no local branch named %q", branch)
		}
	}
	if !repo.Config.IsFeatureBranch(branch) {
		return nil, fmt.Errorf("you can only diff-parent feature branches")
	}
	err = validate.KnowsBranchAncestry(branch, repo.Config.MainBranch(), &repo.Backend)
	if err != nil {
		return nil, err
	}
	return &diffParentConfig{
		branch:       branch,
		parentBranch: repo.Config.ParentBranch(branch),
	}, nil
}
