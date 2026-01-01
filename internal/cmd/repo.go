package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	repoDesc = "Open the repository homepage in the browser"
	repoHelp = `
Supported for repositories hosted on
GitHub, GitLab, Gitea, Bitbucket, and Forgejo.
Derives the Git provider from the "origin" remote.
You can override this detection with
"git config %s <DRIVER>"
where DRIVER is "github", "gitlab", "gitea", or "bitbucket".

When using SSH identities, run
"git config %s <HOSTNAME>"
where HOSTNAME matches what is in your ssh config file.`
)

func repoCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "repo [remote]",
		Args:  cobra.MaximumNArgs(1),
		Short: repoDesc,
		Long:  cmdhelpers.Long(repoDesc, fmt.Sprintf(repoHelp, configdomain.KeyForgeType, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeRepo(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRepo(args []string, cliConfig configdomain.PartialConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
	})
	if err != nil {
		return err
	}
	data, err := determineRepoData(args, repo)
	if err != nil {
		return err
	}
	err = data.connector.BrowseRepository(repo.Frontend)
	print.Footer(repo.UnvalidatedConfig.NormalConfig.Verbose, repo.CommandsCounter.Immutable(), repo.FinalMessages.Result())
	return err
}

func determineRepoData(args []string, repo execute.OpenRepoResult) (data repoData, err error) {
	var remoteOpt Option[gitdomain.Remote]
	if len(args) > 0 {
		remoteOpt = gitdomain.NewRemote(args[0])
	} else {
		remoteOpt = Some(repo.UnvalidatedConfig.NormalConfig.DevRemote)
	}
	remote, hasRemote := remoteOpt.Get()
	if !hasRemote {
		return repoData{connector: nil}, nil
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
		Frontend:             repo.Frontend,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.RemoteURL(repo.Backend, remote),
	})
	if err != nil {
		return data, err
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return data, forgedomain.UnsupportedServiceError()
	}
	return repoData{
		connector: connector,
	}, nil
}

type repoData struct {
	connector forgedomain.Connector
}
