package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/cmd/sync"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err1 := readDryRunFlag(cmd)
			verbose, err2 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2); err != nil {
				return err
			}
			return executeDetach(args, dryRun, verbose)
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeDetach(args []string, dryRun configdomain.DryRun, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineDetachData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
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
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
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
		Connector:               data.connector,
		Detached:                true,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		PendingCommand:          None[string](),
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type detachData struct {
	branchInfosLastRun  Option[gitdomain.BranchInfos]
	branchToDetachInfo  gitdomain.BranchInfo
	branchToDetachName  gitdomain.LocalBranchName
	branchToDetachType  configdomain.BranchType
	branchesSnapshot    gitdomain.BranchesSnapshot
	children            []detachChildBranch
	config              config.ValidatedConfig
	connector           Option[forgedomain.Connector]
	descendents         []detachChildBranch
	dialogTestInputs    dialogcomponents.TestInputs
	dryRun              configdomain.DryRun
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	parentBranch        gitdomain.LocalBranchName
	previousBranch      Option[gitdomain.LocalBranchName]
	stashSize           gitdomain.StashSize
}

type detachChildBranch struct {
	info     *gitdomain.BranchInfo
	name     gitdomain.LocalBranchName
	proposal Option[forgedomain.Proposal]
}

func determineDetachData(args []string, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data detachData, exit dialogdomain.Exit, err error) {
	dialogTestInputs := dialogcomponents.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		CodebergToken:        config.CodebergToken,
		ForgeType:            config.ForgeType,
		Frontend:             repo.Frontend,
		GitHubConnectorType:  config.GitHubConnectorType,
		GitHubToken:          config.GitHubToken,
		GitLabConnectorType:  config.GitLabConnectorType,
		GitLabToken:          config.GitLabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Detached:              true,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchNameToDetach := gitdomain.NewLocalBranchName(slice.FirstElementOr(args, branchesSnapshot.Active.String()))
	branchToDetachInfo, hasBranchToDetachInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToDetach).Get()
	if !hasBranchToDetachInfo {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, branchNameToDetach)
	}
	if branchToDetachInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
		return data, exit, fmt.Errorf(messages.BranchOtherWorktree, branchNameToDetach)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchTypeToDetach := validatedConfig.BranchType(branchNameToDetach)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	previousBranchOpt := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	parentBranch, hasParentBranch := validatedConfig.NormalConfig.Lineage.Parent(branchNameToDetach).Get()
	if !hasParentBranch {
		return data, false, errors.New(messages.DetachNoParent)
	}
	branchHasMergeCommits, err := repo.Git.BranchContainsMerges(repo.Backend, branchNameToDetach, parentBranch)
	if err != nil {
		return data, false, err
	}
	if branchHasMergeCommits {
		return data, false, fmt.Errorf(messages.BranchContainsMergeCommits, branchNameToDetach)
	}
	childBranches := validatedConfig.NormalConfig.Lineage.Children(branchNameToDetach)
	children := make([]detachChildBranch, len(childBranches))
	for c, childBranch := range childBranches {
		proposal := None[forgedomain.Proposal]()
		if connector, hasConnector := connector.Get(); hasConnector {
			if findProposal, canFindProposal := connector.FindProposalFn().Get(); canFindProposal {
				proposal, err = findProposal(childBranch, initialBranch)
				if err != nil {
					return data, false, err
				}
			}
		}
		childInfo, has := branchesSnapshot.Branches.FindByLocalName(childBranch).Get()
		if !has {
			return data, false, fmt.Errorf("cannot find branch info for %q", childBranch)
		}
		children[c] = detachChildBranch{
			info:     childInfo,
			name:     childBranch,
			proposal: proposal,
		}
	}
	descendentNames := validatedConfig.NormalConfig.Lineage.Descendants(branchNameToDetach)
	descendents := make([]detachChildBranch, len(descendentNames))
	for d, descendentName := range descendentNames {
		info, has := branchesSnapshot.Branches.FindByLocalName(descendentName).Get()
		if !has {
			return data, false, fmt.Errorf("cannot find branch info for %q", descendentName)
		}
		descendents[d] = detachChildBranch{
			info:     info,
			name:     descendentName,
			proposal: None[forgedomain.Proposal](),
		}
	}
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	return detachData{
		branchInfosLastRun:  branchInfosLastRun,
		branchToDetachInfo:  *branchToDetachInfo,
		branchToDetachName:  branchNameToDetach,
		branchToDetachType:  branchTypeToDetach,
		branchesSnapshot:    branchesSnapshot,
		children:            children,
		config:              validatedConfig,
		connector:           connector,
		descendents:         descendents,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		nonExistingBranches: nonExistingBranches,
		parentBranch:        parentBranch,
		previousBranch:      previousBranchOpt,
		stashSize:           stashSize,
	}, false, nil
}

func detachProgram(repo execute.OpenRepoResult, data detachData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Frontend)
	prog.Value.Add(
		&opcodes.RebaseOntoRemoveDeleted{
			BranchToRebaseOnto: data.config.ValidatedConfigData.MainBranch,
			CommitsToRemove:    data.parentBranch.BranchName(),
			Upstream:           None[gitdomain.LocalBranchName](),
		},
	)
	if data.branchToDetachInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: data.branchToDetachName, ForceIfIncludes: true},
		)
	}
	lastParent := data.parentBranch
	for _, descendent := range data.descendents {
		sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
			Ancestor:          data.branchToDetachName.BranchName(),
			Branch:            descendent.name,
			HasTrackingBranch: descendent.info.HasTrackingBranch(),
			Program:           prog,
			RebaseOnto:        lastParent,
		})
		if descendent.info.HasTrackingBranch() {
			prog.Value.Add(
				&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: descendent.name, ForceIfIncludes: true},
			)
		}
		lastParent = descendent.name
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.initialBranch})
	if !data.dryRun {
		prog.Value.Add(
			&opcodes.LineageParentSet{
				Branch: data.branchToDetachName,
				Parent: data.config.ValidatedConfigData.MainBranch,
			},
		)
		for _, child := range data.config.NormalConfig.Lineage.Children(data.branchToDetachName) {
			prog.Value.Add(
				&opcodes.LineageParentSet{
					Branch: child,
					Parent: data.parentBranch,
				},
			)
		}
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
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
