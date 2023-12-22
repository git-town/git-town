package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/log"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/hosting"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v11/src/sync/syncprograms"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/spf13/cobra"
)

const proposeDesc = "Creates a proposal to merge a Git branch"

const proposeHelp = `
Syncs the current branch
and opens a browser window to the new proposal page of your repository.

The form is pre-populated for the current branch
so that the proposal only shows the changes made
against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, Gitea and Bitbucket.
When using self-hosted versions this command needs to be configured with
"git config %s <driver>"
where driver is "github", "gitlab", "gitea", or "bitbucket".
When using SSH identities, this command needs to be configured with
"git config %s <hostname>"
where hostname matches what is in your ssh config file.`

func proposeCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "propose",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, configdomain.KeyCodeHostingPlatform, configdomain.KeyCodeHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePropose(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePropose(verbose bool) error {
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
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineProposeConfig(repo, verbose)
	if err != nil || exit {
		return err
	}
	if err != nil {
		return err
	}
	runState := runstate.RunState{
		Command:             "propose",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          proposeProgram(config),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     repo.Runner,
		Connector:               config.connector,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
		NoPushHook:              config.pushHook.Negate(),
	})
}

type proposeConfig struct {
	branches              configdomain.Branches
	branchesToSync        gitdomain.BranchInfos
	connector             hostingdomain.Connector
	hasOpenChanges        bool
	remotes               gitdomain.Remotes
	isOnline              configdomain.Online
	lineage               configdomain.Lineage
	mainBranch            gitdomain.LocalBranchName
	previousBranch        gitdomain.LocalBranchName
	syncPerennialStrategy configdomain.SyncPerennialStrategy
	pushHook              configdomain.PushHook
	syncUpstream          configdomain.SyncUpstream
	syncFeatureStrategy   configdomain.SyncFeatureStrategy
}

func determineProposeConfig(repo *execute.OpenRepoResult, verbose bool) (*proposeConfig, undodomain.BranchesSnapshot, undodomain.StashSnapshot, bool, error) {
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	pushHook := repo.Runner.GitTown.PushHook
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
	mainBranch := repo.Runner.GitTown.MainBranch
	branches.Types, lineage, err = execute.EnsureKnownBranchAncestry(branches.Initial, execute.EnsureKnownBranchAncestryArgs{
		AllBranches:   branches.All,
		BranchTypes:   branches.Types,
		DefaultBranch: mainBranch,
		Lineage:       lineage,
		MainBranch:    mainBranch,
		Runner:        repo.Runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	syncFeatureStrategy := repo.Runner.GitTown.SyncFeatureStrategy
	syncPerennialStrategy := repo.Runner.GitTown.SyncPerennialStrategy
	syncUpstream := repo.Runner.GitTown.SyncUpstream
	originURL := repo.Runner.GitTown.OriginURL()
	hostingService, err := repo.Runner.GitTown.HostingService()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.GitTown.GiteaToken,
		GithubAPIToken:  github.GetAPIToken(repo.Runner.GitTown.GitHubToken),
		GitlabAPIToken:  repo.Runner.GitTown.GitLabToken,
		MainBranch:      mainBranch,
		Log:             log.Printing{},
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	if connector == nil {
		return nil, branchesSnapshot, stashSnapshot, false, hostingdomain.UnsupportedServiceError()
	}
	branchNamesToSync := lineage.BranchAndAncestors(branches.Initial)
	branchesToSync, err := branches.All.Select(branchNamesToSync)
	return &proposeConfig{
		branches:              branches,
		branchesToSync:        branchesToSync,
		connector:             connector,
		hasOpenChanges:        repoStatus.OpenChanges,
		remotes:               remotes,
		isOnline:              repo.IsOffline.ToOnline(),
		lineage:               lineage,
		mainBranch:            mainBranch,
		previousBranch:        previousBranch,
		syncPerennialStrategy: syncPerennialStrategy,
		pushHook:              pushHook,
		syncUpstream:          syncUpstream,
		syncFeatureStrategy:   syncFeatureStrategy,
	}, branchesSnapshot, stashSnapshot, false, err
}

func proposeProgram(config *proposeConfig) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		syncprograms.SyncBranchProgram(branch, syncprograms.SyncBranchProgramArgs{
			BranchInfos:           config.branches.All,
			BranchTypes:           config.branches.Types,
			Remotes:               config.remotes,
			IsOnline:              config.isOnline,
			Lineage:               config.lineage,
			Program:               &prog,
			MainBranch:            config.mainBranch,
			SyncPerennialStrategy: config.syncPerennialStrategy,
			PushBranch:            true,
			PushHook:              config.pushHook,
			SyncUpstream:          config.syncUpstream,
			SyncFeatureStrategy:   config.syncFeatureStrategy,
		})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})
	prog.Add(&opcode.CreateProposal{Branch: config.branches.Initial})
	return prog
}
