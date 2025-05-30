package cmd

import (
	"errors"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/optimizer"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	walkCmd  = "walk"
	walkDesc = "Run a command on each local feature branch"
	walkHelp = `
Executes the given command on each local feature branch in stack order.
Stops if the command exits with an error,
giving you a chance to investigate and fix the issue.

* use "git town continue" to retry the command on the current branch
* use "git town skip" to move on to the next branch
* use "git town undo" to abort the iteration and undo all changes made
* use "git town status reset" to abort the iteration and keep all changes made

If no shell command is provided, drops you into an interactive shell for each branch.
You can manually run any shell commands,
then proceed to the next branch with "git town continue".

Consider this stack:

main
	\
   branch-1
	 	\
     branch-2
		 	\
       branch-3

Running "git town walk --stack make lint" produces this output:

[branch-1] make lint
... output of make lint

[branch-2] make lint
... output of make lint

[branch-3] make lint
... output of make lint
`
)

func walkCommand() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("iterate all local branches")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addStackFlag, readStackFlag := flags.Stack("iterate all branches in the current stack")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     walkCmd,
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   walkDesc,
		Long:    cmdhelpers.Long(walkDesc, walkHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			allBranches, err := readAllFlag(cmd)
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
			return executeWalk(args, dryRun, allBranches, stack, verbose)
		},
	}
	addAllFlag(&cmd)
	addDryRunFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeWalk(args []string, dryRun configdomain.DryRun, allBranches configdomain.AllBranches, fullStack configdomain.FullStack, verbose configdomain.Verbose) error {
	if len(args) == 0 && dryRun {
		return errors.New(messages.WalkNoDryRun)
	}
	err := validateArgs(allBranches, fullStack)
	if err != nil {
		return err
	}
	data, exit, err := determineWalkData(allBranches, dryRun, fullStack, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := walkProgram(args, data, dryRun)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   data.repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               walkCmd,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 data.repo.Backend,
		CommandsCounter:         data.repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		Detached:                true,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           data.repo.FinalMessages,
		Frontend:                data.repo.Frontend,
		Git:                     data.repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   data.repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		PendingCommand:          None[string](),
		RootDir:                 data.repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type walkData struct {
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchesToWalk     gitdomain.LocalBranchNames
	config             config.ValidatedConfig
	connector          Option[forgedomain.Connector]
	dialogTestInputs   components.TestInputs
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	previousBranch     Option[gitdomain.LocalBranchName]
	repo               execute.OpenRepoResult
	stashSize          gitdomain.StashSize
}

func determineWalkData(all configdomain.AllBranches, dryRun configdomain.DryRun, stack configdomain.FullStack, verbose configdomain.Verbose) (walkData, bool, error) {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return walkData{}, false, err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return walkData{}, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Detached:              true,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
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
		return walkData{}, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return walkData{}, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connector, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return walkData{}, false, err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
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
		return walkData{}, exit, err
	}
	perennialBranchNames := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	branchesToWalk := gitdomain.LocalBranchNames{}
	switch {
	case all.Enabled():
		branchesToWalk = localBranches.Remove(perennialBranchNames...)
	case stack.Enabled():
		branchesToWalk = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialBranchNames)
	}
	return walkData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		branchesToWalk:     branchesToWalk,
		config:             validatedConfig,
		connector:          connector,
		dialogTestInputs:   dialogTestInputs,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		previousBranch:     previousBranch,
		repo:               repo,
		stashSize:          stashSize,
	}, false, err
}

func walkProgram(args []string, data walkData, dryRun configdomain.DryRun) program.Program {
	prog := NewMutable(&program.Program{})
	hasCall, executable, callArgs := parseArgs(args)
	for _, branchToWalk := range data.branchesToWalk {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: branchToWalk})
		if hasCall {
			prog.Value.Add(
				&opcodes.ExecuteShellCommand{
					Executable: executable,
					Args:       callArgs,
				},
			)
		} else {
			prog.Value.Add(
				&opcodes.ExitToShell{},
			)
		}
		prog.Value.Add(
			&opcodes.ProgramEndOfBranch{},
		)
	}
	prog.Value.Add(
		&opcodes.CheckoutIfNeeded{
			Branch: data.initialBranch,
		},
		&opcodes.MessageQueue{
			Message: messages.WalkDone,
		},
	)
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   dryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return optimizer.Optimize(prog.Immutable())
}

func validateArgs(all configdomain.AllBranches, stack configdomain.FullStack) error {
	if all.Enabled() == stack.Enabled() {
		return errors.New(messages.WalkAllOrStack)
	}
	return nil
}

func parseArgs(args []string) (bool, string, []string) {
	if len(args) == 0 {
		return false, "", []string{}
	}
	return true, args[0], args[1:]
}
