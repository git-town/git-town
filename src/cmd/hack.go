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
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	configInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/config"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/runstate"
	"github.com/spf13/cobra"
)

const hackDesc = "Create a new feature branch off the main development branch"

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
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineHackData(args, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	appendData, doAppend, makeFeatureBranchData, doMakeFeatureBranch := data.Get()
	if doAppend {
		return createBranch(createBranchArgs{
			appendData:            appendData,
			beginBranchesSnapshot: initialBranchesSnapshot,
			beginConfigSnapshot:   repo.ConfigSnapshot,
			beginStashSize:        initialStashSize,
			dryRun:                dryRun,
			rootDir:               repo.RootDir,
			runner:                appendData.runner,
			verbose:               verbose,
		})
	}
	if doMakeFeatureBranch {
		return makeFeatureBranch(makeFeatureBranchArgs{
			beginConfigSnapshot: repo.ConfigSnapshot,
			config:              repo.Config,
			makeFeatureData:     makeFeatureBranchData,
			repo:                repo,
			rootDir:             repo.RootDir,
			runner:              makeFeatureBranchData.runner,
			verbose:             verbose,
		})
	}
	panic("both config arms were nil")
}

// If set to appendConfig, the user wants to append a new branch to an existing branch.
// If set to makeFeatureConfig, the user wants to make an existing branch a feature branch.
type hackData = Either[appendData, makeFeatureData]

// this configuration is for when "git hack" is used to make contribution, observed, or parked branches feature branches
type makeFeatureData struct {
	runner         *git.ProdRunner
	targetBranches commandconfig.BranchesAndTypes
}

func createBranch(args createBranchArgs) error {
	runState := runstate.RunState{
		BeginBranchesSnapshot: args.beginBranchesSnapshot,
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		BeginStashSize:        args.beginStashSize,
		Command:               "hack",
		DryRun:                args.dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            appendProgram(args.appendData),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Config:                  args.appendData.config,
		Connector:               nil,
		DialogTestInputs:        &args.appendData.dialogTestInputs,
		HasOpenChanges:          args.appendData.hasOpenChanges,
		InitialBranchesSnapshot: args.beginBranchesSnapshot,
		InitialConfigSnapshot:   args.beginConfigSnapshot,
		InitialStashSize:        args.beginStashSize,
		RootDir:                 args.rootDir,
		Run:                     args.runner,
		RunState:                runState,
		Verbose:                 args.verbose,
	})
}

type createBranchArgs struct {
	appendData            appendData
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	beginConfigSnapshot   undoconfig.ConfigSnapshot
	beginStashSize        gitdomain.StashSize
	dryRun                bool
	rootDir               gitdomain.RepoRootDir
	runner                *git.ProdRunner
	verbose               bool
}

func determineHackData(args []string, repo *execute.OpenRepoResult, dryRun, verbose bool) (data hackData, branchesSnapshot gitdomain.BranchesSnapshot, stashSize gitdomain.StashSize, exit bool, err error) {
	fc := execute.FailureCollector{}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	var repoStatus gitdomain.RepoStatus
	repoStatus, err = repo.Backend.RepoStatus()
	if err != nil {
		return
	}
	branchesSnapshot, stashSize, exit, err = execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 len(args) == 1 && !repoStatus.OpenChanges,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return
	}
	previousBranch := repo.Backend.PreviouslyCheckedOutBranch()
	targetBranches := gitdomain.NewLocalBranchNames(args...)
	if len(targetBranches) == 0 {
		data = Right[appendData, makeFeatureData](makeFeatureData{
			targetBranches: commandconfig.NewBranchesAndTypes(gitdomain.LocalBranchNames{branchesSnapshot.Active}, repo.Config.Config),
		})
		return
	}
	if len(targetBranches) > 0 && branchesSnapshot.Branches.HasLocalBranches(targetBranches) {
		data = Right[appendData, makeFeatureData](makeFeatureData{
			targetBranches: commandconfig.NewBranchesAndTypes(targetBranches, repo.Config.Config),
		})
		return
	}
	if len(targetBranches) > 1 {
		err = errors.New(messages.HackTooManyArguments)
		return
	}
	targetBranch := targetBranches[0]
	remotes := fc.Remotes(repo.Backend.Remotes())
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		err = fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
		return
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		err = fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
		return
	}
	branchNamesToSync := gitdomain.LocalBranchNames{repo.Config.Config.MainBranch}
	branchesToSync := fc.BranchInfos(branchesSnapshot.Branches.Select(branchNamesToSync...))
	runner := git.ProdRunner{
		Backend:         repo.Backend,
		CommandsCounter: repo.CommandsCounter,
		Config:          repo.Config,
		FinalMessages:   repo.FinalMessages,
		Frontend:        repo.Frontend,
	}
	data = Left[appendData, makeFeatureData](appendData{
		allBranches:               branchesSnapshot.Branches,
		branchesToSync:            branchesToSync,
		config:                    repo.Config.Config,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             branchesSnapshot.Active,
		newBranchParentCandidates: gitdomain.LocalBranchNames{repo.Config.Config.MainBranch},
		parentBranch:              repo.Config.Config.MainBranch,
		previousBranch:            previousBranch,
		remotes:                   remotes,
		runner:                    &runner,
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
	config              *config.Config
	makeFeatureData     makeFeatureData
	repo                *execute.OpenRepoResult
	rootDir             gitdomain.RepoRootDir
	runner              *git.ProdRunner
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
