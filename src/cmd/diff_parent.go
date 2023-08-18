package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/validate"
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
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		HandleUnfinishedState: true,
		OmitBranchNames:       false,
		ValidateIsOnline:      false,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineDiffParentConfig(args, &repo.Runner)
	if err != nil {
		return err
	}
	err = repo.Runner.Frontend.DiffParent(config.branch, config.parentBranch)
	if err != nil {
		return err
	}
	repo.Runner.Stats.PrintAnalysis()
	return nil
}

type diffParentConfig struct {
	branch       domain.LocalBranchName
	parentBranch domain.LocalBranchName
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentConfig(args []string, run *git.ProdRunner) (*diffParentConfig, error) {
	branches, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
	branch := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	if branch != branches.Initial {
		hasBranch, err := run.Backend.HasLocalBranch(branch)
		if err != nil {
			return nil, err
		}
		if !hasBranch {
			return nil, fmt.Errorf(messages.BranchDoesntExist, branch)
		}
	}
	branchDurations := run.Config.BranchDurations()
	if !branchDurations.IsFeatureBranch(branch) {
		return nil, fmt.Errorf(messages.DiffParentNoFeatureBranch)
	}
	mainBranch := run.Config.MainBranch()
	lineage := run.Config.Lineage()
	updated, err := validate.KnowsBranchAncestors(branch, validate.KnowsBranchAncestorsArgs{
		DefaultBranch:   mainBranch,
		Backend:         &run.Backend,
		AllBranches:     branches.All,
		Lineage:         lineage,
		BranchDurations: branchDurations,
		MainBranch:      mainBranch,
	})
	if err != nil {
		return nil, err
	}
	if updated {
		lineage = run.Config.Lineage()
	}
	return &diffParentConfig{
		branch:       branch,
		parentBranch: lineage.Parent(branch),
	}, nil
}
