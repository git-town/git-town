package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/ship"
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
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/optimizer"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
	"github.com/spf13/cobra"
)

const (
	prependDesc = "Create a new feature branch as the parent of the current branch"
	prependHelp = `
Syncs the parent branch,
cuts a new feature branch with the given name off the parent branch,
makes the new branch the parent of the current branch,
pushes the new feature branch to the origin repository
(if "share-new-branches" is "push"),
and brings over all uncommitted changes to the new feature branch.

See "sync" for upstream remote options.

Consider this stack:

main
 \
* feature-2

We are on the "feature-2" branch.
After running "git town prepend feature-1",
our repository has this stack:

main
 \
* feature-1
   \
    feature-2
`
)

func prependCommand() *cobra.Command {
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addBeamFlag, readBeamFlag := flags.Beam()
	addBodyFlag, readBodyFlag := flags.ProposalBody("")
	addCommitFlag, readCommitFlag := flags.Commit()
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addPushFlag, readPushFlag := flags.Push()
	addStashFlag, readStashFlag := flags.Stash()
	addSyncFlag, readSyncFlag := flags.Sync()
	addTitleFlag, readTitleFlag := flags.ProposalTitle()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "prepend <branch>",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.ExactArgs(1),
		Short:   prependDesc,
		Long:    cmdhelpers.Long(prependDesc, prependHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			beam, errBeam := readBeamFlag(cmd)
			bodyText, errBodyText := readBodyFlag(cmd)
			commit, errCommit := readCommitFlag(cmd)
			commitMessage, errCommitMessage := readCommitMessageFlag(cmd)
			detached, errDetached := readDetachedFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			propose, errPropose := readProposeFlag(cmd)
			prototype, errPrototype := readPrototypeFlag(cmd)
			push, errPush := readPushFlag(cmd)
			stash, errStash := readStashFlag(cmd)
			sync, errSync := readSyncFlag(cmd)
			title, errTitle := readTitleFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAutoResolve, errBeam, errBodyText, errCommit, errCommitMessage, errDetached, errDryRun, errPropose, errPrototype, errPush, errStash, errSync, errTitle, errVerbose); err != nil {
				return err
			}
			if commitMessage.IsSome() {
				commit = true
			}
			if propose.ShouldPropose() && !beam.ShouldBeam() {
				commit = true
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       autoResolve,
				AutoSync:          sync,
				Detached:          detached,
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      push,
				Stash:             stash,
				Verbose:           verbose,
			})
			return executePrepend(prependArgs{
				argv:          args,
				beam:          beam,
				cliConfig:     cliConfig,
				commit:        commit,
				commitMessage: commitMessage,
				proposalBody:  bodyText,
				proposalTitle: title,
				propose:       propose,
				prototype:     prototype,
			})
		},
	}
	addAutoResolveFlag(&cmd)
	addBeamFlag(&cmd)
	addBodyFlag(&cmd)
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addPushFlag(&cmd)
	addStashFlag(&cmd)
	addSyncFlag(&cmd)
	addTitleFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type prependArgs struct {
	argv          []string
	beam          configdomain.Beam
	cliConfig     configdomain.PartialConfig
	commit        configdomain.Commit
	commitMessage Option[gitdomain.CommitMessage]
	proposalBody  Option[gitdomain.ProposalBody]
	proposalTitle Option[gitdomain.ProposalTitle]
	propose       configdomain.Propose
	prototype     configdomain.Prototype
}

func executePrepend(args prependArgs) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determinePrependData(args, repo)
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
	runProgram := prependProgram(repo, data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               "prepend",
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
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
		RootDir:                 repo.RootDir,
		RunState:                runState,
	})
}

type prependData struct {
	beam                configdomain.Beam
	branchInfos         gitdomain.BranchInfos
	branchInfosLastRun  Option[gitdomain.BranchInfos]
	branchesSnapshot    gitdomain.BranchesSnapshot
	branchesToSync      configdomain.BranchesToSync
	commit              configdomain.Commit
	commitMessage       Option[gitdomain.CommitMessage]
	commitsToBeam       gitdomain.Commits
	config              config.ValidatedConfig
	connector           Option[forgedomain.Connector]
	existingParent      gitdomain.LocalBranchName
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	initialBranchInfo   gitdomain.BranchInfo
	inputs              dialogcomponents.Inputs
	newParentCandidates gitdomain.LocalBranchNames
	nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos gitdomain.BranchInfos
	previousBranch      Option[gitdomain.LocalBranchName]
	proposal            Option[forgedomain.Proposal]
	proposalBody        Option[gitdomain.ProposalBody]
	proposalTitle       Option[gitdomain.ProposalTitle]
	propose             configdomain.Propose
	prototype           configdomain.Prototype
	remotes             gitdomain.Remotes
	stashSize           gitdomain.StashSize
	targetBranch        gitdomain.LocalBranchName
}

