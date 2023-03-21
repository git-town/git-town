package cmd

import (
	"errors"

	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const setParentDesc = "Prompts to set the parent branch for the current branch"

func setParentCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:     "set-parent",
		GroupID: "lineage",
		Args:    cobra.NoArgs,
		Short:   setParentDesc,
		Long:    long(setParentDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setParent(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func setParent(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
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
	currentBranch, err := repo.CurrentBranch()
	if err != nil {
		return err
	}
	if !repo.Config.IsFeatureBranch(currentBranch) {
		return errors.New("only feature branches can have parent branches")
	}
	existingParent := repo.Config.ParentBranch(currentBranch)
	if existingParent != "" {
		// TODO: delete the old parent only when the user has entered a new parent
		err = repo.Config.RemoveParent(currentBranch)
		if err != nil {
			return err
		}
	} else {
		existingParent = repo.Config.MainBranch()
	}
	return validate.KnowsBranchAncestry(currentBranch, existingParent, &repo)
}
