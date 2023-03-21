package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

const repoSummary = "Opens the repository homepage"

const repoDesc = `Opens the repository homepage

Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket.
Derives the Git provider from the "origin" remote.
You can override this detection with
"git config %s <DRIVER>"
where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run
"git config %s <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`

func repoCommand() *cobra.Command {
	debug := false
	cmd := cobra.Command{
		Use:   "repo",
		Args:  cobra.NoArgs,
		Short: repoSummary,
		Long:  long(repoSummary, fmt.Sprintf(repoDesc, config.CodeHostingDriverKey, config.CodeHostingOriginHostnameKey)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRepo(debug)
		},
	}
	debugFlagOld(&cmd, &debug)
	return &cmd
}

func runRepo(debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		handleUnfinishedState: false,
		validateGitversion:    true,
		validateIsRepository:  true,
		validateIsConfigured:  true,
		validateIsOnline:      true,
	})
	if err != nil || exit {
		return err
	}
	connector, err := hosting.NewConnector(repo.Config, &repo.InternalRepo, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	if connector == nil {
		return hosting.UnsupportedServiceError()
	}
	browser.Open(connector.RepositoryURL(), repo.Public, repo.InternalRepo.InternalRunner)
	return nil
}
