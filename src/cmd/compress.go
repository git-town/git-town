package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/cli/flags"
	"github.com/git-town/git-town/v13/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/execute"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/gohacks/slice"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/git-town/git-town/v13/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v13/src/vm/interpreter/full"
	"github.com/git-town/git-town/v13/src/vm/opcodes"
	"github.com/git-town/git-town/v13/src/vm/program"
	"github.com/git-town/git-town/v13/src/vm/runstate"
	"github.com/spf13/cobra"
)

const compressDesc = "Squashes all commits on a feature branch down to a single commit"

const compressHelp = `
Compress is a more convenient way of running "git rebase --interactive" and choosing to squash or fixup all commits.
Branches must be synced before you compress them.
`

func compressCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addMessageFlag, readMessageFlag := flags.CommitMessage("customize the commit message")
	cmd := cobra.Command{
		Use:   "compress",
		Args:  cobra.NoArgs,
		Short: compressDesc,
		Long:  cmdhelpers.Long(compressDesc, compressHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeCompress(readDryRunFlag(cmd), readVerboseFlag(cmd), readMessageFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addMessageFlag(&cmd)
	return &cmd
}

func executeCompress(dryRun, verbose bool, message gitdomain.CommitMessage) error {
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineCompressConfig(repo, dryRun, verbose, message)
	if err != nil || exit {
		return err
	}
	err = validateCompressConfig(config)
	if err != nil {
		return err
	}
	steps := compressProgram(config)
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "compress",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            steps,
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

type compressConfig struct {
	*configdomain.FullConfig
	branchType             configdomain.BranchType
	dialogTestInputs       components.TestInputs
	dryRun                 bool
	existingCommitMessages gitdomain.CommitMessages // commit messages of the existing commits
	hasOpenChanges         bool
	initialBranch          gitdomain.BranchInfo
	newCommitMessage       gitdomain.CommitMessage // commit message to use for the new commit
	parentBranch           gitdomain.LocalBranchName
	previousBranch         gitdomain.LocalBranchName
}

func determineCompressConfig(repo *execute.OpenRepoResult, dryRun, verbose bool, message gitdomain.CommitMessage) (*compressConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	branchesSnapshot, stashSize, repoStatus, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FullConfig:            &repo.Runner.Config.FullConfig,
		HandleUnfinishedState: true,
		Repo:                  repo,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	initialBranch := branchesSnapshot.Active
	parentBranch := repo.Runner.Config.FullConfig.Lineage.Parent(initialBranch)
	initialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch)
	commits, err := repo.Runner.Backend.CommitsInBranch(initialBranch, parentBranch)
	if err != nil {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	commitMessages := commits.Messages()
	newCommitMessage := slice.FirstNonEmpty(message, commitMessages...)
	branchType := repo.Runner.Config.FullConfig.BranchType(initialBranch)
	return &compressConfig{
		FullConfig:             &repo.Runner.Config.FullConfig,
		branchType:             branchType,
		dialogTestInputs:       dialogTestInputs,
		dryRun:                 dryRun,
		existingCommitMessages: commitMessages,
		hasOpenChanges:         repoStatus.OpenChanges,
		initialBranch:          *initialBranchInfo,
		newCommitMessage:       newCommitMessage,
		parentBranch:           parentBranch,
		previousBranch:         previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func compressProgram(config *compressConfig) program.Program {
	prog := program.Program{}
	prog.Add(&opcodes.ResetCommitsInCurrentBranch{Parent: config.parentBranch})
	prog.Add(&opcodes.CommitSquashedChanges{Message: config.newCommitMessage})
	if config.initialBranch.HasRemoteBranch() && config.IsOnline() {
		prog.Add(&opcodes.ForcePushCurrentBranch{})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})
	return prog
}

func validateCompressConfig(config *compressConfig) error {
	if config.initialBranch.SyncStatus != gitdomain.SyncStatusUpToDate && config.initialBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return fmt.Errorf(messages.CompressUnsynced, config.initialBranch.LocalName)
	}
	if config.FullConfig.IsMainOrPerennialBranch(config.initialBranch.LocalName) {
	}
	switch len(config.existingCommitMessages) {
	case 0:
		return errors.New(messages.CompressNoCommits)
	case 1:
		return errors.New(messages.CompressAlreadyOneCommit)
	}
	switch config.branchType {
	case configdomain.BranchTypeParkedBranch, configdomain.BranchTypeFeatureBranch:
		// ok
	case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
		return errors.New(messages.CompressIsPerennial)
	case configdomain.BranchTypeObservedBranch:
		return fmt.Errorf(messages.CompressObservedBranch, config.initialBranch.LocalName)
	case configdomain.BranchTypeContributionBranch:
		return fmt.Errorf(messages.CompressContributionBranch, config.initialBranch.LocalName)
	}
	return nil
}
