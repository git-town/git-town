package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/browser"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/log"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/hosting"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/validate"
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
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "repo",
		Args:  cobra.NoArgs,
		Short: repoDesc,
		Long:  long(repoDesc, fmt.Sprintf(repoHelp, configdomain.KeyCodeHostingPlatform, configdomain.KeyCodeHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRepo(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRepo(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: true,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, err := determineRepoConfig(repo)
	if err != nil {
		return err
	}
	browser.Open(config.connector.RepositoryURL(), repo.Runner.Frontend.FrontendRunner, repo.Runner.Backend.BackendRunner)
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), print.NoFinalMessages)
	return nil
}

func determineRepoConfig(repo *execute.OpenRepoResult) (*repoConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return nil, err
	}
	branchTypes := repo.Runner.GitTown.BranchTypes()
	branches := domain.Branches{
		All:     branchesSnapshot.Branches,
		Types:   branchTypes,
		Initial: branchesSnapshot.Active,
	}
	_, err = validate.IsConfigured(&repo.Runner.Backend, branches)
	if err != nil {
		return nil, err
	}
	originURL := repo.Runner.GitTown.OriginURL()
	hostingService, err := repo.Runner.GitTown.HostingService()
	if err != nil {
		return nil, err
	}
	mainBranch := repo.Runner.GitTown.MainBranch()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.GitTown.GiteaToken(),
		GithubAPIToken:  github.GetAPIToken(repo.Runner.GitTown.GitHubToken()),
		GitlabAPIToken:  repo.Runner.GitTown.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             log.Printing{},
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
