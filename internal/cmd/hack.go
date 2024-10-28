package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/cmd/sync"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/gohacks"
	"github.com/git-town/git-town/v16/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	configInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/config"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const hackDesc = "Create a new feature branch off the main branch"

const hackHelp = `
Syncs the main branch, forks a new feature branch with the given name off the main branch, pushes the new feature branch to origin (if and only if "push-new-branches" is true), and brings over all uncommitted changes to the new feature branch.

See "sync" for information regarding upstream remotes.`

func hackCmd() *cobra.Command {
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: "basic",
		Args:    cobra.ArbitraryArgs,
		Short:   hackDesc,
		Long:    cmdhelpers.Long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeHack(args, readDetachedFlag(cmd), readDryRunFlag(cmd), readPrototypeFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addPrototypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeHack(args []string, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineHackData(args, repo, detached, dryRun, prototype, verbose)
	if err != nil || exit {
		return err
	}
	createNewFeatureBranchData, doCreateNewFeatureBranch, convertToFeatureBranchData, doConvertToFeatureBranch := data.Get()
	if doCreateNewFeatureBranch {
		return createFeatureBranch(createFeatureBranchArgs{
			appendData:            createNewFeatureBranchData,
			backend:               repo.Backend,
			beginBranchesSnapshot: createNewFeatureBranchData.branchesSnapshot,
			beginConfigSnapshot:   repo.ConfigSnapshot,
			beginStashSize:        createNewFeatureBranchData.stashSize,
			commandsCounter:       repo.CommandsCounter,
			dryRun:                dryRun,
			finalMessages:         repo.FinalMessages,
			frontend:              repo.Frontend,
			git:                   repo.Git,
			rootDir:               repo.RootDir,
			verbose:               verbose,
		})
	}
	if doConvertToFeatureBranch {
		return convertToFeatureBranch(convertToFeatureBranchArgs{
			beginConfigSnapshot: repo.ConfigSnapshot,
			config:              convertToFeatureBranchData.config,
			makeFeatureData:     convertToFeatureBranchData,
			repo:                repo,
			rootDir:             repo.RootDir,
			verbose:             verbose,
		})
	}
	panic("both config arms were nil")
}

// If set to createNewFeatureData, the user wants to create a new feature branch.
// If set to convertToFeatureData, the user wants to convert an already existing branch into a feature branch.
type hackData = Either[appendFeatureData, convertToFeatureData]

// this configuration is for when "git town hack" is used to make contribution, observed, or parked branches feature branches
type convertToFeatureData struct {
	config         config.ValidatedConfig
	targetBranches configdomain.BranchesAndTypes
}

func createFeatureBranch(args createFeatureBranchArgs) error {
	runProgram := appendProgram(args.appendData)
	runState := runstate.RunState{
		BeginBranchesSnapshot: args.beginBranchesSnapshot,
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		BeginStashSize:        args.beginStashSize,
		Command:               "hack",
		DryRun:                args.dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
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
	dryRun                configdomain.DryRun
	finalMessages         stringslice.Collector
	frontend              gitdomain.Runner
	git                   git.Commands
	rootDir               gitdomain.RepoRootDir
	verbose               configdomain.Verbose
}

func determineHackData(args []string, repo execute.OpenRepoResult, detached configdomain.Detached, dryRun configdomain.DryRun, prototype configdomain.Prototype, verbose configdomain.Verbose) (data hackData, exit bool, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	targetBranches := gitdomain.NewLocalBranchNames(args...)
	var repoStatus gitdomain.RepoStatus
	repoStatus, err = repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
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
		return data, exit, err
	}
	localBranchNames := branchesSnapshot.Branches.LocalBranches().Names()
	var branchesToValidate gitdomain.LocalBranchNames
	shouldCreateBranch := len(targetBranches) == 1 && !slices.Contains(localBranchNames, targetBranches[0])
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
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
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, gitdomain.RemoteOrigin, print.Logger{})
	if err != nil {
		return data, false, err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(localBranchNames)
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: branchesToValidate,
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranchNames,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	if !shouldCreateBranch {
		data = Right[appendFeatureData, convertToFeatureData](convertToFeatureData{
			config:         validatedConfig,
			targetBranches: validatedConfig.BranchesAndTypes(branchesToValidate),
		})
		return data, false, nil
	}
	if len(targetBranches) > 1 {
		return data, false, errors.New(messages.HackTooManyArguments)
	}
	targetBranch := targetBranches[0]
	var remotes gitdomain.Remotes
	remotes, err = repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branchNamesToSync := gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch}
	if detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchesToSync, err := sync.BranchesToSync(branchNamesToSync, branchesSnapshot, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	data = Left[appendFeatureData, convertToFeatureData](appendFeatureData{
		branchInfos:               branchesSnapshot.Branches,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		config:                    validatedConfig,
		dialogTestInputs:          dialogTestInputs,
		dryRun:                    dryRun,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		newBranchParentCandidates: gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch},
		preFetchBranchInfos:       preFetchBranchSnapshot.Branches,
		previousBranch:            previousBranch,
		prototype:                 prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	})
	return data, false, err
}

func convertToFeatureBranch(args convertToFeatureBranchArgs) error {
	err := validateConvertToFeatureData(args.makeFeatureData)
	if err != nil {
		return err
	}
	for branchName, branchType := range args.makeFeatureData.targetBranches {
		switch branchType {
		case configdomain.BranchTypeContributionBranch:
			err = args.config.NormalConfig.RemoveFromContributionBranches(branchName)
		case configdomain.BranchTypeObservedBranch:
			err = args.config.NormalConfig.RemoveFromObservedBranches(branchName)
		case configdomain.BranchTypeParkedBranch:
			err = args.config.NormalConfig.RemoveFromParkedBranches(branchName)
		case configdomain.BranchTypePrototypeBranch:
			err = args.config.NormalConfig.RemoveFromPrototypeBranches(branchName)
		case
			configdomain.BranchTypeFeatureBranch,
			configdomain.BranchTypeMainBranch,
			configdomain.BranchTypePerennialBranch:
			panic(fmt.Sprintf("unchecked branch type: %s", branchType))
		}
		if err != nil {
			return err
		}
		fmt.Printf(messages.HackBranchIsNowFeature, branchName)
	}
	return configInterpreter.Finished(configInterpreter.FinishedArgs{
		Backend:               args.repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		Command:               "observe",
		CommandsCounter:       args.repo.CommandsCounter,
		FinalMessages:         args.repo.FinalMessages,
		Git:                   args.repo.Git,
		RootDir:               args.rootDir,
		TouchedBranches:       args.makeFeatureData.targetBranches.Keys().BranchNames(),
		Verbose:               args.verbose,
	})
}

type convertToFeatureBranchArgs struct {
	beginConfigSnapshot undoconfig.ConfigSnapshot
	config              config.ValidatedConfig
	makeFeatureData     convertToFeatureData
	repo                execute.OpenRepoResult
	rootDir             gitdomain.RepoRootDir
	verbose             configdomain.Verbose
}

func validateConvertToFeatureData(data convertToFeatureData) error {
	for branchName, branchType := range data.targetBranches {
		switch branchType {
		case
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeObservedBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
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