func determinePrependData(args prependArgs, repo execute.OpenRepoResult) (data prependData, flow configdomain.ProgramFlow, err error) {
	prefetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
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
		return data, configdomain.ProgramFlowExit, err
	}
	fetch := true
	if repoStatus.OpenChanges {
		fetch = false
	}
	if args.beam.ShouldBeam() || args.commit.ShouldCommit() {
		fetch = false
	}
	if !config.AutoSync.ShouldSync() {
		fetch = false
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 fetch,
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
		return data, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return data, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return data, configdomain.ProgramFlowExit, errors.New(messages.PrependDetachedHead)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	targetBranch := gitdomain.NewLocalBranchName(args.argv[0])
	if prefix, hasPrefix := config.BranchPrefix.Get(); hasPrefix {
		targetBranch = prefix.Apply(targetBranch)
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
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
		return data, configdomain.ProgramFlowExit, err
	}
	ancestorOpt := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage)
	ancestor, hasAncestor := ancestorOpt.Get()
	if !hasAncestor {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.SetParentNoFeatureBranch, initialBranch)
	}
	commitsToBeam := []gitdomain.Commit{}
	if args.beam {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor.BranchName())
		if err != nil {
			return data, configdomain.ProgramFlowExit, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, targetBranch, repo.Git, repo.Backend, inputs)
		if err != nil || exit {
			return data, configdomain.ProgramFlowExit, err
		}
	}
	branchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	if repo.UnvalidatedConfig.NormalConfig.Detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	parentAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(ancestor)
	slices.Reverse(parentAndAncestors)
	proposalOpt := None[forgedomain.Proposal]()
	if !repo.IsOffline {
		proposalOpt = ship.FindProposal(connector, initialBranch, Some(ancestor))
	}
	propose := args.propose
	if validatedConfig.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPropose {
		propose = true
	}
	return prependData{
		beam:                args.beam,
		branchInfos:         branchesSnapshot.Branches,
		branchInfosLastRun:  branchInfosLastRun,
		branchesSnapshot:    branchesSnapshot,
		branchesToSync:      branchesToSync,
		commit:              args.commit,
		commitMessage:       args.commitMessage,
		commitsToBeam:       commitsToBeam,
		config:              validatedConfig,
		connector:           connector,
		existingParent:      ancestor,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		initialBranchInfo:   *initialBranchInfo,
		inputs:              inputs,
		newParentCandidates: parentAndAncestors,
		nonExistingBranches: nonExistingBranches,
		preFetchBranchInfos: prefetchBranchSnapshot.Branches,
		previousBranch:      previousBranch,
		proposal:            proposalOpt,
		proposalBody:        args.proposalBody,
		proposalTitle:       args.proposalTitle,
		propose:             propose,
		prototype:           args.prototype,
		remotes:             remotes,
		stashSize:           stashSize,
		targetBranch:        targetBranch,
	}, configdomain.ProgramFlowContinue, nil
}

func prependProgram(repo execute.OpenRepoResult, data prependData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	if !data.hasOpenChanges && !data.beam.ShouldBeam() && !data.commit.ShouldCommit() && data.config.NormalConfig.AutoSync.ShouldSync() {
		data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, finalMessages, repo.Backend, data.config.NormalConfig.Order)
		branchesToDelete := set.New[gitdomain.LocalBranchName]()
		sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
			BranchInfos:         data.branchInfos,
			BranchInfosPrevious: data.branchInfosLastRun,
			BranchesToDelete:    NewMutable(&branchesToDelete),
			Config:              data.config,
			InitialBranch:       data.initialBranch,
			PrefetchBranchInfos: data.preFetchBranchInfos,
			Program:             prog,
			Prune:               false,
			PushBranches:        data.config.NormalConfig.PushBranches,
			Remotes:             data.remotes,
		})
	}
	prog.Value.Add(&opcodes.BranchCreateAndCheckoutExistingParent{
		Ancestors: data.newParentCandidates,
		Branch:    data.targetBranch,
	})
	// set the parent of the newly created branch
	prog.Value.Add(&opcodes.LineageParentSetFirstExisting{
		Branch:    data.targetBranch,
		Ancestors: data.newParentCandidates,
	})
	// set the parent of the branch prepended to
	prog.Value.Add(&opcodes.LineageParentSetIfExists{
		Branch: data.initialBranch,
		Parent: data.targetBranch,
	})
	if data.prototype {
		prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: configdomain.BranchTypePrototypeBranch})
	} else if newBranchType, hasNewBranchType := data.config.NormalConfig.NewBranchType.Get(); hasNewBranchType {
		prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: newBranchType.BranchType()})
	}
	proposal, hasProposal := data.proposal.Get()
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.config.NormalConfig.Offline.IsOnline() && (data.config.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPush || hasProposal) {
		prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.targetBranch})
	}
	connector, hasConnector := data.connector.Get()
	_, canUpdateProposalTarget := connector.(forgedomain.ProposalTargetUpdater)
	if hasProposal && hasConnector && canUpdateProposalTarget {
		prog.Value.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: data.targetBranch,
			OldBranch: data.existingParent,
			Proposal:  proposal,
		})
	}
	moveCommitsToPrependedBranch(prog, data)
	if data.commit {
		prog.Value.Add(
			&opcodes.Commit{
				AuthorOverride:                 None[gitdomain.Author](),
				FallbackToDefaultCommitMessage: false,
				Message:                        data.commitMessage,
			},
		)
	}
	if data.propose {
		prog.Value.Add(
			&opcodes.BranchTrackingCreate{
				Branch: data.targetBranch,
			},
			&opcodes.ProposalCreate{
				Branch:        data.targetBranch,
				MainBranch:    data.config.ValidatedConfigData.MainBranch,
				ProposalBody:  data.proposalBody,
				ProposalTitle: data.proposalTitle,
			})
	}
	if data.commit {
		prog.Value.Add(
			&opcodes.Checkout{Branch: data.initialBranch},
		)
	} else {
		previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
		cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
			DryRun:                   repo.UnvalidatedConfig.NormalConfig.DryRun,
			InitialStashSize:         data.stashSize,
			RunInGitRoot:             true,
			StashOpenChanges:         data.hasOpenChanges && data.config.NormalConfig.Stash.ShouldStash(),
			PreviousBranchCandidates: previousBranchCandidates,
		})
	}
	if data.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		sync.AddSyncProposalsProgram(sync.AddSyncProposalsProgramArgs{
			ChangedBranches: gitdomain.LocalBranchNames{data.initialBranch},
			Config:          data.config,
			Program:         prog,
		})
	}
	return optimizer.Optimize(prog.Immutable())
}

