package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/browser"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/spf13/cobra"
)

const repoDesc = "Open the repository homepage in the browser"

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
	data, err := determineRepoData(repo)
	if err != nil {
		return err
	}
	browser.Open(data.connector.RepositoryURL(), repo.Frontend, repo.Backend)
	print.Footer(verbose, repo.CommandsCounter.Count(), repo.FinalMessages.Result())
	return nil
}

func determineRepoData(repo execute.OpenRepoResult) (repoData, error) {
	var err error
	var connectorOpt Option[hostingdomain.Connector]
	if originURL, hasOriginURL := repo.UnvalidatedConfig.OriginURL().Get(); hasOriginURL {
		connectorOpt, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *repo.UnvalidatedConfig.Config,
			HostingPlatform: repo.UnvalidatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return emptyRepoData(), err
		}
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return emptyRepoData(), hostingdomain.UnsupportedServiceError()
	}
	return repoData{
		connector: connector,
	}, err
}

type repoData struct {
	connector hostingdomain.Connector
}

func emptyRepoData() repoData {
	return repoData{} //exhaustruct:ignore
}
