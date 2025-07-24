package cmd

import (
	"cmp"
	"errors"
	"fmt"
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
			message, err1 := readMessageFlag(cmd)
			dryRun, err2 := readDryRunFlag(cmd)
			commitHook, err3 := readNoVerifyFlag(cmd)
			stack, err4 := readStackFlag(cmd)
			verbose, err5 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2, err3, err4, err5); err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  dryRun,
				Verbose: verbose,
			}
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

func executeCompress(cliConfig cliconfig.CliConfig, message Option[gitdomain.CommitMessage], commitHook configdomain.CommitHook, compressEntireStack configdomain.FullStack) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineCompressBranchesData(repo, cliConfig, message, compressEntireStack)
	if err != nil || exit {
		return err
	}
	err = validateCompressBranchesData(data, repo)
	if err != nil {
		return err
	}
	runProgram := compressProgram(data, commitHook)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               compressCommand,
		DryRun:                cliConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
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
		Connector:               None[forgedomain.Connector](),
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
		Verbose:                 cliConfig.Verbose,
	})
}

type compressBranchesData struct {
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
	hasTracking      bool
	name             gitdomain.LocalBranchName
	newCommitMessage gitdomain.CommitMessage // the commit message to use for the compressed commit in this branch
	parentBranch     gitdomain.LocalBranchName
}

func determineCompressBranchesData(repo execute.OpenRepoResult, cliConfig cliconfig.CliConfig, message Option[gitdomain.CommitMessage], compressEntireStack configdomain.FullStack) (data compressBranchesData, exit dialogdomain.Exit, err error) {
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
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
		Detached:              true,
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
		Verbose:               cliConfig.Verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
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
		return data, exit, err
	}
	perennialBranches := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToCompress gitdomain.LocalBranchNames
	if compressEntireStack {
		branchNamesToCompress = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialBranches)
	} else {
		branchNamesToCompress = gitdomain.LocalBranchNames{initialBranch}
	}
	branchesToCompress := []compressBranchData{}
	for _, branchNameToCompress := range branchNamesToCompress {
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToCompress).Get()
		if !hasBranchInfo {
			return data, exit, fmt.Errorf(messages.CompressNoBranchInfo, branchNameToCompress)
		}
		branchType := validatedConfig.BranchType(branchNameToCompress)
		if err := validateCanCompressBranchType(branchNameToCompress, branchType); err != nil {
			if compressEntireStack {
				continue
			}
			return data, exit, err
		}
		if err := validateBranchIsSynced(branchNameToCompress, branchInfo.SyncStatus); err != nil {
			return data, exit, err
		}
		parent := validatedConfig.NormalConfig.Lineage.Parent(branchNameToCompress)
		commits, err := repo.Git.CommitsInBranch(repo.Backend, branchNameToCompress, parent)
		if err != nil {
			return data, exit, err
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
				return data, false, err
			}
		}
		parentBranch, hasParent := parent.Get()
		if !hasParent {
			return data, exit, fmt.Errorf(messages.CompressBranchNoParent, branchNameToCompress)
		}
		hasRemoteBranch, _, _ := branchInfo.HasRemoteBranch()
		branchesToCompress = append(branchesToCompress, compressBranchData{
			branchType:       branchType,
			commitCount:      commitCount,
			hasTracking:      hasRemoteBranch,
			name:             branchNameToCompress,
			newCommitMessage: newCommitMessage,
			parentBranch:     parentBranch,
		})
	}
	if len(branchesToCompress) == 0 {
		return data, exit, fmt.Errorf(messages.CompressNoCommits, branchNamesToCompress[0])
	}
	return compressBranchesData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		branchesToCompress: branchesToCompress,
		config:             validatedConfig,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		inputs:             inputs,
		previousBranch:     previousBranch,
		stashSize:          stashSize,
	}, false, nil
}

func compressProgram(data compressBranchesData, commitHook configdomain.CommitHook) program.Program {
	prog := NewMutable(&program.Program{})
	for _, branchToCompress := range data.branchesToCompress {
		compressBranchProgram(prog, branchToCompress, data.config.NormalConfig.Offline, data.initialBranch, commitHook)
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

func compressBranchProgram(prog Mutable[program.Program], data compressBranchData, offline configdomain.Offline, initialBranch gitdomain.LocalBranchName, commitHook configdomain.CommitHook) {
	if !shouldCompressBranch(data.name, data.branchType, initialBranch) {
		return
	}
	prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.name})
	prog.Value.Add(&opcodes.BranchCurrentReset{Base: data.parentBranch.BranchName()})
	prog.Value.Add(&opcodes.CommitWithMessage{
		AuthorOverride: None[gitdomain.Author](),
		CommitHook:     commitHook,
		Message:        data.newCommitMessage,
	})
	if data.hasTracking && offline.IsOnline() {
		prog.Value.Add(&opcodes.PushCurrentBranchForceIfNeeded{CurrentBranch: data.name, ForceIfIncludes: true})
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

func validateCompressBranchesData(data compressBranchesData, repo execute.OpenRepoResult) error {
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
