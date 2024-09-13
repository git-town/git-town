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
	"github.com/git-town/git-town/v16/internal/cmd/ship"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/undo/undoconfig"
	"github.com/git-town/git-town/v16/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v16/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
	"github.com/git-town/git-town/v16/internal/vm/program"
	"github.com/git-town/git-town/v16/internal/vm/runstate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const renameBranchDesc = "Rename a branch both locally and remotely"

const renameBranchHelp = `
Renames the given branch in the local and origin repository. Aborts if the new branch name already exists or the tracking branch is out of sync.

- creates a branch with the new name
- deletes the old branch

When there is an origin repository:
- syncs the repository

When there is a tracking branch:
- pushes the new branch to the origin repository
- deletes the old branch from the origin repository

When run on a perennial branch:
- confirm with the "--force"/"-f" option
- registers the new perennial branch name in the local Git Town configuration`

func renameBranchCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addForceFlag, readForceFlag := flags.Force("force rename of perennial branch")
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:   "rename-branch [<old_branch_name>] <new_branch_name>",
		Args:  cobra.RangeArgs(1, 2),
		Short: renameBranchDesc,
		Long:  cmdhelpers.Long(renameBranchDesc, renameBranchHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeRenameBranch(args, readDryRunFlag(cmd), readForceFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addForceFlag(&cmd)
	return &cmd
}

func executeRenameBranch(args []string, dryRun configdomain.DryRun, force configdomain.Force, verbose configdomain.Verbose) error {
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
	data, exit, err := determineRenameBranchData(args, force, repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	runProgram := renameBranchProgram(data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               "rename-branch",
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  data.config,
		Connector:               data.connector,
		DialogTestInputs:        data.dialogTestInputs,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranch:           data.initialBranch,
		InitialBranchesSnapshot: data.branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        data.stashSize,
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type renameBranchData struct {
	branchesSnapshot         gitdomain.BranchesSnapshot
	config                   config.ValidatedConfig
	connector                Option[hostingdomain.Connector]
	dialogTestInputs         components.TestInputs
	dryRun                   configdomain.DryRun
	hasOpenChanges           bool
	initialBranch            gitdomain.LocalBranchName
	newBranch                gitdomain.LocalBranchName
	oldBranch                gitdomain.BranchInfo
	previousBranch           Option[gitdomain.LocalBranchName]
	proposal                 Option[hostingdomain.Proposal]
	proposalsOfChildBranches []hostingdomain.Proposal
	stashSize                gitdomain.StashSize
}

func determineRenameBranchData(args []string, force configdomain.Force, repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose) (data renameBranchData, exit bool, err error) {
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
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
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{oldBranchName},
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	if validatedConfig.Config.IsMainBranch(oldBranchName) {
		return data, false, errors.New(messages.RenameMainBranch)
	}
	if force.IsFalse() {
		if validatedConfig.Config.IsPerennialBranch(oldBranchName) {
			return data, false, fmt.Errorf(messages.RenamePerennialBranchWarning, oldBranchName)
		}
	}
	if oldBranchName == newBranchName {
		return data, false, errors.New(messages.RenameToSameName)
	}
	if oldBranch.SyncStatus != gitdomain.SyncStatusUpToDate && oldBranch.SyncStatus != gitdomain.SyncStatusLocalOnly {
		return data, false, fmt.Errorf(messages.RenameBranchNotInSync, oldBranchName)
	}
	if branchesSnapshot.Branches.HasLocalBranch(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsLocally, newBranchName)
	}
	if branchesSnapshot.Branches.HasMatchingTrackingBranchFor(newBranchName) {
		return data, false, fmt.Errorf(messages.BranchAlreadyExistsRemotely, newBranchName)
	}
	var connectorOpt Option[hostingdomain.Connector]
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connectorOpt, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *validatedConfig.Config.UnvalidatedConfig,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			RemoteURL:       originURL,
		})
		if err != nil {
			return data, false, err
		}
	}
	existingParent, hasParent := validatedConfig.Config.Lineage.Parent(initialBranch).Get()
	proposalOpt := None[hostingdomain.Proposal]()
	connector, hasConnector := connectorOpt.Get()
	if hasConnector && connector.CanMakeAPICalls() && hasParent {
		proposalOpt, err = connector.FindProposal(initialBranch, existingParent)
		if err != nil {
			proposalOpt = None[hostingdomain.Proposal]()
		}
	}
	var proposalsOfChildBranches []hostingdomain.Proposal
	childBranches := validatedConfig.Config.Lineage.Children(oldBranchName)
	if hasConnector && connector.CanMakeAPICalls() {
		if !repo.IsOffline.IsTrue() && oldBranch.HasTrackingBranch() {
			for _, childBranch := range childBranches {
				childProposalOpt, err := connector.FindProposal(childBranch, oldBranchName)
				if err != nil {
					return data, false, fmt.Errorf(messages.ProposalNotFoundForBranch, oldBranchName, err)
				}
				childProposal, hasChildProposal := childProposalOpt.Get()
				if hasChildProposal {
					proposalsOfChildBranches = append(proposalsOfChildBranches, childProposal)
				}
			}
		}
	}
	return renameBranchData{
		branchesSnapshot:         branchesSnapshot,
		config:                   validatedConfig,
		connector:                connectorOpt,
		dialogTestInputs:         dialogTestInputs,
		dryRun:                   dryRun,
		hasOpenChanges:           repoStatus.OpenChanges,
		initialBranch:            initialBranch,
		newBranch:                newBranchName,
		oldBranch:                *oldBranch,
		previousBranch:           previousBranch,
		proposal:                 proposalOpt,
		proposalsOfChildBranches: proposalsOfChildBranches,
		stashSize:                stashSize,
	}, false, err
}

