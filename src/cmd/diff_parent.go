package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
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
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:           debug,
		DryRun:          false,
		OmitBranchNames: false,
	})
	if err != nil {
		return err
	}
	_, currentBranch, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		Fetch:                 false,
		HandleUnfinishedState: true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      false,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineDiffParentConfig(args, &run, currentBranch)
	if err != nil {
		return err
	}
	err = run.Frontend.DiffParent(config.branch, config.parentBranch)
	if err != nil {
		return err
	}
	run.Stats.PrintAnalysis()
	return nil
}

type diffParentConfig struct {
	branch       string
	parentBranch string
}

// Does not return error because "Ensure" functions will call exit directly.
func determineDiffParentConfig(args []string, run *git.ProdRunner, initialBranch string) (*diffParentConfig, error) {
	var branch string
	if len(args) > 0 {
		branch = args[0]
	} else {
		branch = initialBranch
	}
	if initialBranch != branch {
		hasBranch, err := run.Backend.HasLocalBranch(branch)
		if err != nil {
			return nil, err
		}
		if !hasBranch {
			return nil, fmt.Errorf(messages.BranchDoesntExist, branch)
		}
	}
	if !run.Config.IsFeatureBranch(branch) {
		return nil, fmt.Errorf(messages.DiffParentNoFeatureBranch)
	}
	err := validate.KnowsBranchAncestors(branch, run.Config.MainBranch(), &run.Backend)
	if err != nil {
		return nil, err
	}
	lineage := run.Config.Lineage()
	return &diffParentConfig{
		branch:       branch,
		parentBranch: lineage.Parent(branch),
	}, nil
}
