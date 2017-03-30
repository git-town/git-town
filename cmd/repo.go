package cmd

import (
	"errors"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/drivers"
	"github.com/Originate/git-town/lib/script"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "View the repository homepage",
	Long:  `View the repository homepage`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := drivers.GetCodeHostingDriver()
		repository := config.GetUrlRepositoryName(config.GetRemoteOriginUrl())
		script.OpenBrowser(driver.GetRepositoryUrl(repository))
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("Too many arguments")
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
