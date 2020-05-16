package cmd

import (
	"github.com/git-town/git-town/src/drivers"
	"github.com/git-town/git-town/src/git"
	"github.com/git-town/git-town/src/script"
	"github.com/git-town/git-town/src/util"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "Opens the repository homepage",
	Long: `Opens the repository homepage

Supported only for repositories hosted on GitHub, GitLab, and Bitbucket.
When using self-hosted versions this command needs to be configured with
"git config git-town.code-hosting-driver <DRIVER>"
where DRIVER is "github", "gitlab", or "bitbucket".
When using SSH identities, this command needs to be configured with
"git config git-town.code-hosting-origin-hostname <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := drivers.GetActiveDriver()
		script.OpenBrowser(driver.GetRepositoryURL())
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return util.FirstError(
			git.ValidateIsRepository,
			validateIsConfigured,
			git.Config().ValidateIsOnline,
			drivers.ValidateHasDriver,
		)
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
