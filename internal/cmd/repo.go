package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/browser"
	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/spf13/cobra"
)

const repoDesc = "Open the repository homepage in the browser"

const repoHelp = `
Supported for repositories hosted on GitHub, GitLab, Gitea, and Bitbucket. Derives the Git provider from the "origin" remote. You can override this detection with "git config %s <DRIVER>" where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run "git config %s <HOSTNAME>" where HOSTNAME matches what is in your ssh config file.`

func repoCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "repo [remote]",
		Args:  cobra.MaximumNArgs(1),
		Short: repoDesc,
		Long:  cmdhelpers.Long(repoDesc, fmt.Sprintf(repoHelp, configdomain.KeyHostingPlatform, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRepo(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRepo(args []string, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, err := determineRepoData(args, repo)
	if err != nil {
		return err
	}
	browser.Open(data.connector.RepositoryURL(), repo.Frontend, repo.Backend)
	print.Footer(verbose, repo.CommandsCounter.Get(), repo.FinalMessages.Result())
	return nil
}

func determineRepoData(args []string, repo execute.OpenRepoResult) (data repoData, err error) {
	var remoteOpt Option[gitdomain.Remote]
	if len(args) > 0 {
		remoteOpt = Some(gitdomain.NewRemote(args[0]))
	} else {
		remoteOpt = Some(gitdomain.RemoteOrigin)
	}
	remote, hasRemote := remoteOpt.Get()
	if !hasRemote {
		return repoData{connector: nil}, nil
	}
	var connectorOpt Option[hostingdomain.Connector]
	if remoteURL, hasRemoteURL := repo.UnvalidatedConfig.RemoteURL(remote).Get(); hasRemoteURL {
		connectorOpt, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          repo.UnvalidatedConfig.Config.Get(),
			HostingPlatform: repo.UnvalidatedConfig.Config.Value.HostingPlatform,
			Log:             print.Logger{},
			RemoteURL:       remoteURL,
		})
		if err != nil {
			return data, err
		}
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return data, hostingdomain.UnsupportedServiceError()
	}
	return repoData{
		connector: connector,
	}, err
}

type repoData struct {
	connector hostingdomain.Connector
}
