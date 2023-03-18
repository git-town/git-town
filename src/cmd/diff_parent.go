package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func diffParentCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:     "diff-parent [<branch>]",
		GroupID: "lineage",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Shows the changes committed to a feature branch",
		Long: `Shows the changes committed to a feature branch

Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiffParent(debug, args)
		},
	}
	debugFlag(&cmd, &debug)
	return &cmd
}

func runDiffParent(debug bool, args []string) error {
	repo, err := Repo(RepoArgs{
		printBranchNames:     false,
		debug:                debug,
		dryRun:               false,
		validateGitversion:   true,
		validateIsRepository: true,
		validateIsConfigured: true,
	})
	if err != nil {
		return err
	}
	config, err := determineDiffParentConfig(args, &repo)
	if err != nil {
		return err
	}
	return repo.DiffParent(config.branch, config.parentBranch)
}

type diffParentConfig struct {
	branch       string
	parentBranch string
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentConfig(args []string, repo *git.PublicRepo) (*diffParentConfig, error) {
	initialBranch, err := repo.CurrentBranch()
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
		hasBranch, err := repo.HasLocalBranch(branch)
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
	err = validate.KnowsBranchAncestry(branch, repo.Config.MainBranch(), repo)
	if err != nil {
		return nil, err
	}
	return &diffParentConfig{
		branch:       branch,
		parentBranch: repo.Config.ParentBranch(branch),
	}, nil
}
