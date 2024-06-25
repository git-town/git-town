package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/commandconfig"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	"github.com/git-town/git-town/v14/src/validate"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const prototypeDesc = "Creates prototype branches"

const prototypeHelp = `
A prototype branch is a local-only feature branch that incorporates updates from its parent branch but is not pushed to the remote repository.
`

func prototypeCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "prototype <branch>",
		GroupID: "types",
		Args:    cobra.ArbitraryArgs,
		Short:   prototypeDesc,
		Long:    cmdhelpers.Long(prototypeDesc, prototypeHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePrototype(args, readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executePrototype(args []string, dryRun, verbose bool) error {
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
	data, exit, err := determinePrototypeData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	createData, doCreate, makePrototypeBranchData, doMakePrototypeBranch := data.Get()
	if doCreate {
		return createPrototypeBranch(appendPrototypeData{
			allBranches:               createData.allBranches,
			backend:                   repo.Backend,
			beginBranchesSnapshot:     createData.beginBranchesSnapshot,
			beginConfigSnapshot:       createData.beginConfigSnapshot,
			beginStashSize:            createData.beginStashSize,
			branchesToSync:            createData.branchesToSync,
			commandsCounter:           repo.CommandsCounter,
			config:                    config.ValidatedConfig{},
			dialogTestInputs:          createData.dialogTestInputs,
			dryRun:                    dryRun,
			finalMessages:             stringslice.Collector{},
			frontend:                  nil,
			git:                       git.Commands{},
			hasOpenChanges:            false,
			initialBranch:             "",
			newBranchParentCandidates: []gitdomain.LocalBranchName{},
			previousBranch:            createData.previousBranch,
			remotes:                   []gitdomain.Remote{},
			rootDir:                   "",
			targetBranch:              "",
			verbose:                   verbose,
		})
	}
	if doMakePrototypeBranch {
		return makePrototypeBranch(makePrototypeData{
			config:         makePrototypeBranchData.config,
			configSnapshot: repo.ConfigSnapshot,
			repo:           repo,
			rootDir:        repo.RootDir,
			verbose:        verbose,
		})
	}
	panic("both config arms were nil")
}

// If set to appendPrototypeData, the user wants to append a new prototype branch to an existing branch.
// If set to makePrototypeData, the user wants to make an existing branch a prototype branch.
type prototypeData = Either[appendPrototypeData, makePrototypeData]

type appendPrototypeData struct {
	allBranches               gitdomain.BranchInfos
	backend                   gitdomain.RunnerQuerier
	beginBranchesSnapshot     gitdomain.BranchesSnapshot
	beginConfigSnapshot       undoconfig.ConfigSnapshot
	beginStashSize            gitdomain.StashSize
	branchesToSync            gitdomain.BranchInfos
	commandsCounter           Mutable[gohacks.Counter]
	config                    config.ValidatedConfig
	dialogTestInputs          Mutable[components.TestInputs]
	dryRun                    bool
	finalMessages             stringslice.Collector
	frontend                  gitdomain.Runner
	git                       git.Commands
	hasOpenChanges            bool
	initialBranch             gitdomain.LocalBranchName
	newBranchParentCandidates gitdomain.LocalBranchNames
	previousBranch            Option[gitdomain.LocalBranchName]
	remotes                   gitdomain.Remotes
	rootDir                   gitdomain.RepoRootDir
	targetBranch              gitdomain.LocalBranchName
	verbose                   bool
}

// this configuration is for when "git prototype" is used to make contribution, observed, or parked branches prototype branches
type makePrototypeData struct {
	config         config.ValidatedConfig
	configSnapshot undoconfig.ConfigSnapshot
	repo           execute.OpenRepoResult
	rootDir        gitdomain.RepoRootDir
	targetBranches commandconfig.BranchesAndTypes
	verbose        bool
}

func createPrototypeBranch(args appendPrototypeData) error {
	program := appendProgram(appendFeatureData{
		allBranches:               args.allBranches,
		branchesSnapshot:          args.beginBranchesSnapshot,
		branchesToSync:            args.branchesToSync,
		config:                    args.config,
		dialogTestInputs:          args.dialogTestInputs,
		dryRun:                    args.dryRun,
		hasOpenChanges:            args.hasOpenChanges,
		initialBranch:             args.initialBranch,
		newBranchParentCandidates: args.newBranchParentCandidates,
		previousBranch:            args.previousBranch,
		remotes:                   args.remotes,
		stashSize:                 args.beginStashSize,
		targetBranch:              args.targetBranch,
	})
	program.Add(&opcodes.AddToPrototypeBranches{})
	runState := runstate.RunState{
		BeginBranchesSnapshot: args.beginBranchesSnapshot,
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		BeginStashSize:        args.beginStashSize,
		Command:               "prototype",
		DryRun:                args.dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            program,
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 args.backend,
		CommandsCounter:         args.commandsCounter,
		Config:                  args.config,
		Connector:               None[hostingdomain.Connector](),
		DialogTestInputs:        args.dialogTestInputs,
		FinalMessages:           args.finalMessages,
		Frontend:                args.frontend,
		Git:                     args.git,
		HasOpenChanges:          args.hasOpenChanges,
		InitialBranch:           args.initialBranch,
		InitialBranchesSnapshot: args.beginBranchesSnapshot,
		InitialConfigSnapshot:   args.beginConfigSnapshot,
		InitialStashSize:        args.beginStashSize,
		RootDir:                 args.rootDir,
		RunState:                runState,
		Verbose:                 args.verbose,
	})
}

func determinePrototypeData(args []string, repo execute.OpenRepoResult, dryRun, verbose bool) (data prototypeData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	targetBranches := gitdomain.NewLocalBranchNames(args...)
	var repoStatus gitdomain.RepoStatus
	repoStatus, err = repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 len(args) == 1 && !repoStatus.OpenChanges,
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
		return
	}
	localBranchNames := branchesSnapshot.Branches.LocalBranches().Names()
	var branchesToValidate gitdomain.LocalBranchNames
	shouldCreateBranch := len(targetBranches) == 1 && !slice.Contains(localBranchNames, targetBranches[0])
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		err = errors.New(messages.CurrentBranchCannotDetermine)
		return
	}
	if shouldCreateBranch {
		branchesToValidate = gitdomain.LocalBranchNames{}
	} else {
		if len(targetBranches) == 0 {
			branchesToValidate = gitdomain.LocalBranchNames{initialBranch}
		} else {
			branchesToValidate = targetBranches
		}
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchesToValidate,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranchNames,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	if !shouldCreateBranch {
		data = Right[appendPrototypeData, makePrototypeData](makePrototypeData{
			config:         validatedConfig,
			targetBranches: commandconfig.NewBranchesAndTypes(branchesToValidate, validatedConfig.Config),
		})
		return
	}
	targetBranch := targetBranches[0]
	var remotes gitdomain.Remotes
	remotes, err = repo.Git.Remotes(repo.Backend)
	if err != nil {
		return
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		err = fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
		return
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		err = fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
		return
	}
	branchNamesToSync := gitdomain.LocalBranchNames{validatedConfig.Config.MainBranch}
	var branchesToSync gitdomain.BranchInfos
	branchesToSync, err = branchesSnapshot.Branches.Select(branchNamesToSync...)
	data = Left[appendPrototypeData, makePrototypeData](appendPrototypeData{
		allBranches:               branchesSnapshot.Branches,
		beginBranchesSnapshot:     branchesSnapshot,
		beginConfigSnapshot:       repo.ConfigSnapshot,
		beginStashSize:            stashSize,
		branchesToSync:            branchesToSync,
		config:                    validatedConfig,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		newBranchParentCandidates: gitdomain.LocalBranchNames{validatedConfig.Config.MainBranch},
		previousBranch:            previousBranch,
		remotes:                   remotes,
		targetBranch:              targetBranch,
	})
	return
}

