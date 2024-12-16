package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cli/print"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/execute"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/slice"
	"github.com/git-town/git-town/v17/internal/hosting"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/validate"
	. "github.com/git-town/git-town/v17/pkg/prelude"
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
		GroupID: "stack",
		Args:    cobra.MaximumNArgs(1),
		Short:   diffParentDesc,
		Long:    cmdhelpers.Long(diffParentDesc, diffParentHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeDiffParent(args, verbose)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeDiffParent(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineDiffParentData(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	err = repo.Git.DiffParent(repo.Frontend, data.branch, data.parentBranch)
	if err != nil {
		return err
	}
	print.Footer(verbose, repo.CommandsCounter.Immutable(), repo.FinalMessages.Result())
	return nil
}

type diffParentData struct {
	branch       gitdomain.LocalBranchName
	parentBranch gitdomain.LocalBranchName
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentData(args []string, repo execute.OpenRepoResult, verbose configdomain.Verbose) (data diffParentData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branch := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, currentBranch.String()))
	if branch != currentBranch {
		if !branchesSnapshot.Branches.HasLocalBranch(branch) {
			return data, false, fmt.Errorf(messages.BranchDoesntExist, branch)
		}
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{branch},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	parentBranch, hasParent := validatedConfig.NormalConfig.Lineage.Parent(branch).Get()
	if !hasParent {
		return data, false, errors.New(messages.DiffParentNoFeatureBranch)
	}
	return diffParentData{
		branch:       branch,
		parentBranch: parentBranch,
	}, false, nil
}
