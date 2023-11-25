package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/cli/log"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/execute"
	"github.com/git-town/git-town/v10/src/gohacks/slice"
	"github.com/git-town/git-town/v10/src/gohacks/stringslice"
	"github.com/git-town/git-town/v10/src/hosting"
	"github.com/git-town/git-town/v10/src/hosting/github"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/git-town/git-town/v10/src/validate"
	"github.com/git-town/git-town/v10/src/vm/interpreter"
	"github.com/git-town/git-town/v10/src/vm/opcode"
	"github.com/git-town/git-town/v10/src/vm/program"
	"github.com/git-town/git-town/v10/src/vm/runstate"
	"github.com/spf13/cobra"
)

const shipDesc = "Deliver a completed feature branch"

const shipHelp = `
Squash-merges the current branch, or <branch_name> if given,
into the main branch, resulting in linear history on the main branch.

- syncs the main branch
- pulls updates for <branch_name>
- merges the main branch into <branch_name>
- squash-merges <branch_name> into the main branch
  with commit message specified by the user
- pushes the main branch to the origin repository
- deletes <branch_name> from the local and origin repositories

Ships direct children of the main branch.
To ship a nested child branch, ship or kill all ancestor branches first.

If you use GitHub, this command can squash merge pull requests via the GitHub API. Setup:
1. Get a GitHub personal access token with the "repo" scope
2. Run 'git config %s <token>' (optionally add the '--global' flag)
Now anytime you ship a branch with a pull request on GitHub, it will squash merge via the GitHub API.
It will also update the base branch for any pull requests against that branch.

If your origin server deletes shipped branches, for example
GitHub's feature to automatically delete head branches,
run "git config %s false"
and Git Town will leave it up to your origin server to delete the remote branch.`

func shipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMessageFlag, readMessageFlag := flags.String("message", "m", "", "Specify the commit message for the squash commit")
	cmd := cobra.Command{
		Use:     "ship",
		GroupID: "basic",
		Args:    cobra.MaximumNArgs(1),
		Short:   shipDesc,
		Long:    long(shipDesc, fmt.Sprintf(shipHelp, config.KeyGithubToken, config.KeyShipDeleteRemoteBranch)),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeShip(args, readMessageFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	return &cmd
}

func executeShip(args []string, message string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSnapshot, exit, err := determineShipConfig(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	if config.branchToShip.LocalName == config.branches.Initial {
		repoStatus, err := repo.Runner.Backend.RepoStatus()
		if err != nil {
			return err
		}
		err = validate.NoOpenChanges(repoStatus.OpenChanges)
		if err != nil {
			return err
		}
	}
	runState := runstate.RunState{
		Command:             "ship",
		InitialActiveBranch: initialBranchesSnapshot.Active,
		RunProgram:          shipProgram(config, message),
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &runState,
		Run:                     &repo.Runner,
		Connector:               config.connector,
		Verbose:                 verbose,
		Lineage:                 config.lineage,
		NoPushHook:              !config.pushHook,
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnapshot,
	})
}

type shipConfig struct {
	branches                 domain.Branches
	branchToShip             domain.BranchInfo
	connector                hosting.Connector
	targetBranch             domain.BranchInfo
	canShipViaAPI            bool
	childBranches            domain.LocalBranchNames
	proposalMessage          string
	deleteOriginBranch       bool
	hasOpenChanges           bool
	remotes                  domain.Remotes
	isShippingInitialBranch  bool
	isOffline                bool
	lineage                  config.Lineage
	mainBranch               domain.LocalBranchName
	previousBranch           domain.LocalBranchName
	proposal                 *domain.Proposal
	proposalsOfChildBranches []domain.Proposal
	pullBranchStrategy       config.PullBranchStrategy
	pushHook                 bool
	shouldSyncUpstream       bool
	syncStrategy             config.SyncStrategy
}

func determineShipConfig(args []string, repo *execute.OpenRepoResult, verbose bool) (*shipConfig, domain.BranchesSnapshot, domain.StashSnapshot, bool, error) {
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
		ValidateNoOpenChanges: len(args) == 0,
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
	deleteOrigin, err := repo.Runner.Config.ShouldShipDeleteOriginBranch()
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	mainBranch := repo.Runner.Config.MainBranch()
	branchNameToShip := domain.NewLocalBranchName(slice.FirstElementOr(args, branches.Initial.String()))
	branchToShip := branches.All.FindByLocalName(branchNameToShip)
	isShippingInitialBranch := branchNameToShip == branches.Initial
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
	if !isShippingInitialBranch {
		if branchToShip == nil {
			return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToShip)
		}
	}
	if !branches.Types.IsFeatureBranch(branchNameToShip) {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.ShipNoFeatureBranch, branchNameToShip)
	}
	branches.Types, lineage, err = execute.EnsureKnownBranchAncestry(branchNameToShip, execute.EnsureKnownBranchAncestryArgs{
		AllBranches:   branches.All,
		BranchTypes:   branches.Types,
		DefaultBranch: mainBranch,
		Lineage:       lineage,
		MainBranch:    mainBranch,
		Runner:        &repo.Runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	err = ensureParentBranchIsMainOrPerennialBranch(branchNameToShip, branches.Types, lineage)
	if err != nil {
		return nil, branchesSnapshot, stashSnapshot, false, err
	}
	targetBranchName := lineage.Parent(branchNameToShip)
	targetBranch := branches.All.FindByLocalName(targetBranchName)
	if targetBranch == nil {
		return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.BranchDoesntExist, targetBranchName)
	}
	canShipViaAPI := false
	proposalMessage := ""
	var proposal *domain.Proposal
	childBranches := lineage.Children(branchNameToShip)
	proposalsOfChildBranches := []domain.Proposal{}
	pushHook, err = repo.Runner.Config.PushHook()
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
	if !repo.IsOffline && connector != nil {
		if branchToShip.HasTrackingBranch() {
			proposal, err = connector.FindProposal(branchNameToShip, targetBranchName)
			if err != nil {
				return nil, branchesSnapshot, stashSnapshot, false, err
			}
			if proposal != nil {
				canShipViaAPI = true
				proposalMessage = connector.DefaultProposalMessage(*proposal)
			}
		}
		for _, childBranch := range childBranches {
			childProposal, err := connector.FindProposal(childBranch, branchNameToShip)
			if err != nil {
				return nil, branchesSnapshot, stashSnapshot, false, fmt.Errorf(messages.ProposalNotFoundForBranch, branchNameToShip, err)
			}
			if childProposal != nil {
				proposalsOfChildBranches = append(proposalsOfChildBranches, *childProposal)
			}
		}
	}
	return &shipConfig{
		branches:                 branches,
		connector:                connector,
		targetBranch:             *targetBranch,
		branchToShip:             *branchToShip,
		canShipViaAPI:            canShipViaAPI,
		childBranches:            childBranches,
		proposalMessage:          proposalMessage,
		deleteOriginBranch:       deleteOrigin,
		hasOpenChanges:           repoStatus.OpenChanges,
		remotes:                  remotes,
		isOffline:                repo.IsOffline,
		isShippingInitialBranch:  isShippingInitialBranch,
		lineage:                  lineage,
		mainBranch:               mainBranch,
		previousBranch:           previousBranch,
		proposal:                 proposal,
		proposalsOfChildBranches: proposalsOfChildBranches,
		pullBranchStrategy:       pullBranchStrategy,
		pushHook:                 pushHook,
		shouldSyncUpstream:       shouldSyncUpstream,
		syncStrategy:             syncStrategy,
	}, branchesSnapshot, stashSnapshot, false, nil
}

