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
		Use:   "compress",
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineCompressBranchesConfig(repo, dryRun, verbose, message, stack)
	if err != nil || exit {
		return err
	}
	err = validateCompressBranchesConfig(config)
	if err != nil {
		return err
	}
	program := compressProgram(config)
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "compress",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            program,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               nil,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

type compressBranchesConfig struct {
	*configdomain.FullConfig
	branchesToCompress  []compressBranchConfig
	compressEntireStack bool
	dialogTestInputs    components.TestInputs
	dryRun              bool
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	previousBranch      gitdomain.LocalBranchName
}

type compressBranchConfig struct {
	branchInfo       gitdomain.BranchInfo
	branchType       configdomain.BranchType
	commitCount      int                     // number of commits in this branch
	newCommitMessage gitdomain.CommitMessage // the commit message to use for the compressed commit in this branch
	parentBranch     gitdomain.LocalBranchName
}

func determineCompressBranchesConfig(repo *execute.OpenRepoResult, dryRun, verbose bool, message gitdomain.CommitMessage, compressEntireStack bool) (*compressBranchesConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FullConfig:            &repo.Runner.Config.FullConfig,
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
		branchNamesToCompress = repo.Runner.Config.FullConfig.Lineage.BranchLineageWithoutRoot(initialBranch)
	} else {
		branchNamesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesToCompress := make([]compressBranchConfig, len(branchNamesToCompress))
	for b, branchNameToCompress := range branchNamesToCompress {
		parentBranchPtr := repo.Runner.Config.FullConfig.Lineage.Parent(branchNameToCompress)
		if parentBranchPtr == nil {
			continue
		}
		parentBranch := *parentBranchPtr
		branchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToCompress)
		branchType := repo.Runner.Config.FullConfig.BranchType(branchNameToCompress.BranchName().LocalName())
		commits, err := repo.Runner.Backend.CommitsInBranch(branchNameToCompress.BranchName().LocalName(), parentBranch)
		if err != nil {
			return nil, branchesSnapshot, stashSize, exit, err
		}
		commitMessages := commits.Messages()
		newCommitMessage := slice.FirstNonEmpty(message, commitMessages...)
		commitCount := len(commitMessages)
		branchesToCompress[b] = compressBranchConfig{
			branchInfo:       *branchInfo,
			branchType:       branchType,
			commitCount:      commitCount,
			newCommitMessage: newCommitMessage,
			parentBranch:     parentBranch,
		}
	}
	return &compressBranchesConfig{
		FullConfig:          &repo.Runner.Config.FullConfig,
		branchesToCompress:  branchesToCompress,
		compressEntireStack: compressEntireStack,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		previousBranch:      previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func compressProgram(config *compressBranchesConfig) program.Program {
	prog := program.Program{}
	for _, branchToCompress := range config.branchesToCompress {
		compressBranchProgram(&prog, branchToCompress, config.Online(), config.initialBranch)
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

func compressBranchProgram(prog *program.Program, branch compressBranchConfig, online configdomain.Online, initialBranch gitdomain.LocalBranchName) {
	if !shouldCompressBranch(branch.branchInfo.LocalName, branch.branchType, initialBranch) {
		return
	}
	prog.Add(&opcodes.Checkout{Branch: branch.branchInfo.LocalName})
	prog.Add(&opcodes.ResetCommitsInCurrentBranch{Parent: branch.parentBranch})
	prog.Add(&opcodes.CommitSquashedChanges{Message: branch.newCommitMessage})
	if branch.branchInfo.HasRemoteBranch() && online.Bool() {
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

func validateCompressBranchesConfig(config *compressBranchesConfig) error {
	ec := execute.FailureCollector{}
	for _, compressBranchConfig := range config.branchesToCompress {
		ec.Check(validateBranchIsSynced(compressBranchConfig.branchInfo.LocalName, compressBranchConfig.branchInfo.SyncStatus))
		ec.Check(validateCanCompressBranchType(compressBranchConfig.branchInfo.LocalName, compressBranchConfig.branchType))
		ec.Check(validateBranchHasMultipleCommits(compressBranchConfig.branchInfo.LocalName, compressBranchConfig.commitCount))
	}
	return ec.Err
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
	panic("unhandled syncstatus" + syncStatus.String())
}
