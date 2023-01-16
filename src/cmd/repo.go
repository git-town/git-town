package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

func repoCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "repo",
		Short: "Opens the repository homepage",
		Long: fmt.Sprintf(`Opens the repository homepage

Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket.
Derives the Git provider from the "origin" remote.
You can override this detection with
"git config %s <DRIVER>"
where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run
"git config %s <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`, config.CodeHostingDriver, config.CodeHostingOriginHostname),
		Run: func(cmd *cobra.Command, args []string) {
			driver := hosting.NewDriver(&prodRepo.Config, &prodRepo.Silent, cli.PrintDriverAction)
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
			if err := prodRepo.Config.ValidateIsOnline(); err != nil {
				return err
			}
			return nil
		},
	}
}