func ensureParentBranchIsMainOrPerennialBranch(branch domain.LocalBranchName, branchTypes domain.BranchTypes, lineage config.Lineage) error {
	parentBranch := lineage.Parent(branch)
	if !branchTypes.IsMainBranch(parentBranch) && !branchTypes.IsPerennialBranch(parentBranch) {
		ancestors := lineage.Ancestors(branch)
		ancestorsWithoutMainOrPerennial := ancestors[1:]
		oldestAncestor := ancestorsWithoutMainOrPerennial[0]
		return fmt.Errorf(`shipping this branch would ship %s as well,
please ship %q first`, stringslice.Connect(ancestorsWithoutMainOrPerennial.Strings()), oldestAncestor)
	}
	return nil
}

func shipProgram(config *shipConfig, commitMessage string) program.Program {
	prog := program.Program{}
	// sync the parent branch
	syncBranchProgram(config.targetBranch, syncBranchProgramArgs{
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
	// sync the branch to ship (local sync only)
	syncBranchProgram(config.branchToShip, syncBranchProgramArgs{
		branchTypes:        config.branches.Types,
		remotes:            config.remotes,
		isOffline:          config.isOffline,
		lineage:            config.lineage,
		program:            &prog,
		mainBranch:         config.mainBranch,
		pullBranchStrategy: config.pullBranchStrategy,
		pushBranch:         false,
		pushHook:           config.pushHook,
		shouldSyncUpstream: config.shouldSyncUpstream,
		syncStrategy:       config.syncStrategy,
	})
	prog.Add(&opcode.EnsureHasShippableChanges{Branch: config.branchToShip.LocalName, Parent: config.mainBranch})
	prog.Add(&opcode.Checkout{Branch: config.targetBranch.LocalName})
	if config.canShipViaAPI {
		// update the proposals of child branches
		for _, childProposal := range config.proposalsOfChildBranches {
			prog.Add(&opcode.UpdateProposalTarget{
				ProposalNumber: childProposal.Number,
				NewTarget:      config.targetBranch.LocalName,
			})
		}
		// push
		prog.Add(&opcode.PushCurrentBranch{CurrentBranch: config.branchToShip.LocalName, NoPushHook: !config.pushHook})
		prog.Add(&opcode.ConnectorMergeProposal{
			Branch:          config.branchToShip.LocalName,
			ProposalNumber:  config.proposal.Number,
			CommitMessage:   commitMessage,
			ProposalMessage: config.proposalMessage,
		})
		prog.Add(&opcode.PullCurrentBranch{})
	} else {
		prog.Add(&opcode.SquashMerge{Branch: config.branchToShip.LocalName, CommitMessage: commitMessage, Parent: config.targetBranch.LocalName})
	}
	if config.remotes.HasOrigin() && !config.isOffline {
		prog.Add(&opcode.PushCurrentBranch{CurrentBranch: config.targetBranch.LocalName, NoPushHook: !config.pushHook})
	}
	// NOTE: when shipping via API, we can always delete the remote branch because:
	// - we know we have a tracking branch (otherwise there would be no PR to ship via API)
	// - we have updated the PRs of all child branches (because we have API access)
	// - we know we are online
	if config.canShipViaAPI || (config.branchToShip.HasTrackingBranch() && len(config.childBranches) == 0 && !config.isOffline) {
		if config.deleteOriginBranch {
			prog.Add(&opcode.DeleteTrackingBranch{Branch: config.branchToShip.RemoteName})
		}
	}
	prog.Add(&opcode.DeleteLocalBranch{Branch: config.branchToShip.LocalName, Force: false})
	prog.Add(&opcode.DeleteParentBranch{Branch: config.branchToShip.LocalName})
	for _, child := range config.childBranches {
		prog.Add(&opcode.ChangeParent{Branch: child, Parent: config.targetBranch.LocalName})
	}
	if !config.isShippingInitialBranch {
		prog.Add(&opcode.Checkout{Branch: config.branches.Initial})
	}
	wrap(&prog, wrapOptions{
		RunInGitRoot:     true,
		StashOpenChanges: !config.isShippingInitialBranch && config.hasOpenChanges,
		MainBranch:       config.mainBranch,
		InitialBranch:    config.branches.Initial,
		PreviousBranch:   config.previousBranch,
	})
	return prog
}
