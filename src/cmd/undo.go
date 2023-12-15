package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cli/log"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/hosting"
	"github.com/git-town/git-town/v11/src/hosting/github"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/git-town/git-town/v11/src/vm/interpreter"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/runstate"
	"github.com/git-town/git-town/v11/src/vm/statefile"
	"github.com/spf13/cobra"
)

const undoDesc = "Undoes the most recent Git Town command"

func undoCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: "errors",
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    long(undoDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeUndo(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUndo(verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           false,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  true,
	})
	if err != nil {
		return err
	}
	config, initialStashSnaphot, lineage, err := determineUndoConfig(repo, verbose)
	if err != nil {
		return err
	}
	undoRunState, err := determineUndoRunState(config, repo)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	return interpreter.Execute(interpreter.ExecuteArgs{
		RunState:                &undoRunState,
		Run:                     repo.Runner,
		Connector:               config.connector,
		Verbose:                 verbose,
		Lineage:                 lineage,
		NoPushHook:              config.pushHook.Negate(),
		RootDir:                 repo.RootDir,
		InitialBranchesSnapshot: config.initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSnapshot:    initialStashSnaphot,
	})
}

type undoConfig struct {
	connector               hosting.Connector
	hasOpenChanges          bool
	initialBranchesSnapshot domain.BranchesSnapshot
	mainBranch              domain.LocalBranchName
	lineage                 configdomain.Lineage
	previousBranch          domain.LocalBranchName
	pushHook                configdomain.PushHook
}

func determineUndoConfig(repo *execute.OpenRepoResult, verbose bool) (*undoConfig, domain.StashSnapshot, configdomain.Lineage, error) {
	lineage := repo.Runner.GitTown.Lineage(repo.Runner.Backend.GitTown.RemoveLocalConfigValue)
	pushHook, err := repo.Runner.GitTown.PushHook()
	if err != nil {
		return nil, domain.EmptyStashSnapshot(), lineage, err
	}
	_, initialBranchesSnapshot, initialStashSnapshot, _, err := execute.LoadBranches(execute.LoadBranchesArgs{
		Repo:                  repo,
		Verbose:               verbose,
		Fetch:                 false,
		HandleUnfinishedState: false,
		Lineage:               lineage,
		PushHook:              pushHook,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
	})
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	mainBranch := repo.Runner.Backend.GitTown.MainBranch()
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	hostingService, err := repo.Runner.GitTown.HostingService()
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	originURL := repo.Runner.GitTown.OriginURL()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		HostingService:  hostingService,
		GetSHAForBranch: repo.Runner.Backend.SHAForBranch,
		OriginURL:       originURL,
		GiteaAPIToken:   repo.Runner.GitTown.GiteaToken(),
		GithubAPIToken:  github.GetAPIToken(repo.Runner.GitTown.GitHubToken()),
		GitlabAPIToken:  repo.Runner.GitTown.GitLabToken(),
		MainBranch:      mainBranch,
		Log:             log.Printing{},
	})
	if err != nil {
		return nil, initialStashSnapshot, lineage, err
	}
	return &undoConfig{
		connector:               connector,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: initialBranchesSnapshot,
		lineage:                 lineage,
		mainBranch:              mainBranch,
		previousBranch:          previousBranch,
		pushHook:                pushHook,
	}, initialStashSnapshot, lineage, nil
}

func determineUndoRunState(config *undoConfig, repo *execute.OpenRepoResult) (runstate.RunState, error) {
	runState, err := statefile.Load(repo.RootDir)
	if err != nil {
		return runstate.EmptyRunState(), fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	if runState == nil {
		return runstate.EmptyRunState(), fmt.Errorf(messages.UndoNothingToDo)
	}
	var undoRunState runstate.RunState
	if runState.IsUnfinished() {
		undoRunState = runState.CreateAbortRunState()
	} else {
		undoRunState = runState.CreateUndoRunState()
	}
	wrap(&undoRunState.RunProgram, wrapOptions{
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: domain.LocalBranchNames{config.previousBranch},
	})
	// If the command to undo failed and was continued,
	// there might be opcodes in the undo stack that became obsolete
	// when the command was continued.
	// Example: the command stashed away uncommitted changes,
	// failed, and remembered in the undo list to pop the stack.
	// When continuing, it finishes and pops the stack as part of the continue list.
	// When we run undo now, it still wants to pop the stack even though that was already done.
	// This seems to apply only to popping the stack and switching back to the initial branch.
	// Hence we consolidate these opcode types here.
	undoRunState.RunProgram.MoveToEnd(&opcode.RestoreOpenChanges{})
	undoRunState.RunProgram.RemoveAllButLast("*opcode.CheckoutIfExists")
	return undoRunState, err
}
