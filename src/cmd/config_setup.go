package cmd

import (
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

const configSetupDesc = "Prompts to setup your Git Town configuration"

func setupConfigCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "setup",
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, isRepository),
		Short:   configSetupDesc,
		Long:    long(configSetupDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return configSetup(repo)
		},
	}
}

func configSetup(repo *git.ProdRepo) error {
	err := validate.EnterMainBranch(repo)
	if err != nil {
		return err
	}
	return validate.EnterPerennialBranches(repo)
}
