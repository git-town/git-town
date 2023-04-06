package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/browser"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/execute"
	"github.com/git-town/git-town/v7/src/flags"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/spf13/cobra"
)

const repoDesc = "Opens the repository homepage"

const repoHelp = `
Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket.
Derives the Git provider from the "origin" remote.
You can override this detection with
"git config %s <DRIVER>"
where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run
"git config %s <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`

func repoCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:   "repo",
		Args:  cobra.NoArgs,
		Short: repoDesc,
		Long:  long(repoDesc, fmt.Sprintf(repoHelp, config.CodeHostingDriverKey, config.CodeHostingOriginHostnameKey)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return repo(readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func repo(debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
		ValidateIsRepository:  true,
		ValidateIsConfigured:  true,
		ValidateIsOnline:      true,
	})
	if err != nil || exit {
		return err
	}
	connector, err := hosting.NewConnector(run.Config.GitTown, &run.Backend, cli.PrintConnectorAction)
	if err != nil {
		return err
	}
	if connector == nil {
		return hosting.UnsupportedServiceError()
	}
	browser.Open(connector.RepositoryURL(), run.Frontend.FrontendRunner, run.Backend.BackendRunner)
	run.Stats.PrintAnalysis()
	return nil
}
