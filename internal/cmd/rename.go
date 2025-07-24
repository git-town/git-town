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
	"github.com/git-town/git-town/v21/internal/cmd/ship"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
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
	renameDesc = "Rename a branch and its tracking branch"
	renameHelp = `
The branch to rename must be fully synced.

Renaming perennial branches requires the --force flag.
`
)

func renameCommand() *cobra.Command {
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addForceFlag, readForceFlag := flags.Force("force rename of perennial branch")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "rename [<old_branch_name>] <new_branch_name>",
		Args:  cobra.RangeArgs(1, 2),
		Short: renameDesc,
		Long:  cmdhelpers.Long(renameDesc, renameHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			dryRun, err1 := readDryRunFlag(cmd)
			force, err2 := readForceFlag(cmd)
			verbose, err3 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2, err3); err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  dryRun,
				Verbose: verbose,
			}
			return executeRename(args, cliConfig, force)
		},
	}
	addDryRunFlag(&cmd)
	addForceFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeRename(args []string, cliConfig cliconfig.CliConfig, force configdomain.Force) error {
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
	data, exit, err := determineRenameData(args, cliConfig, force, repo)
	if err != nil || exit {
		return err
	}
	runProgram := renameProgram(repo, data, repo.FinalMessages)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               "rename",
		DryRun:                cliConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
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
		Verbose:                 cliConfig.Verbose,
	})
}

