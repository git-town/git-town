package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/cmd/sync"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/optimizer"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/git-town/git-town/v21/pkg/set"
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
	addBeamFlag, readBeamFlag := flags.Beam()
	addCommitFlag, readCommitFlag := flags.Commit()
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addStashFlag, readStashFlag := flags.Stash()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "append <branch>",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.ExactArgs(1),
		Short:   appendDesc,
		Long:    cmdhelpers.Long(appendDesc, appendHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			beam, errBeam := readBeamFlag(cmd)
			commit, errCommit := readCommitFlag(cmd)
			commitMessage, errCommitMessage := readCommitMessageFlag(cmd)
			detached, errDetached := readDetachedFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			propose, errPropose := readProposeFlag(cmd)
			prototype, errPrototype := readPrototypeFlag(cmd)
			stash, errStash := readStashFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errBeam, errCommit, errCommitMessage, errDetached, errDryRun, errAutoResolve, errPropose, errPrototype, errStash, errVerbose); err != nil {
				return err
			}
			if commitMessage.IsSome() || propose.IsTrue() {
				commit = true
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:  autoResolve,
				Detached:     detached,
				DryRun:       dryRun,
				PushBranches: None[configdomain.PushBranches](), // TODO: add CLI flag and use its value here
				Stash:        stash,
				Verbose:      verbose,
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
	addBeamFlag(&cmd)
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addAutoResolveFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addStashFlag(&cmd)
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
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineAppendData(determineAppendDataArgs{
		beam:          args.beam,
		commit:        args.commit,
		commitMessage: args.commitMessage,
		propose:       args.propose,
		prototype:     args.prototype,
		targetBranch:  gitdomain.NewLocalBranchName(args.arg),
	}, repo)
	if err != nil || exit {
		return err
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

func determineAppendData(args determineAppendDataArgs, repo execute.OpenRepoResult) (data appendFeatureData, exit dialogdomain.Exit, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
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
		return data, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 !repoStatus.OpenChanges && args.beam.IsFalse() && args.commit.IsFalse(),
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
	if err != nil || exit {
		return data, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	if branchesSnapshot.Branches.HasLocalBranch(args.targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, args.targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(args.targetBranch, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, args.targetBranch)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
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
		LocalBranches:      branchesSnapshot.Branches.LocalBranches().Names(),
		Remotes:            remotes,
		RepoStatus:         repoStatus,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	branchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	if repo.UnvalidatedConfig.NormalConfig.Detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	initialAndAncestors := validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
	slices.Reverse(initialAndAncestors)
	commitsToBeam := []gitdomain.Commit{}
	ancestor, hasAncestor := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage).Get()
	if args.beam.IsTrue() && hasAncestor {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor.BranchName())
		if err != nil {
			return data, false, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, args.targetBranch, repo.Git, repo.Backend, inputs)
		if err != nil || exit {
			return data, exit, err
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
		targetBranch:              args.targetBranch,
	}, false, nil
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
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, finalMessages, frontend)
	if !data.hasOpenChanges && data.beam.IsFalse() && data.commit.IsFalse() {
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
			PushBranches:        true,
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
	if data.prototype {
		prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: configdomain.BranchTypePrototypeBranch})
	} else {
		if newBranchType, hasNewBranchType := data.config.NormalConfig.NewBranchType.Get(); hasNewBranchType {
			switch newBranchType.BranchType() {
			case
				configdomain.BranchTypeContributionBranch,
				configdomain.BranchTypeObservedBranch,
				configdomain.BranchTypeParkedBranch,
				configdomain.BranchTypePerennialBranch,
				configdomain.BranchTypePrototypeBranch,
				configdomain.BranchTypeFeatureBranch:
				prog.Value.Add(&opcodes.BranchTypeOverrideSet{Branch: data.targetBranch, BranchType: newBranchType.BranchType()})
			case configdomain.BranchTypeMainBranch:
			}
		}
	}
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
			StashOpenChanges:         data.hasOpenChanges && data.config.NormalConfig.Stash.IsTrue(),
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
