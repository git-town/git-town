package cmd

import (
	"github.com/Originate/git-town/lib/drivers"
	"github.com/Originate/git-town/lib/gitconfig"
	"github.com/Originate/git-town/lib/script"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "View the repository homepage",
	Long:  `View the repository homepage`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := drivers.GetCodeHostingDriver()
		repository := gitconfig.GetUrlRepositoryName(gitconfig.GetRemoteOriginUrl())
		script.OpenBrowser(driver.GetRepositoryUrl(repository))
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
