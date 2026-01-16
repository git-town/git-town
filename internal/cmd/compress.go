package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/optimizer"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	compressCommand = "compress"
	compressDesc    = "Squash all commits on the current branch down to a single commit"
	compressHelp    = `
Compress is a more convenient way of running "git rebase --interactive"
and choosing to fixup all commits.
Branches must be in sync to compress them, run "git sync" as needed.

Provide the --stack switch to compress all branches in the stack.

The compressed commit uses the commit message of the first commit in the branch.
You can provide a custom commit message with the -m switch.

Assuming you have a feature branch with these commits:

$ git log --format='%s'
commit 1
commit 2
commit 3

Let's compress these three commits into a single commit:

$ git town compress

Now your branch has a single commit with the name of the first commit but
containing the changes of all three commits that existed on the branch before:

$ git log --format='%s'
commit 1
`
)

func compressCmd() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addMessageFlag, readMessageFlag := flags.CommitMessage("customize the commit message")
	addNoVerifyFlag, readNoVerifyFlag := flags.NoVerify()
	addStackFlag, readStackFlag := flags.Stack("Compress the entire stack")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   compressCommand,
		Args:  cobra.NoArgs,
		Short: compressDesc,
		Long:  cmdhelpers.Long(compressDesc, compressHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			commitHook, errCommitHook := readNoVerifyFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			message, errMessage := readMessageFlag(cmd)
			stack, errStack := readStackFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errMessage, errDryRun, errCommitHook, errStack, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          Some(configdomain.Detached(true)),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeCompress(cliConfig, message, commitHook, stack)
		},
	}
	addDryRunFlag(&cmd)
	addMessageFlag(&cmd)
	addNoVerifyFlag(&cmd)
	addStackFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeCompress(cliConfig configdomain.PartialConfig, message Option[gitdomain.CommitMessage], commitHook configdomain.CommitHook, compressEntireStack configdomain.FullStack) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineCompressData(repo, message, compressEntireStack)
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
	err = validateCompressData(data, repo)
	if err != nil {
		return err
	}
	runProgram := compressProgram(data, commitHook)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               compressCommand,
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		BranchInfosLastRun:    data.branchInfosLastRun,
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		ConfigDir:               repo.ConfigDir,
		Connector:               None[forgedomain.Connector](),
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
		RunState:                runState,
	})
}

type compressData struct {
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	branchesToCompress []compressBranchData
	config             config.ValidatedConfig
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	inputs             dialogcomponents.Inputs
	previousBranch     Option[gitdomain.LocalBranchName]
	stashSize          gitdomain.StashSize
}

type compressBranchData struct {
	branchType       configdomain.BranchType
	commitCount      int // number of commits in this branch
	name             gitdomain.LocalBranchName
	newCommitMessage gitdomain.CommitMessage // the commit message to use for the compressed commit in this branch
	parentBranch     gitdomain.LocalBranchName
	trackingBranch   Option[gitdomain.RemoteBranchName]
}

func determineCompressData(repo execute.OpenRepoResult, message Option[gitdomain.CommitMessage], compressEntireStack configdomain.FullStack) (compressData, configdomain.ProgramFlow, error) {
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	inputs := dialogcomponents.LoadInputs(os.Environ())
	var emptyResult compressData
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ConfigDir:            repo.ConfigDir,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
		Frontend:             repo.Frontend,
		GiteaToken:           config.GiteaToken,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 true,
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
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return emptyResult, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.CompressDetachedHead)
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		ConfigDir:          repo.ConfigDir,
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
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	perennialBranches := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToCompress gitdomain.LocalBranchNames
	if compressEntireStack {
		branchNamesToCompress = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialBranches, validatedConfig.NormalConfig.Order)
	} else {
		branchNamesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesToCompress := []compressBranchData{}
	for _, branchNameToCompress := range branchNamesToCompress {
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToCompress).Get()
		if !hasBranchInfo {
			return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.CompressNoBranchInfo, branchNameToCompress)
		}
		branchType := validatedConfig.BranchType(branchNameToCompress)
		if err := validateCanCompressBranchType(branchNameToCompress, branchType); err != nil {
			if compressEntireStack {
				continue
			}
			return emptyResult, configdomain.ProgramFlowExit, err
		}
		if err := validateBranchIsSynced(branchNameToCompress, branchInfo.SyncStatus); err != nil {
			return emptyResult, configdomain.ProgramFlowExit, err
		}
		parent := validatedConfig.NormalConfig.Lineage.Parent(branchNameToCompress)
		commits, err := repo.Git.CommitsInBranch(repo.Backend, branchNameToCompress, parent)
		if err != nil {
			return emptyResult, configdomain.ProgramFlowExit, err
		}
		commitCount := len(commits)
		if commitCount == 0 {
			continue
		}
		var newCommitMessage gitdomain.CommitMessage
		if messageContent, has := message.Get(); has {
			newCommitMessage = messageContent
		} else {
			newCommitMessage, err = repo.Git.CommitMessage(repo.Backend, commits[0].SHA)
			if err != nil {
				return emptyResult, configdomain.ProgramFlowExit, err
			}
		}
		parentBranch, hasParent := parent.Get()
		if !hasParent {
			return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.CompressBranchNoParent, branchNameToCompress)
		}
		branchesToCompress = append(branchesToCompress, compressBranchData{
			branchType:       branchType,
			commitCount:      commitCount,
			name:             branchNameToCompress,
			newCommitMessage: newCommitMessage,
			parentBranch:     parentBranch,
			trackingBranch:   branchInfo.RemoteName,
		})
	}
	if len(branchesToCompress) == 0 {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.CompressNoCommits, branchNamesToCompress[0])
	}
	return compressData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		branchesToCompress: branchesToCompress,
		config:             validatedConfig,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		inputs:             inputs,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, configdomain.ProgramFlowContinue, nil
}