// provides the strategy to use to sync a branch after beaming some of its commits to its new parent branch
func afterBeamToParentSyncStrategy(branchType configdomain.BranchType, config config.NormalConfig) Option[configdomain.SyncStrategy] {
	switch branchType {
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return None[configdomain.SyncStrategy]()
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch:
		return Some(config.SyncFeatureStrategy.SyncStrategy())
	case configdomain.BranchTypePrototypeBranch:
		return Some(config.SyncPrototypeStrategy.SyncStrategy())
	}
	panic("unhandled branch type: " + branchType.String())
}

// provides the name of the youngest ancestor branch of the given branch that actually exists,
// either locally or remotely.
func latestExistingAncestor(branch gitdomain.LocalBranchName, branchInfos gitdomain.BranchInfos, lineage configdomain.Lineage) Option[gitdomain.LocalBranchName] {
	for {
		parent, hasParent := lineage.Parent(branch).Get()
		if !hasParent {
			return None[gitdomain.LocalBranchName]()
		}
		if branchInfos.HasBranch(parent) {
			return Some(parent)
		}
		branch = parent
	}
}

func moveCommitsToPrependedBranch(prog Mutable[program.Program], data prependData) {
	if len(data.commitsToBeam) == 0 {
		return
	}
	// cherry-pick the commits into the new branch
	for _, commitToBeam := range data.commitsToBeam {
		prog.Value.Add(
			&opcodes.CherryPick{SHA: commitToBeam.SHA},
		)
	}
	// manually delete the beamed commits from the old branch
	prog.Value.Add(
		&opcodes.Checkout{Branch: data.initialBranch},
	)
	for _, commitToBeam := range data.commitsToBeam {
		prog.Value.Add(
			&opcodes.CommitRemove{SHA: commitToBeam.SHA},
		)
	}
	// sync the initial branch with the new parent branch to remove the moved commits from the initial branch
	initialBranchType := data.config.BranchType(data.initialBranch)
	syncWithParent(prog, data.targetBranch, data.initialBranchInfo, initialBranchType, data.config.NormalConfig)
	// go back to the target branch
	prog.Value.Add(
		&opcodes.Checkout{Branch: data.targetBranch},
	)
}

// basic sync of the current branch with its parent after beaming some commits into the parent
func syncWithParent(prog Mutable[program.Program], parentName gitdomain.LocalBranchName, initialBranchInfo gitdomain.BranchInfo, initialBranchType configdomain.BranchType, config config.NormalConfig) {
	if syncStrategy, hasSyncStrategy := afterBeamToParentSyncStrategy(initialBranchType, config).Get(); hasSyncStrategy {
		switch syncStrategy {
		case configdomain.SyncStrategyCompress, configdomain.SyncStrategyMerge:
			prog.Value.Add(
				&opcodes.MergeIntoCurrentBranch{BranchToMerge: parentName.BranchName()},
			)
			if initialBranchInfo.HasTrackingBranch() {
				prog.Value.Add(
					&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
				)
			}
		case configdomain.SyncStrategyRebase:
			prog.Value.Add(
				&opcodes.RebaseBranch{Branch: parentName.BranchName()},
			)
			if initialBranchInfo.HasTrackingBranch() {
				prog.Value.Add(
					&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
				)
			}
		case configdomain.SyncStrategyFFOnly:
			// the ff-only sync strategy does not sync with the parent
		}
	}
}
