package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/sync"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/programs"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	detachCommandName = "detach"
	detachDesc        = "Move a branch out of a stack"
	detachHelp        = `
The "detach" command removes the current branch from the stack it is in
and makes it a stand-alone top-level branch
that ships directly into your main branch.
This is useful when a branch in a stack makes changes
that are independent from the changes made by other branches in this stack.
Detaching such independent branches
reduces your stack to changes that belong together,
and gets more of your changes reviewed and shipped concurrently.

Consider this stack:

main
 \
  branch-1
   \
*   branch-2
     \
      branch-3

We are on the "branch-2" branch.
After running "git town detach",
we end up with this stack:

main
 \
  branch-1
   \
    branch-3
 \
* branch-2
`
)

func detachCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     detachCommandName,
		Args:    cobra.NoArgs,
		Short:   detachDesc,
		GroupID: cmdhelpers.GroupIDStack,
		Long:    cmdhelpers.Long(detachDesc, detachHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			dryRun, errDryRun := readDryRunFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDryRun, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          Some(configdomain.Detached(true)),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeDetach(cliConfig)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeDetach(cliConfig configdomain.PartialConfig) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineDetachData(repo)
	if err != nil {
		return err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit:
		return nil
	case configdomain.ProgramFlowRestart:
		goto Start
	}
	err = validateDetachData(data)
	if err != nil {
		return err
	}
	runProgram := detachProgram(repo, data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               detachCommandName,
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		BranchInfosLastRun:    data.branchInfosLastRun,
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		ConfigDir:               repo.ConfigDir,
		Connector:               data.connector,
		DryRun:                  data.config.NormalConfig.DryRun,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		Inputs:                  data.inputs,
		PendingCommand:          None[string](),
		RunState:                runState,
	})
}

type detachData struct {
	branchInfosLastRun     Option[gitdomain.BranchInfos]
	branchToDetachInfo     gitdomain.BranchInfo
	branchToDetachName     gitdomain.LocalBranchName
	branchToDetachProposal Option[forgedomain.Proposal]
	branchToDetachType     configdomain.BranchType
	branchesSnapshot       gitdomain.BranchesSnapshot
	children               []detachChildBranch
	config                 config.ValidatedConfig
	connector              Option[forgedomain.Connector]
	descendents            []detachChildBranch
	hasOpenChanges         bool
	initialBranch          gitdomain.LocalBranchName
	inputs                 dialogcomponents.Inputs
	nonExistingBranches    gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch           gitdomain.LocalBranchName
	previousBranch         Option[gitdomain.LocalBranchName]
	stashSize              gitdomain.StashSize
}

