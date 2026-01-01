package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/sync"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	addCommitMessageFlag, readCommitMessageFlag := flags.CommitMessage("the commit message")
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addProposeFlag, readProposeFlag := flags.Propose()
	addPrototypeFlag, readPrototypeFlag := flags.Prototype()
	addStashFlag, readStashFlag := flags.Stash()
	addSyncFlag, readSyncFlag := flags.Sync()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "hack <branch>",
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.ExactArgs(1),
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
			stash, errStash := readStashFlag(cmd)
			sync, errSync := readSyncFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errAutoResolve, errBeam, errCommit, errCommitMessage, errDetached, errDryRun, errPropose, errPrototype, errStash, errSync, errVerbose); err != nil {
				return err
			}
			if commitMessage.IsSome() || propose.ShouldPropose() {
				commit = true
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       autoResolve,
				AutoSync:          sync,
				Detached:          detached,
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             stash,
				Verbose:           verbose,
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
	addAutoResolveFlag(&cmd)
	addBeamFlag(&cmd)
	addCommitFlag(&cmd)
	addCommitMessageFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addProposeFlag(&cmd)
	addPrototypeFlag(&cmd)
	addStashFlag(&cmd)
	addSyncFlag(&cmd)
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
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineHackData(args, repo)
	if err != nil {
		return err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit:
		return nil
	case configdomain.ProgramFlowRestart:
		goto Start
	}
	runProgram := appendProgram(repo.Backend, data, repo.FinalMessages, true)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               "hack",
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

func determineHackData(args hackArgs, repo execute.OpenRepoResult) (data appendFeatureData, flow configdomain.ProgramFlow, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	targetBranches := gitdomain.NewLocalBranchNames(args.argv...)
	var repoStatus gitdomain.RepoStatus
	repoStatus, err = repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
		Frontend:             repo.Frontend,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	fetch := true
	if repoStatus.OpenChanges {
		fetch = false
	}
	if args.beam.ShouldBeam() || args.commit.ShouldCommit() {
		fetch = false
	}
	if !config.AutoSync {
		fetch = false
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 fetch,
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
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return data, flow, nil
	}
	localBranchNames := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	initialBranchInfo, hasInitialBranchInfo := branchesSnapshot.Branches.FindByLocalName(initialBranch).Get()
	if !hasInitialBranchInfo {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	branchesToValidate := gitdomain.LocalBranchNames{}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(localBranchNames)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
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
		return data, configdomain.ProgramFlowExit, err
	}
	if len(targetBranches) > 1 {
		return data, configdomain.ProgramFlowExit, errors.New(messages.HackTooManyArguments)
	}
	targetBranch := targetBranches[0]
	if prefix, hasPrefix := validatedConfig.NormalConfig.BranchPrefix.Get(); hasPrefix {
		targetBranch = prefix.Apply(targetBranch)
	}
	if branchesSnapshot.Branches.HasLocalBranch(targetBranch) {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchAlreadyExistsLocally, targetBranch)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(targetBranch) {
		return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchAlreadyExistsRemotely, targetBranch, config.DevRemote)
	}
	branchNamesToSync := gitdomain.LocalBranchNames{validatedConfig.ValidatedConfigData.MainBranch}
	if validatedConfig.NormalConfig.Detached {
		branchNamesToSync = validatedConfig.RemovePerennials(branchNamesToSync)
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	commitsToBeam := []gitdomain.Commit{}
	ancestor, hasAncestor := latestExistingAncestor(initialBranch, branchesSnapshot.Branches, validatedConfig.NormalConfig.Lineage).Get()
	if args.beam.ShouldBeam() && !hasAncestor {
		// ask the user for the parent branch
		excludeBranches := append(
			gitdomain.LocalBranchNames{initialBranch},
			validatedConfig.NormalConfig.Lineage.Children(initialBranch, validatedConfig.NormalConfig.Order)...,
		)
		noneEntry := dialog.SwitchBranchEntry{
			Branch:        messages.SetParentNoneOption,
			Indentation:   "",
			OtherWorktree: false,
			Type:          configdomain.BranchTypeFeatureBranch,
		}
		entriesArgs := dialog.NewSwitchBranchEntriesArgs{
			BranchInfos:       branchesSnapshot.Branches,
			BranchTypes:       []configdomain.BranchType{},
			BranchesAndTypes:  branchesAndTypes,
			ExcludeBranches:   excludeBranches,
			Lineage:           validatedConfig.NormalConfig.Lineage,
			MainBranch:        Some(validatedConfig.ValidatedConfigData.MainBranch),
			Order:             validatedConfig.NormalConfig.Order,
			Regexes:           []*regexp.Regexp{},
			ShowAllBranches:   true,
			UnknownBranchType: validatedConfig.NormalConfig.UnknownBranchType,
		}
		entriesAll := append(dialog.SwitchBranchEntries{noneEntry}, dialog.NewSwitchBranchEntries(entriesArgs)...)
		entriesArgs.ShowAllBranches = false
		entriesLocal := append(dialog.SwitchBranchEntries{noneEntry}, dialog.NewSwitchBranchEntries(entriesArgs)...)
		newParent, exit, err := dialog.SwitchBranch(dialog.SwitchBranchArgs{
			CurrentBranch:      None[gitdomain.LocalBranchName](),
			Cursor:             1, // select the "main branch" entry, below the "make perennial" entry
			DisplayBranchTypes: validatedConfig.NormalConfig.DisplayTypes,
			EntryData: dialog.EntryData{
				EntriesAll:      entriesAll,
				EntriesLocal:    entriesLocal,
				ShowAllBranches: false,
			},
			InputName:          fmt.Sprintf("parent-branch-for-%q", initialBranch),
			Inputs:             inputs,
			Title:              Some(fmt.Sprintf(messages.ParentBranchTitle, initialBranch)),
			UncommittedChanges: false,
		})
		if err != nil || exit {
			return data, configdomain.ProgramFlowExit, err
		}
		// store the new parent
		if err = validatedConfig.NormalConfig.SetParent(repo.Backend, initialBranch, newParent); err != nil {
			return data, configdomain.ProgramFlowContinue, err
		}
		ancestor = newParent
		hasAncestor = true
	}
	if args.beam.ShouldBeam() && hasAncestor {
		commitsInBranch, err := repo.Git.CommitsInFeatureBranch(repo.Backend, initialBranch, ancestor.BranchName())
		if err != nil {
			return data, configdomain.ProgramFlowExit, err
		}
		commitsToBeam, exit, err = dialog.CommitsToBeam(commitsInBranch, targetBranch, repo.Git, repo.Backend, inputs)
		if err != nil || exit {
			return data, configdomain.ProgramFlowExit, err
		}
	}
	if validatedConfig.NormalConfig.ShareNewBranches == configdomain.ShareNewBranchesPropose {
		args.propose = true
	}
	data = appendFeatureData{
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
	}
	return data, configdomain.ProgramFlowContinue, nil
}