func makePrototypeBranch(args makePrototypeData) error {
	err := validateMakePrototypeData(args)
	if err != nil {
		return err
	}
	for branchName, branchType := range args.targetBranches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			err = args.config.RemoveFromContributionBranches(branchName)
		case configdomain.BranchTypeObservedBranch:
			err = args.config.RemoveFromObservedBranches(branchName)
		case configdomain.BranchTypeParkedBranch:
			err = args.config.RemoveFromParkedBranches(branchName)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
			panic(fmt.Sprintf("unchecked branch type: %s", branchType))
		}
		if err != nil {
			return err
		}
		fmt.Printf(messages.PrototypeBranchIsNowPrototype, branchName)
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             args.repo.Backend,
		BeginConfigSnapshot: args.configSnapshot,
		Command:             "observe",
		CommandsCounter:     args.repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       args.repo.FinalMessages,
		RootDir:             args.rootDir,
		Verbose:             args.verbose,
	})
}

type makePrototypeBranchArgs struct {
	beginConfigSnapshot undoconfig.ConfigSnapshot
	config              config.ValidatedConfig
	makeFeatureData     makeFeatureData
	repo                execute.OpenRepoResult
	rootDir             gitdomain.RepoRootDir
	verbose             bool
}

func validateMakePrototypeData(data makePrototypeData) error {
	for branchName, branchType := range data.targetBranches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
			return nil
		case configdomain.BranchTypePrototypeBranch:
			return fmt.Errorf(messages.PrototypeBranchIsAlreadyPrototype, branchName)
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.PrototypeCannotPrototypeMainBranch)
		case configdomain.BranchTypePerennialBranch:
			return fmt.Errorf(messages.PrototypeCannotPrototypePerennialBranch, branchName)
		}
		panic(fmt.Sprintf("unhandled branch type: %s", branchType))
	}
	return nil
}
