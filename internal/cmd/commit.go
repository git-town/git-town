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
	"github.com/git-town/git-town/v22/internal/cmd/sync"
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
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
	"github.com/spf13/cobra"
)

const (
	commitDesc = "Commit into another branch"
	commitHelp = `
Allows you to commit the currently staged changes
into another branch without needing to change branches.`
)

func commitCmd() *cobra.Command {
	addDownFlag, readDownFlag := flags.Down()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "commit",
		Args:    cobra.NoArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   commitDesc,
		Long:    cmdhelpers.Long(commitDesc, commitHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			down, errDown := readDownFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			message, errMessage := readMessageFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDown, errDryRun, errMessage, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeCommit(cliConfig, message, down)
		},
	}
	addDownFlag(&cmd)
	addDryRunFlag(&cmd)
	addMessageFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeCommit(cliConfig configdomain.PartialConfig, commitMessage Option[gitdomain.CommitMessage], down Option[configdomain.Down]) error {
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
	data, flow, err := determineCommitData(repo, commitMessage, down)
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
	runProgram := commitProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "commit",
		DryRun:                data.config.NormalConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[configdomain.EndConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		FinalUndoProgram:      program.Program{},
		BranchInfosLastRun:    data.branchInfosLastRun,
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

type commitData struct {
	branchInfosLastRun       Option[gitdomain.BranchInfos]
	branchToCommitInto       gitdomain.LocalBranchName
	branchesSnapshot         gitdomain.BranchesSnapshot
	branchesToSync           configdomain.BranchesToSync
	commitMessage            Option[gitdomain.CommitMessage]
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	inputs                   dialogcomponents.Inputs
	prefetchBranchesSnapshot gitdomain.BranchesSnapshot
	previousBranch           Option[gitdomain.LocalBranchName]
	remotes                  gitdomain.Remotes
	stashSize                gitdomain.StashSize
}

func determineCommitData(repo execute.OpenRepoResult, commitMessage Option[gitdomain.CommitMessage], down Option[configdomain.Down]) (commitData, configdomain.ProgramFlow, error) {
	var emptyCommitData commitData
	inputs := dialogcomponents.LoadInputs(os.Environ())
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	preFetchBranchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return emptyCommitData, configdomain.ProgramFlowExit, err
	}
	initialBranch, hasInitialBranch := preFetchBranchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return emptyCommitData, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return emptyCommitData, configdomain.ProgramFlowExit, err
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
		GiteaToken:           config.GiteaToken,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return emptyCommitData, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
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
	if err != nil {
		return emptyCommitData, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return emptyCommitData, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return emptyCommitData, configdomain.ProgramFlowExit, errors.New(messages.CommitDetachedHead)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return emptyCommitData, configdomain.ProgramFlowExit, err
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
		return emptyCommitData, configdomain.ProgramFlowExit, err
	}
	var branchToCommitIntoOpt Option[gitdomain.LocalBranchName]
	if down, hasDown := down.Get(); hasDown {
		if down {
			parent, hasParent := validatedConfig.NormalConfig.Lineage.Parent(initialBranch).Get()
			if !hasParent {
				return emptyCommitData, configdomain.ProgramFlowExit, fmt.Errorf(messages.CommitDownNoParent, initialBranch)
			}
			branchToCommitIntoOpt = Some(parent)
		}
	}
	branchToCommitInto, hasBranchToCommitInto := branchToCommitIntoOpt.Get()
	if !hasBranchToCommitInto {
		return emptyCommitData, configdomain.ProgramFlowExit, errors.New(messages.CommitNoBranchToCommitInto)
	}
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	branchNamesToSync := gitdomain.LocalBranchNames{initialBranch}
	allBranchNamesToSync := validatedConfig.NormalConfig.Lineage.BranchesAndAncestors(branchNamesToSync, validatedConfig.NormalConfig.Order)
	allBranchNamesToSync = allBranchNamesToSync.Remove(perennialAndMain...)
	allBranchNamesToSync = allBranchNamesToSync.Remove(branchToCommitInto)
	branchInfosToSync, _ := branchesSnapshot.Branches.Select(allBranchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return emptyCommitData, configdomain.ProgramFlowExit, err
	}
	return commitData{
		branchInfosLastRun:       branchInfosLastRun,
		branchToCommitInto:       branchToCommitInto,
		branchesSnapshot:         branchesSnapshot,
		branchesToSync:           branchesToSync,
		commitMessage:            commitMessage,
		config:                   validatedConfig,
		connector:                connector,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		inputs:                   inputs,
		prefetchBranchesSnapshot: preFetchBranchesSnapshot,
		previousBranch:           previousBranch,
		remotes:                  remotes,
		stashSize:                stashSize,
	}, configdomain.ProgramFlowContinue, nil
}

func commitProgram(data commitData) (runProgram program.Program) {
	prog := NewMutable(&program.Program{})
	// checkout the branch to commit into
	prog.Value.Add(
		&opcodes.Checkout{Branch: data.branchToCommitInto},
		&opcodes.Commit{
			AuthorOverride:                 Option[gitdomain.Author]{},
			FallbackToDefaultCommitMessage: false,
			Message:                        data.commitMessage,
		},
		&opcodes.Checkout{
			Branch: data.initialBranch,
		},
	)
	// git sync --detached --no-push

	sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
		BranchInfos:         data.branchesSnapshot.Branches,
		BranchInfosPrevious: data.branchInfosLastRun,
		BranchesToDelete:    NewMutable(&set.Set[gitdomain.LocalBranchName]{}),
		Config:              data.config,
		InitialBranch:       data.initialBranch,
		PrefetchBranchInfos: data.prefetchBranchesSnapshot.Branches,
		Program:             prog,
		Prune:               false,
		PushBranches:        false,
		Remotes:             data.remotes,
	})

	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             true,
		StashOpenChanges:         false,
		PreviousBranchCandidates: []Option[gitdomain.LocalBranchName]{data.previousBranch},
	})
	return prog.Immutable()
}
