package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
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
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addMessageFlag, readMessageFlag := flags.CommitMessage("customize the commit message")
	addStackFlag, readStackFlag := flags.Bool("stack", "s", "Compress the entire stack", flags.FlagTypeNonPersistent)
	cmd := cobra.Command{
		Use:   compressCommand,
		Args:  cobra.NoArgs,
		Short: compressDesc,
		Long:  cmdhelpers.Long(compressDesc, compressHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeCompress(readDryRunFlag(cmd), readVerboseFlag(cmd), readMessageFlag(cmd), readStackFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	addStackFlag(&cmd)
	return &cmd
}

func executeCompress(dryRun, verbose bool, message gitdomain.CommitMessage, stack bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineCompressBranchesData(repo, dryRun, verbose, message, stack)
	if err != nil || exit {
		return err
	}
	program := compressProgram(config)
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               compressCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            program,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Config:                  config.config,
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type compressBranchesData struct {
	branchesToCompress  []compressBranchData
	compressEntireStack bool
	config              configdomain.FullConfig
	dialogTestInputs    components.TestInputs
	dryRun              bool
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	previousBranch      gitdomain.LocalBranchName
}

type compressBranchData struct {
	branchInfo       gitdomain.BranchInfo
	branchType       configdomain.BranchType
	commitCount      int                     // number of commits in this branch
	newCommitMessage gitdomain.CommitMessage // the commit message to use for the compressed commit in this branch
	parentBranch     gitdomain.LocalBranchName
}

func determineCompressBranchesData(repo *execute.OpenRepoResult, dryRun, verbose bool, message gitdomain.CommitMessage, compressEntireStack bool) (*compressBranchesData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	initialBranch := branchesSnapshot.Active.BranchName().LocalName()
	var branchNamesToCompress gitdomain.LocalBranchNames
	if compressEntireStack {
		branchNamesToCompress = repo.Runner.Config.Config.Lineage.BranchLineageWithoutRoot(initialBranch)
	} else {
		branchNamesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesToCompress := make([]compressBranchData, len(branchNamesToCompress))
	for b, branchNameToCompress := range branchNamesToCompress {
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToCompress).Get()
		if !hasBranchInfo {
			return nil, branchesSnapshot, stashSize, exit, fmt.Errorf(messages.CompressNoBranchInfo, branchNameToCompress)
		}
		branchType := repo.Runner.Config.Config.BranchType(branchNameToCompress.BranchName().LocalName())
		if err := validateCanCompressBranchType(branchInfo.LocalName, branchType); err != nil {
			return nil, branchesSnapshot, stashSize, exit, err
		}
		if err := validateBranchIsSynced(branchInfo.LocalName, branchInfo.SyncStatus); err != nil {
			return nil, branchesSnapshot, stashSize, exit, err
		}
		parent := repo.Runner.Config.Config.Lineage.Parent(branchNameToCompress)
		commits, err := repo.Runner.Backend.CommitsInBranch(branchNameToCompress.BranchName().LocalName(), parent)
		if err != nil {
			return nil, branchesSnapshot, stashSize, exit, err
		}
		commitMessages := commits.Messages()
		newCommitMessage := slice.FirstNonEmpty(message, commitMessages...)
		commitCount := len(commitMessages)
		if err := validateBranchHasMultipleCommits(branchInfo.LocalName, commitCount); err != nil {
			return nil, branchesSnapshot, stashSize, exit, err
		}
		parentBranch, hasParent := parent.Get()
		if !hasParent {
			return nil, branchesSnapshot, stashSize, exit, fmt.Errorf(messages.CompressBranchNoParent, branchNameToCompress)
		}
		branchesToCompress[b] = compressBranchData{
			branchInfo:       branchInfo,
			branchType:       branchType,
			commitCount:      commitCount,
			newCommitMessage: newCommitMessage,
			parentBranch:     parentBranch,
		}
	}
	return &compressBranchesData{
		branchesToCompress:  branchesToCompress,
		compressEntireStack: compressEntireStack,
		config:              repo.Runner.Config.Config,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		previousBranch:      previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func compressProgram(config *compressBranchesData) program.Program {
	prog := program.Program{}
	for _, branchToCompress := range config.branchesToCompress {
		compressBranchProgram(&prog, branchToCompress, config.config.Online(), config.initialBranch)
	}
	prog.Add(&opcodes.Checkout{Branch: config.initialBranch.BranchName().LocalName()})
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})
	return prog
}

func compressBranchProgram(prog *program.Program, data compressBranchData, online configdomain.Online, initialBranch gitdomain.LocalBranchName) {
	if !shouldCompressBranch(data.branchInfo.LocalName, data.branchType, initialBranch) {
		return
	}
	prog.Add(&opcodes.Checkout{Branch: data.branchInfo.LocalName})
	prog.Add(&opcodes.ResetCommitsInCurrentBranch{Parent: data.parentBranch})
	prog.Add(&opcodes.CommitSquashedChanges{Message: data.newCommitMessage})
	if data.branchInfo.HasRemoteBranch() && online.Bool() {
		prog.Add(&opcodes.ForcePushCurrentBranch{})
	}
}

func shouldCompressBranch(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType, initialBranchName gitdomain.LocalBranchName) bool {
	isInitialBranch := branchName == initialBranchName
	if branchType == configdomain.BranchTypeParkedBranch && !isInitialBranch {
		return false
	}
	return true
}

func validateCanCompressBranchType(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	switch branchType {
	case configdomain.BranchTypeParkedBranch, configdomain.BranchTypeFeatureBranch:
		return nil
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		return errors.New(messages.CompressIsPerennial)
	case configdomain.BranchTypeObservedBranch:
		return fmt.Errorf(messages.CompressObservedBranch, branchName)
	case configdomain.BranchTypeContributionBranch:
		return fmt.Errorf(messages.CompressContributionBranch, branchName)
	}
	return nil
}

func validateBranchHasMultipleCommits(branch gitdomain.LocalBranchName, commitCount int) error {
	switch commitCount {
	case 0:
		return fmt.Errorf(messages.CompressNoCommits, branch)
	case 1:
		return fmt.Errorf(messages.CompressAlreadyOneCommit, branch)
	}
	return nil
}

func validateBranchIsSynced(branchName gitdomain.LocalBranchName, syncStatus gitdomain.SyncStatus) error {
	switch syncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusLocalOnly:
		return nil
	case gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.CompressUnsynced, branchName)
	}
	panic("unhandled syncstatus: " + syncStatus.String())
}
