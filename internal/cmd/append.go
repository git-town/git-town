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
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
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
	appendDesc = "Create a new feature branch as a child of the current branch"
	appendHelp = `
Consider this stack:

main
 \
* feature-1

We are on the "feature-1" branch,
which is a child of branch "main".
After running "git town append feature-2",
the repository will have these branches:

main
 \
  feature-1
   \
*   feature-2

The new branch "feature-2"
is a child of "feature-1".

If there are no uncommitted changes,
it also syncs all affected branches.
`
)

func appendCmd() *cobra.Command {
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addBeamFlag, readBeamFlag := flags.Beam()
	addCommitFlag, readCommitFlag := flags.Commit()
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addPushFlag, readPushFlag := flags.Push()
	addStashFlag, readStashFlag := flags.Stash()
	addSyncFlag, readSyncFlag := flags.Sync()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			beam, errBeam := readBeamFlag(cmd)
			commit, errCommit := readCommitFlag(cmd)
			commitMessage, errCommitMessage := readCommitMessageFlag(cmd)
			detached, errDetached := readDetachedFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			propose, errPropose := readProposeFlag(cmd)
			prototype, errPrototype := readPrototypeFlag(cmd)
			push, errPush := readPushFlag(cmd)
			stash, errStash := readStashFlag(cmd)
			sync, errSync := readSyncFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAutoResolve, errBeam, errCommit, errCommitMessage, errDetached, errDryRun, errPropose, errPrototype, errPush, errStash, errSync, errVerbose); err != nil {
				return err
			}
			if commitMessage.IsSome() || propose.ShouldPropose() {
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
			return executeAppend(executeAppendArgs{
				arg:           args[0],
				beam:          beam,
				cliConfig:     cliConfig,
				commit:        commit,
				commitMessage: commitMessage,
				propose:       propose,
				prototype:     prototype,
			})
		},
	}
	addAutoResolveFlag(&cmd)
	addBeamFlag(&cmd)
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addPushFlag(&cmd)
	addStashFlag(&cmd)
	addSyncFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeAppendArgs struct {
	arg           string
	beam          configdomain.Beam
	cliConfig     configdomain.PartialConfig
	commit        configdomain.Commit
	commitMessage Option[gitdomain.CommitMessage]
	propose       configdomain.Propose
	prototype     configdomain.Prototype
}

func executeAppend(args executeAppendArgs) error {
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
	data, flow, err := determineAppendData(determineAppendDataArgs{
		beam:          args.beam,
		commit:        args.commit,
		commitMessage: args.commitMessage,
		propose:       args.propose,
		prototype:     args.prototype,
		targetBranch:  gitdomain.NewLocalBranchName(args.arg),
	}, repo)
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
	runProgram := appendProgram(repo.Backend, data, repo.FinalMessages, false)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "append",
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

type appendFeatureData struct {
	beam                      configdomain.Beam
	branchInfos               gitdomain.BranchInfos
	branchInfosLastRun        Option[gitdomain.BranchInfos]
	branchesSnapshot          gitdomain.BranchesSnapshot
	branchesToSync            configdomain.BranchesToSync
	commit                    configdomain.Commit
	commitMessage             Option[gitdomain.CommitMessage]
	commitsToBeam             gitdomain.Commits
	config                    config.ValidatedConfig
	connector                 Option[forgedomain.Connector]
	hasOpenChanges            bool
	initialBranch             gitdomain.LocalBranchName
	initialBranchInfo         *gitdomain.BranchInfo
	inputs                    dialogcomponents.Inputs
	newBranchParentCandidates gitdomain.LocalBranchNames
	nonExistingBranches       gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos       gitdomain.BranchInfos
	previousBranch            Option[gitdomain.LocalBranchName]
	propose                   configdomain.Propose
	prototype                 configdomain.Prototype
	remotes                   gitdomain.Remotes
	stashSize                 gitdomain.StashSize
	targetBranch              gitdomain.LocalBranchName
}

func determineAppendData(args determineAppendDataArgs, repo execute.OpenRepoResult) (data appendFeatureData, flow configdomain.ProgramFlow, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
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
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	shouldFetch := true
	if repoStatus.OpenChanges {
		shouldFetch = false
	}
	if args.beam.ShouldBeam() || args.commit.ShouldCommit() {
		shouldFetch = false
	}
	if !config.AutoSync.ShouldSync() {
		shouldFetch = false
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 shouldFetch,
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
		return data, configdomain.ProgramFlowExit, errors.New(messages.AppendDetachedHead)
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	targetBranch := args.targetBranch
	if prefix, hasPrefix := config.BranchPrefix.Get(); hasPrefix {
		targetBranch = prefix.Apply(targetBranch)
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch, config.DevRemote)
	}
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
		LocalBranches:      branchesSnapshot.Branches.LocalBranches().NamesLocalBranches(),
		Remotes:            remotes,
		RepoStatus:         repoStatus,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, configdomain.ProgramFlowExit, err
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
	initialAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	slices.Reverse(initialAndAncestors)
	commitsToBeam := []gitdomain.Commit{}
	ancestor, hasAncestor := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage).Get()
	if args.beam.ShouldBeam() && hasAncestor {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor.BranchName())
		if err != nil {
			return data, configdomain.ProgramFlowExit, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, targetBranch, repo.Git, repo.Backend, inputs)
		if err != nil || exit {
			return data, configdomain.ProgramFlowExit, err
		}
	}
	if validatedConfig.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPropose {
		args.propose = true
	}
	return appendFeatureData{
		beam:                      args.beam,
		branchInfos:               branchesSnapshot.Branches,
		branchInfosLastRun:        branchInfosLastRun,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		commit:                    args.commit,
		commitMessage:             args.commitMessage,
		commitsToBeam:             commitsToBeam,
		config:                    validatedConfig,
		connector:                 connector,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		initialBranchInfo:         initialBranchInfo,
		inputs:                    inputs,
		newBranchParentCandidates: initialAndAncestors,
		nonExistingBranches:       nonExistingBranches,
		preFetchBranchInfos:       preFetchBranchSnapshot.Branches,
		previousBranch:            previousBranch,
		propose:                   args.propose,
		prototype:                 args.prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	}, configdomain.ProgramFlowContinue, nil
}