func renameBranchProgram(data renameBranchData) program.Program {
	result := NewMutable(&program.Program{})
	oldLocalBranch, hasOldLocalBranch := data.oldBranch.LocalName.Get()
	if !hasOldLocalBranch {
		return result.Get()
	}
	result.Value.Add(&opcodes.CreateBranch{Branch: data.newBranch, StartingPoint: oldLocalBranch.Location()})
	if data.initialBranch == oldLocalBranch {
		result.Value.Add(&opcodes.Checkout{Branch: data.newBranch})
	}
	if !data.dryRun {
		if data.config.Config.IsPerennialBranch(data.initialBranch) {
			result.Value.Add(&opcodes.RemoveFromPerennialBranches{Branch: oldLocalBranch})
			result.Value.Add(&opcodes.AddToPerennialBranches{Branch: data.newBranch})
		} else {
			if slices.Contains(data.config.Config.PrototypeBranches, data.initialBranch) {
				result.Value.Add(&opcodes.RemoveFromPrototypeBranches{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.AddToPrototypeBranches{Branch: data.newBranch})
			}
			if slices.Contains(data.config.Config.ObservedBranches, data.initialBranch) {
				result.Value.Add(&opcodes.RemoveFromObservedBranches{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.AddToObservedBranches{Branch: data.newBranch})
			}
			if slices.Contains(data.config.Config.ContributionBranches, data.initialBranch) {
				result.Value.Add(&opcodes.RemoveFromContributionBranches{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.AddToContributionBranches{Branch: data.newBranch})
			}
			if slices.Contains(data.config.Config.ParkedBranches, data.initialBranch) {
				result.Value.Add(&opcodes.RemoveFromParkedBranches{Branch: oldLocalBranch})
				result.Value.Add(&opcodes.AddToParkedBranches{Branch: data.newBranch})
			}
			if parentBranch, hasParent := data.config.Config.Lineage.Parent(oldLocalBranch).Get(); hasParent {
				result.Value.Add(&opcodes.SetParent{Branch: data.newBranch, Parent: parentBranch})
			}
			result.Value.Add(&opcodes.DeleteParentBranch{Branch: oldLocalBranch})
		}
	}
	for _, child := range data.config.Config.Lineage.Children(oldLocalBranch) {
		result.Value.Add(&opcodes.SetParent{Branch: child, Parent: data.newBranch})
	}
	if oldTrackingBranch, hasOldTrackingBranch := data.oldBranch.RemoteName.Get(); hasOldTrackingBranch {
		if data.oldBranch.HasTrackingBranch() && data.config.Config.IsOnline() {
			result.Value.Add(&opcodes.CreateTrackingBranch{Branch: data.newBranch})
			proposal, hasProposal := data.proposal.Get()
			if data.connector.IsSome() && hasProposal {
				result.Value.Add(&opcodes.UpdateProposalHead{
					NewTarget:      data.newBranch,
					ProposalNumber: proposal.Number,
				})
			}
			ship.UpdateChildBranchProposals(result.Value, data.proposalsOfChildBranches, data.newBranch)
			result.Value.Add(&opcodes.DeleteTrackingBranch{Branch: oldTrackingBranch})
		}
	}
	result.Value.Add(&opcodes.DeleteLocalBranch{Branch: oldLocalBranch})
	previousBranchCandidates := []Option[gitdomain.LocalBranchName]{Some(data.newBranch), data.previousBranch}
	cmdhelpers.Wrap(result, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             false,
		StashOpenChanges:         false,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	return result.Get()
}
