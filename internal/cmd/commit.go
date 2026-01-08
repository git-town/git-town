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
	"github.com/spf13/cobra"
)

const (
	commitDesc = "Commit into another branch"
	commitHelp = `
Allows you to commit the currently staged changes
into another branch without needing to change branches.`
)

func commitCmd() *cobra.Command {
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "commit",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   commitDesc,
		Long:    cmdhelpers.Long(commitDesc, commitHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			message, errMessage := readMessageFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errMessage, errVerbose); err != nil {
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
			return executeCommit(cliConfig, message)
		},
	}
	addMessageFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeCommit(cliConfig configdomain.PartialConfig, message Option[gitdomain.CommitMessage]) error {
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
	data, flow, err := determineCommitData(repo)
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
	err = validateCommitData(data)
	if err != nil {
		return err
	}
	runProgram, finalUndoProgram := commitProgram(repo, data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "commit",
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		FinalUndoProgram:      finalUndoProgram,
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

type commitData struct {
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchToCommitInto gitdomain.LocalBranchName
	commitMessage      Option[gitdomain.CommitMessage]
	config             config.ValidatedConfig
	connector          Option[forgedomain.Connector]
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	inputs             dialogcomponents.Inputs
	stashSize          gitdomain.StashSize
}

func determineCommitData(args []string, repo execute.OpenRepoResult) (data commitData, flow configdomain.ProgramFlow, err error) {
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
		GiteaToken:           config.GiteaToken,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
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
		return data, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return data, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return data, configdomain.ProgramFlowExit, errors.New(messages.DeleteRepoHasDetachedHead)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{},
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
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	return commitData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		config:             validatedConfig,
		connector:          connector,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		inputs:             inputs,
		stashSize:          stashSize,
	}, configdomain.ProgramFlowContinue, nil
}

func commitProgram(repo execute.OpenRepoResult, data commitData, finalMessages stringslice.Collector) (runProgram, finalUndoProgram program.Program) {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Backend, data.config.NormalConfig.Order)
	undoProg := NewMutable(&program.Program{})
	switch data.branchToDeleteType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		deleteFeatureBranch(prog, undoProg, data)
	case
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeContributionBranch:
		deleteLocalBranch(prog, undoProg, data)
	case
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		panic(fmt.Sprintf("this branch type should have been filtered in validation: %s", data.branchToDeleteType))
	}
	localBranchNameToDelete := data.branchToDeleteInfo.LocalBranchName()
	if _, hasOverride := data.config.NormalConfig.BranchTypeOverrides[localBranchNameToDelete]; hasOverride {
		prog.Value.Add(&opcodes.BranchTypeOverrideRemove{
			Branch: localBranchNameToDelete,
		})
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch, Some(data.initialBranch)},
	})
	return optimizer.Optimize(prog.Immutable()), undoProg.Immutable()
}

func deleteFeatureBranch(prog, finalUndoProgram Mutable[program.Program], data deleteData) {
	trackingBranchToDelete, hasTrackingBranchToDelete := data.branchToDeleteInfo.RemoteName.Get()
	if data.branchToDeleteInfo.SyncStatus != gitdomain.SyncStatusDeletedAtRemote && hasTrackingBranchToDelete && data.config.NormalConfig.Offline.IsOnline() {
		ship.UpdateChildBranchProposalsToGrandParent(prog.Value, data.proposalsOfChildBranches)
		prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: trackingBranchToDelete})
	}
	deleteLocalBranch(prog, finalUndoProgram, data)
	if data.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		sync.AddSyncProposalsProgram(sync.AddSyncProposalsProgramArgs{
			ChangedBranches: data.oldClan.Remove(data.branchToDeleteInfo.GetLocalOrRemoteNameAsLocalName()),
			Config:          data.config,
			Program:         prog,
		})
	}
}

