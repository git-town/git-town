package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/cli/flags"
	"github.com/git-town/git-town/v9/src/cli/print"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const diffParentDesc = "Shows the changes committed to a feature branch"

const diffParentHelp = `
Works on either the current branch or the branch name provided.

Exits with error code 1 if the given branch is a perennial branch or the main branch.`

func diffParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "diff-parent [<branch>]",
		GroupID: "lineage",
		Args:    cobra.MaximumNArgs(1),
		Short:   diffParentDesc,
		Long:    long(diffParentDesc, diffParentHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDiffParent(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeDiffParent(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, exit, err := determineDiffParentConfig(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	err = repo.Runner.Frontend.DiffParent(config.branch, config.parentBranch)
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}

type diffParentConfig struct {
	branch       domain.LocalBranchName
	parentBranch domain.LocalBranchName
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*diffParentConfig, bool, error) {
	lineage := repo.Runner.Config.Lineage()
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, false, err
	}
	branches, _, _, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	branch := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	if branch != branches.Initial {
		if !branches.All.HasLocalBranch(branch) {
			return nil, false, fmt.Errorf(messages.BranchDoesntExist, branch)
		}
	}
	branchTypes := repo.Runner.Config.BranchTypes()
	if !branchTypes.IsFeatureBranch(branch) {
		return nil, false, fmt.Errorf(messages.DiffParentNoFeatureBranch)
	}
	mainBranch := repo.Runner.Config.MainBranch()
	updated, err := validate.KnowsBranchAncestors(branch, validate.KnowsBranchAncestorsArgs{
		DefaultBranch: mainBranch,
		Backend:       &repo.Runner.Backend,
		AllBranches:   branches.All,
		BranchTypes:   branchTypes,
		MainBranch:    mainBranch,
	})
	if err != nil {
		return nil, false, err
	}
	if updated {
		lineage = repo.Runner.Config.Lineage()
	}
	return &diffParentConfig{
		branch:       branch,
		parentBranch: lineage.Parent(branch),
	}, false, nil
}