type detachChildBranch struct {
	info     *gitdomain.BranchInfo
	name     gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func determineDetachData(repo execute.OpenRepoResult) (detachData, configdomain.ProgramFlow, error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	var emptyResult detachData
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ConfigDir:            repo.ConfigDir,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
		Frontend:             repo.Frontend,
		GiteaToken:           config.GiteaToken,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Inputs:                inputs,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return emptyResult, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.DetachRepoHasDetachedHead)
	}
	currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchNameToDetach := currentBranch
	branchToDetachInfo, hasBranchToDetachInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToDetach).Get()
	if !hasBranchToDetachInfo {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchDoesntExist, branchNameToDetach)
	}
	if branchToDetachInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchOtherWorktree, branchNameToDetach)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		ConfigDir:          repo.ConfigDir,
		ConfigSnapshot:     repo.ConfigSnapshot,
		Connector:          connector,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		Inputs:             inputs,
		LocalBranches:      localBranches,
		Remotes:            remotes,
		RepoStatus:         repoStatus,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	branchTypeToDetach := validatedConfig.BranchType(branchNameToDetach)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(branchNameToDetach).Get()
	if !hasParentBranch {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.DetachNoParent)
	}
	branchHasMergeCommits, err := repo.Git.BranchContainsMerges(repo.Backend, branchNameToDetach, parentBranch)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	if branchHasMergeCommits {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchContainsMergeCommits, branchNameToDetach)
	}

	connectorProposalFinder := None[forgedomain.ProposalFinder]()
	branchToDetachProposal := None[forgedomain.Proposal]()
	if connector, hasConnector := connector.Get(); hasConnector {
		if proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder); canFindProposals {
			connectorProposalFinder = Some(proposalFinder)
			branchToDetachProposal, err = proposalFinder.FindProposal(branchNameToDetach, parentBranch)
			if err != nil {
				return emptyResult, configdomain.ProgramFlowExit, err
			}
		}
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(branchNameToDetach, validatedConfig.NormalConfig.Order)
	children := make([]detachChildBranch, len(childBranches))
	proposalFinder, hasProposalFinder := connectorProposalFinder.Get()
	for c, childBranch := range childBranches {
		proposal := None[forgedomain.Proposal]()
		if hasProposalFinder {
			proposal, err = proposalFinder.FindProposal(childBranch, initialBranch)
			if err != nil {
				return emptyResult, configdomain.ProgramFlowExit, err
			}
		}
		childInfo, has := branchesSnapshot.Branches.FindByLocalName(childBranch).Get()
		if !has {
			return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf("cannot find branch info for %q", childBranch)
		}
		children[c] = detachChildBranch{
			info:     childInfo,
			name:     childBranch,
			proposal: proposal,
		}
	}
	descendentNames := validatedConfig.NormalConfig.Lineage.Descendants(branchNameToDetach, validatedConfig.NormalConfig.Order)
	descendents := make([]detachChildBranch, len(descendentNames))
	for d, descendentName := range descendentNames {
		info, has := branchesSnapshot.Branches.FindByLocalName(descendentName).Get()
		if !has {
			return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf("cannot find branch info for %q", descendentName)
		}
		descendents[d] = detachChildBranch{
			info:     info,
			name:     descendentName,
			proposal: None[forgedomain.Proposal](),
		}
	}
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(lineageBranches...)
	return detachData{
		branchInfosLastRun:     branchInfosLastRun,
		branchToDetachInfo:     *branchToDetachInfo,
		branchToDetachName:     branchNameToDetach,
		branchToDetachProposal: branchToDetachProposal,
		branchToDetachType:     branchTypeToDetach,
		branchesSnapshot:       branchesSnapshot,
		children:               children,
		config:                 validatedConfig,
		connector:              connector,
		descendents:            descendents,
		hasOpenChanges:         repoStatus.OpenChanges,
		initialBranch:          initialBranch,
		inputs:                 inputs,
		nonExistingBranches:    nonExistingBranches,
		parentBranch:           parentBranch,
		previousBranch:         previousBranchOpt,
		stashSize:              stashSize,
	}, configdomain.ProgramFlowContinue, nil
}