type determineAppendDataArgs struct {
	beam          configdomain.Beam
	commit        configdomain.Commit
	commitMessage Option[gitdomain.CommitMessage]
	propose       configdomain.Propose
	prototype     configdomain.Prototype
	targetBranch  gitdomain.LocalBranchName
}

func appendProgram(frontend subshelldomain.Runner, data appendFeatureData, finalMessages stringslice.Collector, beamCherryPick bool) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, finalMessages, frontend, data.config.NormalConfig.Order)
	if !data.hasOpenChanges && !data.beam.ShouldBeam() && !data.commit.ShouldCommit() && data.config.NormalConfig.AutoSync.ShouldSync() {
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
			Remotes:             data.remotes,
			PushBranches:        data.config.NormalConfig.PushBranches,
		})
	}
	prog.Value.Add(&opcodes.BranchCreateAndCheckoutExistingParent{
		Ancestors: data.newBranchParentCandidates,
		Branch:    data.targetBranch,
	})
	if data.remotes.HasRemote(data.config.NormalConfig.DevRemote) && data.config.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPush && data.config.NormalConfig.Offline.IsOnline() {
		prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.targetBranch})
	}
	prog.Value.Add(&opcodes.LineageParentSetFirstExisting{
		Branch:    data.targetBranch,
		Ancestors: data.newBranchParentCandidates,
	})
	var branchType configdomain.BranchType
	if data.prototype {
		branchType = configdomain.BranchTypePrototypeBranch
	} else if newBranchType, hasNewBranchType := data.config.NormalConfig.NewBranchType.Get(); hasNewBranchType {
		branchType = newBranchType.BranchType()
	} else {
		branchType = configdomain.BranchTypeFeatureBranch
	}
	prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: branchType})
	if data.commit {
		prog.Value.Add(
			&opcodes.Commit{
				AuthorOverride:                 None[gitdomain.Author](),
				FallbackToDefaultCommitMessage: false,
				Message:                        data.commitMessage,
			},
		)
	}
	moveCommitsToAppendedBranch(prog, data, beamCherryPick)
	if data.propose {
		title := None[gitdomain.ProposalTitle]()
		if commitMessage, has := data.commitMessage.Get(); has {
			title = Some(gitdomain.ProposalTitle(string(commitMessage)))
		}
		prog.Value.Add(
			&opcodes.BranchTrackingCreate{
				Branch: data.targetBranch,
			},
			&opcodes.ProposalCreate{
				Branch:        data.targetBranch,
				MainBranch:    data.config.ValidatedConfigData.MainBranch,
				ProposalBody:  None[gitdomain.ProposalBody](),
				ProposalTitle: title,
			},
		)
	}
	if data.commit {
		prog.Value.Add(
			&opcodes.Checkout{Branch: data.initialBranch},
		)
	} else {
		previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.initialBranch), data.previousBranch}
		cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
			DryRun:                   data.config.NormalConfig.DryRun,
			InitialStashSize:         data.stashSize,
			RunInGitRoot:             true,
			StashOpenChanges:         data.hasOpenChanges && data.config.NormalConfig.Stash.ShouldStash(),
			PreviousBranchCandidates: previousBranchCandidates,
		})
	}
	return optimizer.Optimize(prog.Immutable())
}

func moveCommitsToAppendedBranch(prog Mutable[program.Program], data appendFeatureData, performCherryPick bool) {
	if len(data.commitsToBeam) == 0 {
		return
	}
	if performCherryPick {
		for _, commitToBeam := range data.commitsToBeam {
			prog.Value.Add(
				&opcodes.CherryPick{SHA: commitToBeam.SHA},
			)
		}
	}
	prog.Value.Add(
		&opcodes.Checkout{
			Branch: data.initialBranch,
		},
	)
	for c := len(data.commitsToBeam) - 1; c >= 0; c-- {
		commitToBeam := data.commitsToBeam[c]
		prog.Value.Add(
			&opcodes.CommitRemove{
				SHA: commitToBeam.SHA,
			},
		)
	}
	if data.initialBranchInfo.HasTrackingBranch() {
		prog.Value.Add(
			&opcodes.PushCurrentBranchForceIgnoreError{},
		)
	}
	prog.Value.Add(
		&opcodes.Checkout{
			Branch: data.targetBranch,
		},
	)
	if !performCherryPick {
		prog.Value.Add(
			&opcodes.RebaseBranch{
				Branch: data.initialBranch.BranchName(),
			},
		)
	}
}
