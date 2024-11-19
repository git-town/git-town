package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const compressCommand = "compress"

const compressDesc = "Squash all commits on feature branches down to a single commit"

const compressHelp = `
Compress is a more convenient way of running "git rebase --interactive"
and choosing to fixup all commits.
Branches must be synced before you compress them.

By default, this command compresses only the branch you are on.
With the --stack switch it compresses all branches in the current stack.

The compressed commit uses the commit message of the first commit in the branch.
You can provide a custom commit message with the -m switch.
`

func compressCmd() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addMessageFlag, readMessageFlag := flags.CommitMessage("customize the commit message")
	addStackFlag, readStackFlag := flags.Stack("Compress the entire stack")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   compressCommand,
		Args:  cobra.NoArgs,
		Short: compressDesc,
		Long:  cmdhelpers.Long(compressDesc, compressHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			message, err := readMessageFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			stack, err := readStackFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeCompress(dryRun, verbose, message, stack)
		},
	}
	addDryRunFlag(&cmd)
	addMessageFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeCompress(dryRun configdomain.DryRun, verbose configdomain.Verbose, message Option[gitdomain.CommitMessage], compressEntireStack configdomain.FullStack) error {
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
	data, exit, err := determineCompressBranchesData(repo, dryRun, verbose, message, compressEntireStack)
	if err != nil || exit {
		return err
	}
	runProgram := compressProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               compressCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               None[hostingdomain.Connector](),
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type compressBranchesData struct {
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchesToCompress []compressBranchData
	config             config.ValidatedConfig
	dialogTestInputs   components.TestInputs
	dryRun             configdomain.DryRun
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	previousBranch     Option[gitdomain.LocalBranchName]
	stashSize          gitdomain.StashSize
}

type compressBranchData struct {
	branchType       configdomain.BranchType
	commitCount      int // number of commits in this branch
	hasTracking      bool
	name             gitdomain.LocalBranchName
	newCommitMessage gitdomain.CommitMessage // the commit message to use for the compressed commit in this branch
	parentBranch     gitdomain.LocalBranchName
}

func determineCompressBranchesData(repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose, message Option[gitdomain.CommitMessage], compressEntireStack configdomain.FullStack) (data compressBranchesData, exit bool, err error) {
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
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
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
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
	var branchNamesToCompress gitdomain.LocalBranchNames
	if compressEntireStack {
		branchNamesToCompress = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch)
	} else {
		branchNamesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesToCompress := make([]compressBranchData, 0, len(branchNamesToCompress))
	for _, branchNameToCompress := range branchNamesToCompress {
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToCompress).Get()
		if !hasBranchInfo {
			return data, exit, fmt.Errorf(messages.CompressNoBranchInfo, branchNameToCompress)
		}
		branchType := validatedConfig.BranchType(branchNameToCompress.BranchName().LocalName())
		if err := validateCanCompressBranchType(branchNameToCompress, branchType); err != nil {
			return data, exit, err
		}
		if err := validateBranchIsSynced(branchNameToCompress, branchInfo.SyncStatus); err != nil {
			return data, exit, err
		}
		parent := validatedConfig.NormalConfig.Lineage.Parent(branchNameToCompress)
		commits, err := repo.Git.CommitsInBranch(repo.Backend, branchNameToCompress.BranchName().LocalName(), parent)
		if err != nil {
			return data, exit, err
		}
		commitCount := len(commits)
		if commitCount == 0 {
			continue
		}
		var newCommitMessage gitdomain.CommitMessage
		if messageContent, has := message.Get(); has {
			newCommitMessage = messageContent
		} else {
			newCommitMessage = commits.Messages()[0]
		}
		parentBranch, hasParent := parent.Get()
		if !hasParent {
			return data, exit, fmt.Errorf(messages.CompressBranchNoParent, branchNameToCompress)
		}
		hasRemoteBranch, _, _ := branchInfo.HasRemoteBranch()
		branchesToCompress = append(branchesToCompress, compressBranchData{
			branchType:       branchType,
			commitCount:      commitCount,
			hasTracking:      hasRemoteBranch,
			name:             branchNameToCompress,
			newCommitMessage: newCommitMessage,
			parentBranch:     parentBranch,
		})
	}
	if len(branchesToCompress) == 0 {
		return data, exit, fmt.Errorf(messages.CompressNoCommits, branchNamesToCompress[0])
	}
	if len(branchesToCompress) == 1 && branchesToCompress[0].commitCount == 1 {
		return data, exit, fmt.Errorf(messages.CompressAlreadyOneCommit, branchNamesToCompress[0])
	}
	return compressBranchesData{
		branchesSnapshot:   branchesSnapshot,
		branchesToCompress: branchesToCompress,
		config:             validatedConfig,
		dialogTestInputs:   dialogTestInputs,
		dryRun:             dryRun,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, false, nil
}

func compressProgram(data compressBranchesData) program.Program {
	prog := NewMutable(&program.Program{})
	for _, branchToCompress := range data.branchesToCompress {
		compressBranchProgram(prog, branchToCompress, data.config.NormalConfig.Online(), data.initialBranch)
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.initialBranch.BranchName().LocalName()})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return prog.Immutable()
}

func compressBranchProgram(prog Mutable[program.Program], data compressBranchData, online configdomain.Online, initialBranch gitdomain.LocalBranchName) {
	if !shouldCompressBranch(data.name, data.branchType, initialBranch) {
		return
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.name})
	prog.Value.Add(&opcodes.BranchCurrentReset{Base: data.parentBranch.BranchName()})
	prog.Value.Add(&opcodes.CommitWithMessage{
		AuthorOverride: None[gitdomain.Author](),
		Message:        data.newCommitMessage,
	})
	if data.hasTracking && online.IsTrue() {
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{ForceIfIncludes: true})
	}
}

func shouldCompressBranch(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType, initialBranchName gitdomain.LocalBranchName) bool {
	if branchName == initialBranchName {
		return true
	}
	return branchType != configdomain.BranchTypeParkedBranch
}

func validateCanCompressBranchType(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	switch branchType {
	case
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypePrototypeBranch:
		return nil
	case
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return errors.New(messages.CompressIsPerennial)
	case configdomain.BranchTypeObservedBranch:
		return fmt.Errorf(messages.CompressObservedBranch, branchName)
	case configdomain.BranchTypeContributionBranch:
		return fmt.Errorf(messages.CompressContributionBranch, branchName)
	}
	return nil
}

func validateBranchIsSynced(branchName gitdomain.LocalBranchName, syncStatus gitdomain.SyncStatus) error {
	switch syncStatus {
	case
		gitdomain.SyncStatusUpToDate,
		gitdomain.SyncStatusLocalOnly:
		return nil
	case
		gitdomain.SyncStatusNotInSync,
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind,
		gitdomain.SyncStatusDeletedAtRemote,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.CompressUnsynced, branchName)
	}
	panic("unhandled syncstatus: " + syncStatus.String())
}
