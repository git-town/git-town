package cmd

import (
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/validate"
	"github.com/spf13/cobra"
)

const setupConfigDesc = "Prompts to setup your Git Town configuration"

func setupConfigCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "setup",
		Args:  cobra.NoArgs,
		Short: setupConfigDesc,
		Long:  long(setupConfigDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setup(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func setup(debug bool) error {
	run, err := execute.LoadProdRunner(execute.LoadArgs{
		OmitBranchNames: true,
		Debug:           debug,
		DryRun:          false,
	})
	if err != nil {
		return err
	}
	_, _, exit, err := execute.LoadGitRepo(&run, execute.LoadGitArgs{
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsConfigured:  false,
		ValidateIsOnline:      false,
		ValidateIsRepository:  true,
	})
	if err != nil || exit {
		return err
	}
	mainBranch, err := validate.EnterMainBranch(&run.Backend)
	if err != nil {
		return err
	}
	return validate.EnterPerennialBranches(&run.Backend, mainBranch)
}