func deleteLocalBranch(prog, finalUndoProgram Mutable[program.Program], data deleteData) {
	if localBranchToDelete, hasLocalBranchToDelete := data.branchToDeleteInfo.LocalName.Get(); hasLocalBranchToDelete {
		if data.initialBranch == localBranchToDelete {
			if data.hasOpenChanges {
				prog.Value.Add(&opcodes.ChangesStage{})
				prog.Value.Add(&opcodes.CommitWithMessage{
					AuthorOverride: None[gitdomain.Author](),
					CommitHook:     configdomain.CommitHookEnabled,
					Message:        "Committing open changes on deleted branch",
				})
				// update the registered initial SHA for this branch so that undo restores the just committed changes
				prog.Value.Add(&opcodes.SnapshotInitialUpdateLocalSHAIfNeeded{Branch: data.initialBranch})
				// when undoing, manually undo the just committed changes so that they are uncommitted again
				finalUndoProgram.Value.Add(&opcodes.CheckoutIfNeeded{Branch: localBranchToDelete})
				finalUndoProgram.Value.Add(&opcodes.UndoLastCommit{})
			}
		}
		// delete the commits of this branch from all descendents
		if data.config.NormalConfig.SyncFeatureStrategy == configdomain.SyncFeatureStrategyRebase {
			descendents := data.config.NormalConfig.Lineage.Descendants(localBranchToDelete, data.config.NormalConfig.Order)
			for _, descendent := range descendents {
				if branchInfo, hasBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get(); hasBranchInfo {
					parent := data.config.NormalConfig.Lineage.Parent(descendent).GetOr(data.config.ValidatedConfigData.MainBranch)
					if parent == localBranchToDelete {
						parent = data.config.NormalConfig.Lineage.Parent(parent).GetOr(data.config.ValidatedConfigData.MainBranch)
					}
					sync.RemoveAncestorCommits(sync.RemoveAncestorCommitsArgs{
						Ancestor:          localBranchToDelete.BranchName(),
						Branch:            descendent,
						HasTrackingBranch: branchInfo.HasTrackingBranch(),
						Program:           prog,
						RebaseOnto:        parent,
					})
				}
			}
		}
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.branchWhenDone})
		prog.Value.Add(&opcodes.BranchLocalDelete{
			Branch: localBranchToDelete,
		})
		if !data.config.NormalConfig.DryRun {
			sync.RemoveBranchConfiguration(sync.RemoveBranchConfigurationArgs{
				Branch:  localBranchToDelete,
				Lineage: data.config.NormalConfig.Lineage,
				Order:   data.config.NormalConfig.Order,
				Program: prog,
			})
		}
	}
}

func determineBranchWhenDone(args branchWhenDoneArgs) gitdomain.LocalBranchName {
	if args.branchToDelete != args.initialBranch {
		return args.initialBranch
	}
	// here we are deleting the initial branch
	previousBranch, hasPreviousBranch := args.previousBranch.Get()
	if !hasPreviousBranch || previousBranch == args.initialBranch {
		return args.mainBranch
	}
	// here we could return the previous branch
	if previousBranchInfo, hasPreviousBranchInfo := args.branches.FindByLocalName(previousBranch).Get(); hasPreviousBranchInfo {
		if previousBranchInfo.SyncStatus != gitdomain.SyncStatusOtherWorktree {
			return previousBranch
		}
	}
	// here the previous branch is checked out in another worktree --> cannot return it
	return args.mainBranch
}

type branchWhenDoneArgs struct {
	branchToDelete gitdomain.LocalBranchName
	branches       gitdomain.BranchInfos
	initialBranch  gitdomain.LocalBranchName
	mainBranch     gitdomain.LocalBranchName
	previousBranch Option[gitdomain.LocalBranchName]
}

func validateDeleteData(data deleteData) error {
	switch data.branchToDeleteType {
	case
		configdomain.BranchTypeContributionBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeObservedBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.DeleteCannotDeleteMainBranch)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.DeleteCannotDeletePerennialBranches)
	}
	return nil
}
