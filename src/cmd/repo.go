package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v13/src/browser"
	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/cli/flags"
	"github.com/git-town/git-town/v13/src/cli/print"
	"github.com/git-town/git-town/v13/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v13/src/config/gitconfig"
	"github.com/git-town/git-town/v13/src/execute"
	"github.com/git-town/git-town/v13/src/hosting"
	"github.com/git-town/git-town/v13/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v13/src/validate"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeRepo(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRepo(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, err := determineRepoConfig(repo)
	if err != nil {
		return err
	}
	browser.Open(config.connector.RepositoryURL(), repo.Runner.Frontend.Runner, repo.Runner.Backend.Runner)
	print.Footer(verbose, repo.Runner.CommandsCounter.Count(), repo.Runner.FinalMessages.Result())
	return nil
}

func determineRepoConfig(repo *execute.OpenRepoResult) (*repoConfig, error) {
	branchesSnapshot, err := repo.Runner.Backend.BranchesSnapshot()
	if err != nil {
		return nil, err
	}
	dialogInputs := components.LoadTestInputs(os.Environ())
	err = validate.IsConfigured(&repo.Runner.Backend, &repo.Runner.Config.FullConfig, branchesSnapshot.Branches.LocalBranches().Names(), &dialogInputs)
	if err != nil {
		return nil, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:      &repo.Runner.Config.FullConfig,
		HostingPlatform: repo.Runner.Config.FullConfig.HostingPlatform,
		Log:             print.Logger{},
		OriginURL:       repo.Runner.Config.OriginURL(),
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
