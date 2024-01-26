package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v11/src/browser"
	"github.com/git-town/git-town/v11/src/cli/dialog"
	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/hosting"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v11/src/validate"
	"github.com/spf13/cobra"
)

const repoDesc = "Opens the repository homepage"

const repoHelp = `
Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket. Derives the Git provider from the "origin" remote. You can override this detection with "git config %s <DRIVER>" where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run "git config %s <HOSTNAME>" where HOSTNAME matches what is in your ssh config file.`

func repoCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "repo",
		Args:  cobra.NoArgs,
		Short: repoDesc,
		Long:  cmdhelpers.Long(repoDesc, fmt.Sprintf(repoHelp, gitconfig.KeyHostingPlatform, gitconfig.KeyHostingOriginHostname)),
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
	dialogInputs := dialog.LoadTestInputs(os.Environ())
	err = validate.IsConfigured(&repo.Runner.Backend, &repo.Runner.FullConfig, branchesSnapshot.Branches.LocalBranches().Names(), &dialogInputs)
	if err != nil {
		return nil, err
	}
	hostingService, err := repo.Runner.HostingService()
	if err != nil {
		return nil, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:     &repo.Runner.FullConfig,
		HostingService: hostingService,
		OriginURL:      repo.Runner.OriginURL(),
		Log:            print.Logger{},
	})
	if err != nil {
		return nil, err
	}
	if connector == nil {
		return nil, hostingdomain.UnsupportedServiceError()
	}
	return &repoConfig{
		connector: connector,
	}, err
}

type repoConfig struct {
	connector hostingdomain.Connector
}
