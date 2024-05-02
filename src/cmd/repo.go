package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/browser"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/validate"
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
	data, abort, err := determineRepoData(repo)
	if err != nil || abort {
		return err
	}
	browser.Open(data.connector.RepositoryURL(), repo.Frontend.Runner, repo.Backend.Runner)
	print.Footer(verbose, repo.CommandsCounter.Count(), repo.FinalMessages.Result())
	return nil
}

func determineRepoData(repo *execute.OpenRepoResult) (*repoData, bool, error) {
	branchesSnapshot, err := repo.Backend.BranchesSnapshot()
	if err != nil {
		return nil, false, err
	}
	dialogInputs := components.LoadTestInputs(os.Environ())
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, abort, err := validate.Config(validate.ConfigArgs{
		Unvalidated:        repo.UnvalidatedConfig,
		BranchesToValidate: localBranches,
		LocalBranches:      localBranches,
		Backend:            &repo.Backend,
		TestInputs:         &dialogInputs,
	})
	if err != nil || abort {
		return nil, abort, err
	}
	var connector hostingdomain.Connector
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          &validatedConfig.Config,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, false, err
		}
	}
	if connector == nil {
		return nil, false, hostingdomain.UnsupportedServiceError()
	}
	return &repoData{
		connector: connector,
	}, false, err
}

type repoData struct {
	connector hostingdomain.Connector
}
