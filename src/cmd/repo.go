package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/browser"
	"github.com/git-town/git-town/v9/src/cli"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/hosting"
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
	repo, exit, err := execute.OpenRepo(execute.OpenShellArgs{
		Debug:                 debug,
		DryRun:                false,
		Fetch:                 false,
		HandleUnfinishedState: false,
		OmitBranchNames:       false,
		ValidateIsOnline:      true,
		ValidateGitRepo:       true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return err
	}
	config, err := determineRepoConfig(&repo.Runner)
	if err != nil {
		return err
	}
	browser.Open(config.connector.RepositoryURL(), repo.Runner.Frontend.FrontendRunner, repo.Runner.Backend.BackendRunner)
	repo.Runner.Stats.PrintAnalysis()
	return nil
}

func determineRepoConfig(run *git.ProdRunner) (*repoConfig, error) {
	_, err := execute.LoadBranches(run, execute.LoadBranchesArgs{
		ValidateIsConfigured: true,
	})
	if err != nil {
		return nil, err
	}
	originURL := run.Config.OriginURL()
	hostingService, err := run.Config.HostingService()
	if err != nil {
		return nil, err
	}
	mainBranch := run.Config.MainBranch()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetShaForBranch: run.Backend.ShaForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   run.Config.GiteaToken(),
		GithubAPIToken:  run.Config.GitHubToken(),
		GitlabAPIToken:  run.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             cli.PrintingLog{},
	})
	if err != nil {
		return nil, err
	}
	if connector == nil {
		return nil, hosting.UnsupportedServiceError()
	}
	return &repoConfig{
		connector: connector,
	}, err
}

type repoConfig struct {
	connector hosting.Connector
}
