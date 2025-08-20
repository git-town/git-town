package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/cmd/sync"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/configinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	hackDesc = "Create a new feature branch off the main branch"
	hackHelp = `
Consider this stack:

main
 \
  branch-1
   \
*   branch-2

We are on the "branch-2" branch. After running "git hack branch-3", our
workspace contains these branches:

main
 \
  branch-1
   \
    branch-2
 \
* branch-3

The new branch "feature-2"
is a child of the main branch.

If there are no uncommitted changes,
it also syncs all affected branches.
`
)

func hackCmd() *cobra.Command {
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addBeamFlag, readBeamFlag := flags.Beam()
	addCommitFlag, readCommitFlag := flags.Commit()
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.ArbitraryArgs,
		Short:   hackDesc,
		Long:    cmdhelpers.Long(hackDesc, hackHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			beam, errBeam := readBeamFlag(cmd)
			commit, errCommit := readCommitFlag(cmd)
			commitMessage, errCommitMessage := readCommitMessageFlag(cmd)
			detached, errDetached := readDetachedFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			propose, errPropose := readProposeFlag(cmd)
			prototype, errPrototype := readPrototypeFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAutoResolve, errBeam, errCommit, errCommitMessage, errDetached, errDryRun, errPropose, errPrototype, errVerbose); err != nil {
				return err
			}
			if commitMessage.IsSome() || propose.IsTrue() {
				commit = true
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: autoResolve,
				Detached:    detached,
				DryRun:      dryRun,
				Verbose:     verbose,
			})
			return executeHack(hackArgs{
				argv:          args,
				beam:          beam,
				cliConfig:     cliConfig,
				commit:        commit,
				commitMessage: commitMessage,
				propose:       propose,
				prototype:     prototype,
			})
		},
	}
	addBeamFlag(&cmd)
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addAutoResolveFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type hackArgs struct {
	argv          []string
	beam          configdomain.Beam
	cliConfig     configdomain.PartialConfig
	commit        configdomain.Commit
	commitMessage Option[gitdomain.CommitMessage]
	propose       configdomain.Propose
	prototype     configdomain.Prototype
}

