package cmd

import (
	"github.com/Originate/git-town/lib/drivers"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/prompt"
	"github.com/Originate/git-town/lib/script"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "Opens the repository homepage",
	Run: func(cmd *cobra.Command, args []string) {
		git.EnsureIsRepository()
		prompt.EnsureIsConfigured()
		driver := drivers.GetCodeHostingDriver()
		repository := git.GetURLRepositoryName(git.GetRemoteOriginURL())
		script.OpenBrowser(driver.GetRepositoryURL(repository))
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