func detachProgram(repo execute.OpenRepoResult, data detachData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Frontend, data.config.NormalConfig.Order)
	oldParent := data.parentBranch
	// step 1: delete the commits of the branch to detach from all descendents,
	// while that branch is still in the form it had inside the stack
	for _, descendent := range data.descendents {
		// Determine the correct parent to rebase onto.
		// If the descendant's parent is the branch being detached, rebase onto that branch's parent.
		// Otherwise, rebase onto the descendant's actual parent.
		rebaseOnto := data.parentBranch
		if descendentParent, hasParent := data.config.NormalConfig.Lineage.Parent(descendent.name).Get(); hasParent {
			if descendentParent != data.branchToDetachName {
				rebaseOnto = descendentParent
			}
		}
		sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
			Ancestor:          data.branchToDetachName.BranchName(),
			Branch:            descendent.name,
			HasTrackingBranch: descendent.info.HasTrackingBranch(),
			Program:           prog,
			RebaseOnto:        rebaseOnto,
		})
		if descendentTracking, descendentHasTracking := descendent.info.RemoteName.Get(); descendentHasTracking {
			prog.Value.Add(
				&opcodes.PushCurrentBranchForceIfNeeded{
					CurrentBranch:   descendent.name,
					ForceIfIncludes: true,
					TrackingBranch:  descendentTracking,
				},
			)
		}
	}
	// step 2: delete the commits of parent branches from the detached branch
	prog.Value.Add(
		&opcodes.CheckoutIfNeeded{
			Branch: data.branchToDetachName,
		},
		&opcodes.RebaseOnto{
			BranchToRebaseOnto: data.config.ValidatedConfigData.MainBranch.BranchName(),
			CommitsToRemove:    data.parentBranch.BranchName().Location(),
		},
	)
	if trackingBranch, hasTrackingBranch := data.branchToDetachInfo.RemoteName.Get(); hasTrackingBranch {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{
				CurrentBranch:   data.branchToDetachName,
				ForceIfIncludes: true,
				TrackingBranch:  trackingBranch,
			},
		)
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.initialBranch})
	if !data.config.NormalConfig.DryRun {
		prog.Value.Add(
			&opcodes.LineageParentSet{
				Branch: data.branchToDetachName,
				Parent: data.config.ValidatedConfigData.MainBranch,
			},
		)
	}
	if proposal, hasProposal := data.branchToDetachProposal.Get(); hasProposal {
		prog.Value.Add(
			&opcodes.ProposalUpdateTarget{
				NewBranch: data.config.ValidatedConfigData.MainBranch,
				OldBranch: data.parentBranch,
				Proposal:  proposal,
			},
		)
	}
	for _, child := range data.children {
		prog.Value.Add(
			&opcodes.LineageParentSet{
				Branch: child.name,
				Parent: data.parentBranch,
			},
		)

		if proposal, hasProposal := child.proposal.Get(); hasProposal {
			prog.Value.Add(
				&opcodes.ProposalUpdateTarget{
					NewBranch: data.parentBranch,
					OldBranch: data.branchToDetachName,
					Proposal:  proposal,
				},
			)
		}
	}
	updateBreadcrumb := data.config.NormalConfig.ProposalBreadcrumb.Enabled()
	isOnline := data.config.NormalConfig.Offline.IsOnline()
	if updateBreadcrumb && isOnline {
		programs.UpdateBreadcrumbsProgram(programs.UpdateBreadcrumbsArgs{
			Config:          data.config,
			Program:         prog,
			TouchedBranches: gitdomain.LocalBranchNames{data.branchToDetachName, oldParent},
		})
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return prog.Immutable()
}

func validateDetachData(data detachData) error {
	switch data.branchToDetachInfo.SyncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusAhead, gitdomain.SyncStatusLocalOnly:
	case gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusBehind:
		return errors.New(messages.DetachNeedsSync)
	case gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.DetachOtherWorkTree, data.branchToDetachName)
	case gitdomain.SyncStatusRemoteOnly:
		return errors.New(messages.DetachRemoteBranch)
	}
	switch data.branchToDetachType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return fmt.Errorf(messages.DetachUnsupportedBranchType, data.branchToDetachType)
	}
	for _, child := range data.children {
		switch child.info.SyncStatus {
		case
			gitdomain.SyncStatusAhead,
			gitdomain.SyncStatusLocalOnly,
			gitdomain.SyncStatusUpToDate:
		case
			gitdomain.SyncStatusBehind,
			gitdomain.SyncStatusDeletedAtRemote,
			gitdomain.SyncStatusNotInSync,
			gitdomain.SyncStatusRemoteOnly:
			return errors.New(messages.DetachNeedsSync)
		case gitdomain.SyncStatusOtherWorktree:
			return fmt.Errorf(messages.DetachOtherWorkTree, child.name)
		}
	}
	return nil
}
