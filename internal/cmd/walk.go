package cmd

import (
	"cmp"
	"errors"
	"os"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
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
			allBranches, errAllBranches := readAllFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			stack, errStack := readStackFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAllBranches, errDryRun, errStack, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				DryRun:      dryRun,
				Verbose:     verbose,
			})
			return executeWalk(executeWalkArgs{
				allBranches: allBranches,
				argv:        args,
				cliConfig:   cliConfig,
				stack:       stack,
			})
		},
	}
	addAllFlag(&cmd)
	addDryRunFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeWalkArgs struct {
	allBranches configdomain.AllBranches
	argv        []string
	cliConfig   configdomain.PartialConfig
	stack       configdomain.FullStack
}

func executeWalk(args executeWalkArgs) error {
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
	if len(args.argv) == 0 && repo.UnvalidatedConfig.NormalConfig.DryRun {
		return errors.New(messages.WalkNoDryRun)
	}
	if err := validateArgs(args.allBranches, args.stack); err != nil {
		return err
	}
	data, exit, err := determineWalkData(repo, args.allBranches, args.stack)
	if err != nil || exit {
		return err
	}
	runProgram := walkProgram(args.argv, data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               walkCmd,
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		Detached:                true,
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

type walkData struct {
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchesToWalk     gitdomain.LocalBranchNames
	config             config.ValidatedConfig
	connector          Option[forgedomain.Connector]
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	inputs             dialogcomponents.Inputs
	previousBranch     Option[gitdomain.LocalBranchName]
	stashSize          gitdomain.StashSize
}

func determineWalkData(repo execute.OpenRepoResult, all configdomain.AllBranches, stack configdomain.FullStack) (walkData, dialogdomain.Exit, error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return walkData{}, false, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		CodebergToken:        config.CodebergToken,
		ForgeType:            config.ForgeType,
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
		return walkData{}, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Detached:              true,
		Fetch:                 false,
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
		return walkData{}, exit, err
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return walkData{}, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return walkData{}, false, err
	}
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
		LocalBranches:      localBranches,
		Remotes:            remotes,
		RepoStatus:         repoStatus,
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
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		inputs:             inputs,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, false, nil
}

func walkProgram(args []string, data walkData) program.Program {
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
		DryRun:                   data.config.NormalConfig.DryRun,
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