func compressProgram(data compressData, commitHook configdomain.CommitHook) program.Program {
	prog := NewMutable(&program.Program{})
	for _, branchToCompress := range data.branchesToCompress {
		compressBranchProgram(compressBranchProgramArgs{
			commitHook:    commitHook,
			data:          branchToCompress,
			initialBranch: data.initialBranch,
			offline:       data.config.NormalConfig.Offline,
			prog:          prog,
		})
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.initialBranch})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

type compressBranchProgramArgs struct {
	commitHook    configdomain.CommitHook
	data          compressBranchData
	initialBranch gitdomain.LocalBranchName
	offline       configdomain.Offline
	prog          Mutable[program.Program]
}

func compressBranchProgram(args compressBranchProgramArgs) {
	if !shouldCompressBranch(args.data.name, args.data.branchType, args.initialBranch) {
		return
	}
	args.prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: args.data.name})
	args.prog.Value.Add(&opcodes.BranchCurrentReset{Base: args.data.parentBranch.BranchName()})
	args.prog.Value.Add(&opcodes.CommitWithMessage{
		AuthorOverride: None[gitdomain.Author](),
		CommitHook:     args.commitHook,
		Message:        args.data.newCommitMessage,
	})
	trackingBranch, hasTrackingBranch := args.data.trackingBranch.Get()
	if hasTrackingBranch && args.offline.IsOnline() {
		args.prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{
			CurrentBranch:   args.data.name,
			ForceIfIncludes: true,
			TrackingBranch:  trackingBranch,
		})
	}
}

func shouldCompressBranch(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType, initialBranchName gitdomain.LocalBranchName) bool {
	if branchName == initialBranchName {
		return true
	}
	return branchType != configdomain.BranchTypeParkedBranch
}

func validateCanCompressBranchType(branchName gitdomain.LocalBranchName, branchType configdomain.BranchType) error {
	switch branchType {
	case
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypePrototypeBranch:
		return nil
	case
		configdomain.BranchTypeMainBranch,
		configdomain.BranchTypePerennialBranch:
		return errors.New(messages.CompressIsPerennial)
	case configdomain.BranchTypeObservedBranch:
		return fmt.Errorf(messages.CompressObservedBranch, branchName)
	case configdomain.BranchTypeContributionBranch:
		return fmt.Errorf(messages.CompressContributionBranch, branchName)
	}
	return nil
}

func validateBranchIsSynced(branchName gitdomain.LocalBranchName, syncStatus gitdomain.SyncStatus) error {
	switch syncStatus {
	case
		gitdomain.SyncStatusUpToDate,
		gitdomain.SyncStatusLocalOnly:
		return nil
	case
		gitdomain.SyncStatusNotInSync,
		gitdomain.SyncStatusAhead,
		gitdomain.SyncStatusBehind,
		gitdomain.SyncStatusDeletedAtRemote,
		gitdomain.SyncStatusRemoteOnly,
		gitdomain.SyncStatusOtherWorktree:
		return fmt.Errorf(messages.CompressUnsynced, branchName)
	}
	panic("unhandled syncstatus: " + syncStatus.String())
}

func validateCompressData(data compressData, repo execute.OpenRepoResult) error {
	for _, branch := range data.branchesToCompress {
		if parent, hasParent := data.config.NormalConfig.Lineage.Parent(branch.name).Get(); hasParent {
			isInSyncWithParent, err := repo.Git.BranchInSyncWithParent(repo.Backend, branch.name, parent.BranchName())
			if err != nil {
				return err
			}
			if !isInSyncWithParent {
				return fmt.Errorf(messages.BranchNotInSyncWithParent, branch.name)
			}
		}
	}
	return nil
}
