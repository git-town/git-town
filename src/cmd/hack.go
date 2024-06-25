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
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const hackDesc = "Create a new feature branch off the main branch"

const hackHelp = `
Syncs the main branch, forks a new feature branch with the given name off the main branch, pushes the new feature branch to origin (if and only if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func hackCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: "basic",
		Args:    cobra.ArbitraryArgs,
		Short:   hackDesc,
		Long:    cmdhelpers.Long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeHack(args, readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeHack(args []string, dryRun, verbose bool) error {
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
	data, exit, err := determineHackData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	appendData, doAppend, makeFeatureBranchData, doMakeFeatureBranch := data.Get()
	if doAppend {
		return createFeatureBranch(createFeatureBranchArgs{
			appendData:            appendData,
			backend:               repo.Backend,
			beginBranchesSnapshot: appendData.branchesSnapshot,
			beginConfigSnapshot:   repo.ConfigSnapshot,
			beginStashSize:        appendData.stashSize,
			commandsCounter:       repo.CommandsCounter,
			dryRun:                dryRun,
			finalMessages:         repo.FinalMessages,
			frontend:              repo.Frontend,
			git:                   repo.Git,
			rootDir:               repo.RootDir,
			verbose:               verbose,
		})
	}
	if doMakeFeatureBranch {
		return makeFeatureBranch(makeFeatureBranchArgs{
			beginConfigSnapshot: repo.ConfigSnapshot,
			config:              makeFeatureBranchData.config,
			makeFeatureData:     makeFeatureBranchData,
			repo:                repo,
			rootDir:             repo.RootDir,
			verbose:             verbose,
		})
	}
	panic("both config arms were nil")
}

// If set to appendData, the user wants to append a new branch to an existing branch.
// If set to makeFeatureData, the user wants to make an existing branch a feature branch.
type hackData = Either[appendFeatureData, makeFeatureData]

// this configuration is for when "git hack" is used to make contribution, observed, or parked branches feature branches
type makeFeatureData struct {
	config         config.ValidatedConfig
	targetBranches commandconfig.BranchesAndTypes
}

func createFeatureBranch(args createFeatureBranchArgs) error {
	runState := runstate.RunState{
		BeginBranchesSnapshot: args.beginBranchesSnapshot,
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		BeginStashSize:        args.beginStashSize,
		Command:               "hack",
		DryRun:                args.dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            appendProgram(args.appendData),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 args.backend,
		CommandsCounter:         args.commandsCounter,
		Config:                  args.appendData.config,
		Connector:               None[hostingdomain.Connector](),
		DialogTestInputs:        args.appendData.dialogTestInputs,
		FinalMessages:           args.finalMessages,
		Frontend:                args.frontend,
		Git:                     args.git,
		HasOpenChanges:          args.appendData.hasOpenChanges,
		InitialBranch:           args.appendData.initialBranch,
		InitialBranchesSnapshot: args.beginBranchesSnapshot,
		InitialConfigSnapshot:   args.beginConfigSnapshot,
		InitialStashSize:        args.beginStashSize,
		RootDir:                 args.rootDir,
		RunState:                runState,
		Verbose:                 args.verbose,
	})
}

type createFeatureBranchArgs struct {
	appendData            appendFeatureData
	backend               gitdomain.RunnerQuerier
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	beginConfigSnapshot   undoconfig.ConfigSnapshot
	beginStashSize        gitdomain.StashSize
	commandsCounter       Mutable[gohacks.Counter]
	dryRun                bool
	finalMessages         stringslice.Collector
	frontend              gitdomain.Runner
	git                   git.Commands
	rootDir               gitdomain.RepoRootDir
	verbose               bool
}

func determineHackData(args []string, repo execute.OpenRepoResult, dryRun, verbose bool) (data hackData, exit bool, err error) {
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
		data = Right[appendFeatureData, makeFeatureData](makeFeatureData{
			config:         validatedConfig,
			targetBranches: commandconfig.NewBranchesAndTypes(branchesToValidate, validatedConfig.Config),
		})
		return data, false, nil
	}
	if len(targetBranches) > 1 {
		err = errors.New(messages.HackTooManyArguments)
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
	data = Left[appendFeatureData, makeFeatureData](appendFeatureData{
		allBranches:               branchesSnapshot.Branches,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		config:                    validatedConfig,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		newBranchParentCandidates: gitdomain.LocalBranchNames{validatedConfig.Config.MainBranch},
		previousBranch:            previousBranch,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	})
	return
}

func makeFeatureBranch(args makeFeatureBranchArgs) error {
	err := validateMakeFeatureData(args.makeFeatureData)
	if err != nil {
		return err
	}
	for branchName, branchType := range args.makeFeatureData.targetBranches {
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
		fmt.Printf(messages.HackBranchIsNowFeature, branchName)
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:             args.repo.Backend,
		BeginConfigSnapshot: args.beginConfigSnapshot,
		Command:             "observe",
		CommandsCounter:     args.repo.CommandsCounter,
		EndConfigSnapshot:   undoconfig.EmptyConfigSnapshot(),
		FinalMessages:       args.repo.FinalMessages,
		RootDir:             args.rootDir,
		Verbose:             args.verbose,
	})
}

type makeFeatureBranchArgs struct {
	beginConfigSnapshot undoconfig.ConfigSnapshot
	config              config.ValidatedConfig
	makeFeatureData     makeFeatureData
	repo                execute.OpenRepoResult
	rootDir             gitdomain.RepoRootDir
	verbose             bool
}

func validateMakeFeatureData(data makeFeatureData) error {
	for branchName, branchType := range data.targetBranches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
			return nil
		case configdomain.BranchTypeFeatureBranch:
			return fmt.Errorf(messages.HackBranchIsAlreadyFeature, branchName)
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.HackCannotFeatureMainBranch)
		case configdomain.BranchTypePerennialBranch:
			return fmt.Errorf(messages.HackCannotFeaturePerennialBranch, branchName)
		}
		panic(fmt.Sprintf("unhandled branch type: %s", branchType))
	}
	return nil
}
