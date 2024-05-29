package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
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

func executeCompress(dryRun, verbose bool, message Option[gitdomain.CommitMessage], stack bool) error {
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
	data, exit, err := determineCompressBranchesData(repo, dryRun, verbose, message, stack)
	if err != nil || exit {
		return err
	}
	program := compressProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               compressCommand,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            program,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               nil,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
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
	branchesSnapshot    gitdomain.BranchesSnapshot
	branchesToCompress  []compressBranchData
	compressEntireStack bool
	config              config.ValidatedConfig
	dialogTestInputs    components.TestInputs
	dryRun              bool
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	previousBranch      gitdomain.LocalBranchName
	stashSize           gitdomain.StashSize
}

type compressBranchData struct {
	branchType       configdomain.BranchType
	commitCount      int // number of commits in this branch
	hasTracking      bool
	name             gitdomain.LocalBranchName
	newCommitMessage gitdomain.CommitMessage // the commit message to use for the compressed commit in this branch
	parentBranch     gitdomain.LocalBranchName
}

func determineCompressBranchesData(repo execute.OpenRepoResult, dryRun, verbose bool, message Option[gitdomain.CommitMessage], compressEntireStack bool) (*compressBranchesData, bool, error) {
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Backend.RepoStatus()
	if err != nil {
		return nil, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		RootDir:               repo.RootDir,
		UnvalidatedConfig:     repo.UnvalidatedConfig,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return nil, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return nil, exit, err
	}
	var branchNamesToCompress gitdomain.LocalBranchNames
	if compressEntireStack {
		branchNamesToCompress = validatedConfig.Config.Lineage.BranchLineageWithoutRoot(initialBranch)
	} else {
		branchNamesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesToCompress := make([]compressBranchData, 0, len(branchNamesToCompress))
	for _, branchNameToCompress := range branchNamesToCompress {
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToCompress).Get()
		if !hasBranchInfo {
			return nil, exit, fmt.Errorf(messages.CompressNoBranchInfo, branchNameToCompress)
		}
		branchType := validatedConfig.Config.BranchType(branchNameToCompress.BranchName().LocalName())
		if err := validateCanCompressBranchType(branchNameToCompress, branchType); err != nil {
			return nil, exit, err
		}
		if err := validateBranchIsSynced(branchNameToCompress, branchInfo.SyncStatus); err != nil {
			return nil, exit, err
		}
		parent := validatedConfig.Config.Lineage.Parent(branchNameToCompress)
		commits, err := repo.Backend.CommitsInBranch(branchNameToCompress.BranchName().LocalName(), parent)
		if err != nil {
			return nil, exit, err
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
			return nil, exit, fmt.Errorf(messages.CompressBranchNoParent, branchNameToCompress)
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
		return nil, exit, fmt.Errorf(messages.CompressNoCommits, branchNamesToCompress[0])
	}
	if len(branchesToCompress) == 1 && branchesToCompress[0].commitCount == 1 {
		return nil, exit, fmt.Errorf(messages.CompressAlreadyOneCommit, branchNamesToCompress[0])
	}
	return &compressBranchesData{
		branchesSnapshot:    branchesSnapshot,
		branchesToCompress:  branchesToCompress,
		compressEntireStack: compressEntireStack,
		config:              validatedConfig,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		previousBranch:      previousBranch,
		stashSize:           stashSize,
	}, false, nil
}

func compressProgram(data *compressBranchesData) program.Program {
	prog := program.Program{}
	for _, branchToCompress := range data.branchesToCompress {
		compressBranchProgram(&prog, branchToCompress, data.config.Config.Online(), data.initialBranch)
	}
	prog.Add(&opcodes.Checkout{Branch: data.initialBranch.BranchName().LocalName()})
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{data.previousBranch},
	})
	return prog
}

func compressBranchProgram(prog *program.Program, data compressBranchData, online configdomain.Online, initialBranch gitdomain.LocalBranchName) {
	if !shouldCompressBranch(data.name, data.branchType, initialBranch) {
		return
	}
	prog.Add(&opcodes.Checkout{Branch: data.name})
	prog.Add(&opcodes.ResetCommitsInCurrentBranch{Parent: data.parentBranch})
	prog.Add(&opcodes.CommitSquashedChanges{Message: Some(data.newCommitMessage)})
	if data.hasTracking && online.Bool() {
		prog.Add(&opcodes.ForcePushCurrentBranch{})
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

func validateBranchIsSynced(branchName gitdomain.LocalBranchName, syncStatus gitdomain.SyncStatus) error {
	switch syncStatus {
	case gitdomain.SyncStatusUpToDate, gitdomain.SyncStatusLocalOnly:
		return nil
	case gitdomain.SyncStatusNotInSync, gitdomain.SyncStatusDeletedAtRemote, gitdomain.SyncStatusRemoteOnly, gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.CompressUnsynced, branchName)
	}
	panic("unhandled syncstatus: " + syncStatus.String())
}
