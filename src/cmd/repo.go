package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

func repoCommand(repo *git.ProdRepo) *cobra.Command {
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
where HOSTNAME matches what is in your ssh config file.`, config.CodeHostingDriverKey, config.CodeHostingOriginHostnameKey),
		RunE: func(cmd *cobra.Command, args []string) error {
			connector, err := hosting.NewConnector(&repo.Config, &repo.Silent, cli.PrintConnectorAction)
			if err != nil {
				return err
			}
			if connector == nil {
				return hosting.UnsupportedServiceError()
			}
			browser.Open(connector.RepositoryURL(), repo.LoggingShell)
			return nil
		},
		Args:    cobra.NoArgs,
		PreRunE: ensure(repo, hasGitVersion, isRepository, isConfigured, isOnline),
	}
}
