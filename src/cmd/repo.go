package cmd

import (
	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

var repoCommand = &cobra.Command{
	Use:   "repo",
	Short: "Opens the repository homepage",
	Long: `Opens the repository homepage

Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket.
Derives the Git provider from the "origin" remote.
You can override this detection with
"git config git-town.code-hosting-driver <DRIVER>"
where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run
"git config git-town.code-hosting-origin-hostname <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		driver := hosting.NewDriver(prodRepo.Config.Hosting, &prodRepo.Silent, cli.PrintDriverAction)
		if driver == nil {
			cli.Exit(hosting.UnsupportedServiceError())
		}
		browser.Open(driver.RepositoryURL(), prodRepo.LoggingShell)
	},
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := ValidateIsRepository(prodRepo); err != nil {
			return err
		}
		if err := validateIsConfigured(prodRepo); err != nil {
			return err
		}
		if err := prodRepo.Config.Offline.Validate(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(repoCommand)
}