func executeHack(args hackArgs) error {
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
	data, exit, err := determineHackData(args, repo)
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
			branchInfosLastRun:    createNewFeatureBranchData.branchInfosLastRun,
			commandsCounter:       repo.CommandsCounter,
			dryRun:                createNewFeatureBranchData.config.NormalConfig.DryRun,
			finalMessages:         repo.FinalMessages,
			frontend:              repo.Frontend,
			git:                   repo.Git,
			rootDir:               repo.RootDir,
		})
	}
	if doConvertToFeatureBranch {
		return convertToFeatureBranch(repo, convertToFeatureBranchArgs{
			beginConfigSnapshot: repo.ConfigSnapshot,
			config:              convertToFeatureBranchData.config,
			makeFeatureData:     convertToFeatureBranchData,
			verbose:             convertToFeatureBranchData.config.NormalConfig.Verbose,
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

type createFeatureBranchArgs struct {
	appendData            appendFeatureData
	backend               subshelldomain.RunnerQuerier
	beginBranchesSnapshot gitdomain.BranchesSnapshot
	beginConfigSnapshot   configdomain.BeginConfigSnapshot
	beginStashSize        gitdomain.StashSize
	branchInfosLastRun    Option[gitdomain.BranchInfos]
	commandsCounter       Mutable[gohacks.Counter]
	dryRun                configdomain.DryRun
	finalMessages         stringslice.Collector
	frontend              subshelldomain.Runner
	git                   git.Commands
	rootDir               gitdomain.RepoRootDir
}

func createFeatureBranch(args createFeatureBranchArgs) error {
	runProgram := appendProgram(args.backend, args.appendData, args.finalMessages, true)
	runState := runstate.RunState{
		BeginBranchesSnapshot: args.beginBranchesSnapshot,
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		BeginStashSize:        args.beginStashSize,
		BranchInfosLastRun:    args.branchInfosLastRun,
		Command:               "hack",
		DryRun:                args.dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 args.backend,
		CommandsCounter:         args.commandsCounter,
		Config:                  args.appendData.config,
		Connector:               args.appendData.connector,
		FinalMessages:           args.finalMessages,
		Frontend:                args.frontend,
		Git:                     args.git,
		HasOpenChanges:          args.appendData.hasOpenChanges,
		InitialBranch:           args.appendData.initialBranch,
		InitialBranchesSnapshot: args.beginBranchesSnapshot,
		InitialConfigSnapshot:   args.beginConfigSnapshot,
		InitialStashSize:        args.beginStashSize,
		Inputs:                  args.appendData.inputs,
		PendingCommand:          None[string](),
		RootDir:                 args.rootDir,
		RunState:                runState,
	})
}

func determineHackData(args hackArgs, repo execute.OpenRepoResult) (data hackData, exit dialogdomain.Exit, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	targetBranches := gitdomain.NewLocalBranchNames(args.argv...)
	var repoStatus gitdomain.RepoStatus
	repoStatus, err = repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
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
		return data, false, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 len(args.argv) == 1 && !repoStatus.OpenChanges && args.beam.IsFalse() && args.commit.IsFalse(),
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
		return data, exit, err
	}
	localBranchNames := branchesSnapshot.Branches.LocalBranches().Names()
	var branchesToValidate gitdomain.LocalBranchNames
	shouldCreateBranch := len(targetBranches) == 1 && !slices.Contains(localBranchNames, targetBranches[0])
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
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
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(localBranchNames)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: branchesToValidate,
		ConfigSnapshot:     repo.ConfigSnapshot,
		Connector:          connector,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		Inputs:             inputs,
		LocalBranches:      localBranchNames,
		Remotes:            remotes,
		RepoStatus:         repoStatus,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	if !shouldCreateBranch {
		data = Right[appendFeatureData](convertToFeatureData{
			config:         validatedConfig,
			targetBranches: validatedConfig.BranchesAndTypes(branchesToValidate),
		})
		return data, false, nil
	}
	if len(targetBranches) > 1 {
		return data, false, errors.New(messages.HackTooManyArguments)
	}
	targetBranch := targetBranches[0]
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch)
	}
	branchNamesToSync := gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch}
	if validatedConfig.NormalConfig.Detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	commitsToBeam := []gitdomain.Commit{}
	ancestor, hasAncestor := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage).Get()
	if args.beam.IsTrue() && hasAncestor {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor.BranchName())
		if err != nil {
			return data, false, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, targetBranch, repo.Git, repo.Backend, inputs)
		if err != nil || exit {
			return data, exit, err
		}
	}
	if validatedConfig.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPropose {
		args.propose = true
	}
	data = Left[appendFeatureData, convertToFeatureData](appendFeatureData{
		beam:                      args.beam,
		branchInfos:               branchesSnapshot.Branches,
		branchInfosLastRun:        branchInfosLastRun,
		branchesSnapshot:          branchesSnapshot,
		branchesToSync:            branchesToSync,
		commit:                    args.commit,
		commitMessage:             args.commitMessage,
		commitsToBeam:             commitsToBeam,
		config:                    validatedConfig,
		connector:                 connector,
		hasOpenChanges:            repoStatus.OpenChanges,
		initialBranch:             initialBranch,
		initialBranchInfo:         initialBranchInfo,
		inputs:                    inputs,
		newBranchParentCandidates: gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch},
		nonExistingBranches:       nonExistingBranches,
		preFetchBranchInfos:       preFetchBranchSnapshot.Branches,
		previousBranch:            previousBranch,
		propose:                   args.propose,
		prototype:                 args.prototype,
		remotes:                   remotes,
		stashSize:                 stashSize,
		targetBranch:              targetBranch,
	})
	return data, false, err
}

type convertToFeatureBranchArgs struct {
	beginConfigSnapshot configdomain.BeginConfigSnapshot
	config              config.ValidatedConfig
	makeFeatureData     convertToFeatureData
	verbose             configdomain.Verbose
}

func convertToFeatureBranch(repo execute.OpenRepoResult, args convertToFeatureBranchArgs) error {
	for branchName, branchType := range args.makeFeatureData.targetBranches {
		switch branchType {
		case
			configdomain.BranchTypeContributionBranch,
			configdomain.BranchTypeObservedBranch,
			configdomain.BranchTypeParkedBranch,
			configdomain.BranchTypePrototypeBranch:
			if err := gitconfig.SetBranchTypeOverride(repo.Backend, configdomain.BranchTypeFeatureBranch, branchName); err != nil {
				return err
			}
		case configdomain.BranchTypeFeatureBranch:
			return fmt.Errorf(messages.HackBranchIsAlreadyFeature, branchName)
		case configdomain.BranchTypeMainBranch:
			return errors.New(messages.HackCannotFeatureMainBranch)
		case configdomain.BranchTypePerennialBranch:
			return fmt.Errorf(messages.HackCannotFeaturePerennialBranch, branchName)
		}
		fmt.Printf(messages.BranchIsNowFeature, branchName)
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   args.beginConfigSnapshot,
		Command:               "observe",
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       args.makeFeatureData.targetBranches.Keys().BranchNames(),
		Verbose:               args.verbose,
	})
}
