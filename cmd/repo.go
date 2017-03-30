package cmd

import (
	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/drivers"
	"github.com/Originate/git-town/lib/script"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "Opens the repository homepage",
	Run: func(cmd *cobra.Command, args []string) {
		driver := drivers.GetCodeHostingDriver()
		repository := config.GetUrlRepositoryName(config.GetRemoteOriginUrl())
		script.OpenBrowser(driver.GetRepositoryUrl(repository))
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
