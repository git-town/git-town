package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/validate"
	"github.com/spf13/cobra"
)

const diffParentDesc = "Show the changes committed to a feature branch"

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
		Long:    cmdhelpers.Long(diffParentDesc, diffParentHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDiffParent(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeDiffParent(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, exit, err := determineDiffParentData(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	err = repo.Frontend.DiffParent(config.branch, config.parentBranch)
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.CommandsCounter.Count(), repo.FinalMessages.Result())
	return nil
}

type diffParentData struct {
	branch       gitdomain.LocalBranchName
	parentBranch gitdomain.LocalBranchName
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentData(args []string, repo *execute.OpenRepoResult, verbose bool) (*diffParentData, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                &repo.UnvalidatedConfig.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	branch := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	if branch != branchesSnapshot.Active {
		if !branchesSnapshot.Branches.HasLocalBranch(branch) {
			return nil, false, fmt.Errorf(messages.BranchDoesntExist, branch)
		}
	}
	branchesToDiff := gitdomain.LocalBranchNames{branch}
	localBranches := branchesSnapshot.Branches.LocalBranches()
	validatedConfig, err := validate.Config(repo.UnvalidatedConfig, branchesToDiff, localBranches, &repo.Backend, &dialogTestInputs)
	if err != nil {
		return nil, false, err
	}
	parentBranch, hasParent := validatedConfig.Config.Lineage.Parent(branch).Get()
	if !hasParent {
		return nil, false, errors.New(messages.DiffParentNoFeatureBranch)
	}
	return &diffParentData{
		branch:       branch,
		parentBranch: parentBranch,
	}, false, nil
}
