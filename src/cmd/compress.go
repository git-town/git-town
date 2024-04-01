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
	"github.com/git-town/git-town/v13/src/git"
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineCompressConfig(repo, dryRun, verbose, message, stack)
	if err != nil || exit {
		return err
	}
	err = validateCompressConfig(config, repo.Runner)
	if err != nil {
		return err
	}
	program := compressProgram(config, initialBranchesSnapshot.Branches, repo.Runner)
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

type compressConfig struct {
	*configdomain.FullConfig
	branchesToCompress  gitdomain.LocalBranchNames
	commitMessage       gitdomain.CommitMessage
	compressEntireStack bool
	dialogTestInputs    components.TestInputs
	dryRun              bool
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	previousBranch      gitdomain.LocalBranchName
}

func determineCompressConfig(repo *execute.OpenRepoResult, dryRun, verbose bool, message gitdomain.CommitMessage, compressEntireStack bool) (*compressConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
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
	var branchesToCompress gitdomain.LocalBranchNames
	if compressEntireStack {
		branchesToCompress = repo.Runner.Config.FullConfig.Lineage.BranchLineage(initialBranch)
	} else {
		branchesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	return &compressConfig{
		FullConfig:          &repo.Runner.Config.FullConfig,
		branchesToCompress:  branchesToCompress,
		commitMessage:       message,
		compressEntireStack: compressEntireStack,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       branchesSnapshot.Active,
		previousBranch:      previousBranch,
	}, branchesSnapshot, stashSize, false, nil
}

func compressProgram(config *compressConfig, branchInfos gitdomain.BranchInfos, runner *git.ProdRunner) program.Program {
	prog := program.Program{}
	for _, branch := range config.Lineage.BranchLineage(config.initialBranch) {
		compressBranchProgram(&prog, branch, config, branchInfos, runner)
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})
	return prog
}

func compressBranchProgram(prog *program.Program, branch gitdomain.LocalBranchName, config *compressConfig, branchInfos gitdomain.BranchInfos, runner *git.ProdRunner) {
	branchType := config.BranchType(branch)
	if err := canCompressBranchType(branch, branchType); err != nil {
		return
	}
	parent := config.Lineage.Parent(branch)
	commits, err := runner.Backend.CommitsInBranch(branch, parent)
	if err != nil {
		panic(err.Error())
	}
	if len(commits) < 2 {
		return
	}
	branchCommitMessages := commits.Messages()
	commitMessage := slice.FirstNonEmpty(config.commitMessage, branchCommitMessages...)
	prog.Add(&opcodes.Checkout{Branch: branch})
	prog.Add(&opcodes.ResetCommitsInCurrentBranch{Parent: parent})
	prog.Add(&opcodes.CommitSquashedChanges{Message: commitMessage})
	branchInfo := branchInfos.FindByLocalName(branch)
	if branchInfo.HasRemoteBranch() && config.IsOnline() {
		prog.Add(&opcodes.ForcePushCurrentBranch{})
	}
}

func validateCompressConfig(config *compressConfig, run *git.ProdRunner) error {
	if !config.compressEntireStack {
		for _, branchToCompress := range config.branchesToCompress {
			branchType := config.BranchType(branchToCompress)
			if err := canCompressBranchType(branchToCompress, branchType); err != nil {
				return err
			}
			parentBranch := config.Lineage.Parent(branchToCompress)
			commits, err := run.Backend.CommitsInBranch(branchToCompress, parentBranch)
			if err != nil {
				return err
			}
			commitMessages := commits.Messages()
			switch len(commitMessages) {
			case 0:
				return errors.New(messages.CompressNoCommits)
			case 1:
				return errors.New(messages.CompressAlreadyOneCommit)
			}
		}
	}
	return nil
}

func canCompressBranchType(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
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