type renameData struct {
	branchInfosLastRun       Option[gitdomain.BranchInfos]
	branchesSnapshot         gitdomain.BranchesSnapshot
	config                   config.ValidatedConfig
	connector                Option[forgedomain.Connector]
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	inputs                   dialogcomponents.Inputs
	newBranch                gitdomain.LocalBranchName
	nonExistingBranches      gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	oldBranch                gitdomain.BranchInfo
	previousBranch           Option[gitdomain.LocalBranchName]
	proposal                 Option[forgedomain.Proposal]
	proposalsOfChildBranches []forgedomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineRenameData(args []string, cliConfig cliconfig.CliConfig, force configdomain.Force, repo execute.OpenRepoResult) (data renameData, exit dialogdomain.Exit, err error) {
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
	var oldBranchName gitdomain.LocalBranchName
	var newBranchName gitdomain.LocalBranchName
	if len(args) == 1 {
		oldBranchName = initialBranch
		newBranchName = gitdomain.NewLocalBranchName(args[0])
	} else {
		oldBranchName = gitdomain.NewLocalBranchName(args[0])
		newBranchName = gitdomain.NewLocalBranchName(args[1])
	}
	oldBranch, hasOldBranch := branchesSnapshot.Branches.FindByLocalName(oldBranchName).Get()
	if !hasOldBranch {
		return data, false, fmt.Errorf(messages.BranchDoesntExist, oldBranchName)
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
		BranchesToValidate: gitdomain.LocalBranchNames{oldBranchName},
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
	if validatedConfig.ValidatedConfigData.IsMainBranch(oldBranchName) {
		return data, false, errors.New(messages.RenameMainBranch)
	}
	if !force {
		if validatedConfig.BranchType(oldBranchName) == configdomain.BranchTypePerennialBranch {
			return data, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return data, false, errors.New(messages.RenameToSameName)
	}
	if oldBranch.SyncStatus != gitdomain.SyncStatusUpToDate && oldBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return data, false, fmt.Errorf(messages.BranchNotInSyncWithParent, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName, repo.UnvalidatedConfig.NormalConfig.DevRemote) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	parentOpt := validatedConfig.NormalConfig.Lineage.Parent(initialBranch)
	lineageBranches := validatedConfig.NormalConfig.Lineage.BranchNames()
	_, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, lineageBranches...)
	proposalOpt := None[forgedomain.Proposal]()
	if !repo.IsOffline {
		proposalOpt = ship.FindProposal(connector, initialBranch, parentOpt)
	}
	proposalsOfChildBranches := ship.LoadProposalsOfChildBranches(ship.LoadProposalsOfChildBranchesArgs{
		ConnectorOpt:               connector,
		Lineage:                    validatedConfig.NormalConfig.Lineage,
		Offline:                    false,
		OldBranch:                  oldBranchName,
		OldBranchHasTrackingBranch: oldBranch.HasTrackingBranch(),
	})
	return renameData{
		branchInfosLastRun:       branchInfosLastRun,
		branchesSnapshot:         branchesSnapshot,
		config:                   validatedConfig,
		connector:                connector,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		inputs:                   inputs,
		newBranch:                newBranchName,
		nonExistingBranches:      nonExistingBranches,
		oldBranch:                *oldBranch,
		previousBranch:           previousBranch,
		proposal:                 proposalOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, err
}

func renameProgram(repo execute.OpenRepoResult, data renameData, finalMessages stringslice.Collector) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchesSnapshot.Branches, data.nonExistingBranches, finalMessages, repo.Backend)
	oldLocalBranch, hasOldLocalBranch := data.oldBranch.LocalName.Get()
	if !hasOldLocalBranch {
		return prog.Immutable()
	}
	prog.Value.Add(&opcodes.BranchLocalRename{OldName: oldLocalBranch, NewName: data.newBranch})
	if data.initialBranch == oldLocalBranch {
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: data.newBranch})
	}
	if !data.config.NormalConfig.DryRun {
		if override, hasBranchTypeOverride := data.config.NormalConfig.BranchTypeOverrides[oldLocalBranch]; hasBranchTypeOverride {
			prog.Value.Add(
				&opcodes.ConfigSet{
					Key:   configdomain.NewBranchTypeOverrideKeyForBranch(data.newBranch).Key,
					Scope: configdomain.ConfigScopeLocal,
					Value: override.String(),
				},
				&opcodes.ConfigRemove{
					Key:   configdomain.NewBranchTypeOverrideKeyForBranch(oldLocalBranch).Key,
					Scope: configdomain.ConfigScopeLocal,
				},
			)
		}
		if parentBranch, hasParent := data.config.NormalConfig.Lineage.Parent(oldLocalBranch).Get(); hasParent {
			prog.Value.Add(&opcodes.LineageParentSet{Branch: data.newBranch, Parent: parentBranch})
		}
		prog.Value.Add(&opcodes.LineageParentRemove{Branch: oldLocalBranch})
	}
	for _, child := range data.config.NormalConfig.Lineage.Children(oldLocalBranch) {
		prog.Value.Add(&opcodes.LineageParentSet{Branch: child, Parent: data.newBranch})
	}
	if oldTrackingBranch, hasOldTrackingBranch := data.oldBranch.RemoteName.Get(); hasOldTrackingBranch {
		if data.oldBranch.HasTrackingBranch() && data.config.NormalConfig.Offline.IsOnline() {
			prog.Value.Add(&opcodes.BranchTrackingCreate{Branch: data.newBranch})
			updateChildBranchProposalsToBranch(prog.Value, data.proposalsOfChildBranches, data.newBranch)
			proposal, hasProposal := data.proposal.Get()
			connector, hasConnector := data.connector.Get()
			connectorCanUpdateProposalSource := hasConnector && connector.UpdateProposalSourceFn().IsSome()
			if hasProposal && hasConnector && connectorCanUpdateProposalSource {
				prog.Value.Add(&opcodes.ProposalUpdateSource{
					NewBranch: data.newBranch,
					OldBranch: data.oldBranch.LocalBranchName(),
					Proposal:  proposal,
				})
			}
			prog.Value.Add(&opcodes.BranchTrackingDelete{Branch: oldTrackingBranch})
		}
	}
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.newBranch), data.previousBranch}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.config.NormalConfig.DryRun,
		InitialStashSize:         data.stashSize,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return optimizer.Optimize(prog.Immutable())
}

func updateChildBranchProposalsToBranch(prog *program.Program, proposals []forgedomain.Proposal, target gitdomain.LocalBranchName) {
	for _, childProposal := range proposals {
		prog.Add(&opcodes.ProposalUpdateTarget{
			NewBranch: target,
			OldBranch: childProposal.Data.Data().Target,
			Proposal:  childProposal,
		})
	}
}
