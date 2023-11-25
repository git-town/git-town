package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/cli/log"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/hosting"
	"github.com/git-town/git-town/v10/src/hosting/github"
	"github.com/git-town/git-town/v10/src/validate"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/program"
	"github.com/git-town/git-town/v10/src/vm/runstate"
	"github.com/spf13/cobra"
)

const newPullRequestDesc = "Creates a new pull request"

const newPullRequestHelp = `
Syncs the current branch
and opens a browser window to the new pull request page of your repository.

The form is pre-populated for the current branch
so that the pull request only shows the changes made
against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, Gitea and Bitbucket.
When using self-hosted versions this command needs to be configured with
"git config %s <driver>"
where driver is "github", "gitlab", "gitea", or "bitbucket".
When using SSH identities, this command needs to be configured with
"git config %s <hostname>"
where hostname matches what is in your ssh config file.`

func newPullRequestCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   newPullRequestDesc,
		Long:    long(newPullRequestDesc, fmt.Sprintf(newPullRequestHelp, config.KeyCodeHostingDriver, config.KeyCodeHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeNewPullRequest(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeNewPullRequest(verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineNewPullRequestConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:             "new-pull-request",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          newPullRequestProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               config.connector,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              !config.pushHook,
	})
}

type newPullRequestConfig struct {
	branches           domain.Branches
	branchesToSync     domain.BranchInfos
	connector          hosting.Connector
	hasOpenChanges     bool
	remotes            domain.Remotes
	isOffline          bool
	lineage            config.Lineage
	mainBranch         domain.LocalBranchName
	previousBranch     domain.LocalBranchName
	pullBranchStrategy config.PullBranchStrategy
	pushHook           bool
	shouldSyncUpstream bool
	syncStrategy       config.SyncStrategy
}

func determineNewPullRequestConfig(repo *execute.OpenRepoResult, verbose bool) (*newPullRequestConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
	lineage := repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.Config.PushHook()
	if err != nil {
		return nil, domain.EmptyBranchesSnapshot(), domain.EmptyStashSnapshot(), false, err
	}
	branches, branchesSnapshot, stashSnapshot, exit, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSnapshot, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	updated, err := validate.KnowsBranchAncestors(branches.Initial, validate.KnowsBranchAncestorsArgs{
		AllBranches:   branches.All,
		Backend:       &repo.Runner.Backend,
		BranchTypes:   branches.Types,
		DefaultBranch: mainBranch,
		MainBranch:    mainBranch,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	if updated {
		lineage = repo.Runner.Config.Lineage(repo.Runner.Backend.Config.RemoveLocalConfigValue)
	}
	syncStrategy, err := repo.Runner.Config.SyncStrategy()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	pullBranchStrategy, err := repo.Runner.Config.PullBranchStrategy()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	shouldSyncUpstream, err := repo.Runner.Config.ShouldSyncUpstream()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	originURL := repo.Runner.Config.OriginURL()
	hostingService, err := repo.Runner.Config.HostingService()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.Config.GiteaToken(),
		GithubAPIToken:  github.GetAPIToken(repo.Runner.Config),
		GitlabAPIToken:  repo.Runner.Config.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             log.Printing{},
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	if connector == nil {
		return nil, branchesSnapshot, stashSnapshot, false, hosting.UnsupportedServiceError()
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync, err := branches.All.Select(branchNamesToSync)
	return &newPullRequestConfig{
		branches:           branches,
		branchesToSync:     branchesToSync,
		connector:          connector,
		hasOpenChanges:     repoStatus.OpenChanges,
		remotes:            remotes,
		isOffline:          repo.IsOffline,
		lineage:            lineage,
		mainBranch:         mainBranch,
		previousBranch:     previousBranch,
		pullBranchStrategy: pullBranchStrategy,
		pushHook:           pushHook,
		shouldSyncUpstream: shouldSyncUpstream,
		syncStrategy:       syncStrategy,
	}, branchesSnapshot, stashSnapshot, false, err
}

func newPullRequestProgram(config *newPullRequestConfig) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		syncBranchProgram(branch, syncBranchProgramArgs{
			branchTypes:        config.branches.Types,
			remotes:            config.remotes,
			isOffline:          config.isOffline,
			lineage:            config.lineage,
			program:            &prog,
			mainBranch:         config.mainBranch,
			pullBranchStrategy: config.pullBranchStrategy,
			pushBranch:         true,
			pushHook:           config.pushHook,
			shouldSyncUpstream: config.shouldSyncUpstream,
			syncStrategy:       config.syncStrategy,
		})
	}
	wrap(&prog, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	prog.Add(&opcode.CreateProposal{Branch: config.branches.Initial})
	return prog
}
