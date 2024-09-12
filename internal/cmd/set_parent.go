package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
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

const setParentCmd = "set-parent"

const setParentDesc = "Prompt to set the parent branch for the current branch"

func setParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     setParentCmd,
		GroupID: "lineage",
		Args:    cobra.NoArgs,
		Short:   setParentDesc,
		Long:    cmdhelpers.Long(setParentDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeSetParent(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSetParentData(repo, verbose)
	if err != nil || exit {
		return err
	}
	err = verifySetParentData(data)
	if err != nil {
		return err
	}
	outcome, selectedBranch, err := dialog.Parent(dialog.ParentArgs{
		Branch:          data.initialBranch,
		DefaultChoice:   data.defaultChoice,
		DialogTestInput: data.dialogTestInputs.Next(),
		Lineage:         data.config.Config.Lineage,
		LocalBranches:   data.branchesSnapshot.Branches.LocalBranches().Names(),
		MainBranch:      data.mainBranch,
	})
	if err != nil {
		return err
	}
	runProgram, aborted := setParentProgram(outcome, selectedBranch, data.initialBranch)
	if aborted {
		return nil
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               setParentCmd,
		DryRun:                false,
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
		Connector:               None[hostingdomain.Connector](),
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

type setParentData struct {
	branchesSnapshot gitdomain.BranchesSnapshot
	config           config.ValidatedConfig
	connector        Option[hostingdomain.Connector]
	defaultChoice    gitdomain.LocalBranchName
	dialogTestInputs components.TestInputs
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	mainBranch       gitdomain.LocalBranchName
	proposal         Option[hostingdomain.Proposal]
	stashSize        gitdomain.StashSize
}

func determineSetParentData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data setParentData, exit bool, err error) {
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
		Fetch:                 false,
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
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{},
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
	mainBranch := validatedConfig.Config.MainBranch
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	existingParent, hasParent := validatedConfig.Config.Lineage.Parent(initialBranch).Get()
	var defaultChoice gitdomain.LocalBranchName
	if hasParent {
		defaultChoice = existingParent
	} else {
		defaultChoice = mainBranch
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
	proposalOpt := None[hostingdomain.Proposal]()
	connector, hasConnector := connectorOpt.Get()
	if hasConnector && connector.CanMakeAPICalls() && hasParent {
		proposalOpt, err = connector.FindProposal(initialBranch, existingParent)
		if err != nil {
			proposalOpt = None[hostingdomain.Proposal]()
		}
	}
	return setParentData{
		branchesSnapshot: branchesSnapshot,
		config:           validatedConfig,
		defaultChoice:    defaultChoice,
		dialogTestInputs: dialogTestInputs,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    initialBranch,
		mainBranch:       mainBranch,
		proposal:         proposalOpt,
		stashSize:        stashSize,
	}, false, nil
}

func verifySetParentData(data setParentData) error {
	if data.config.Config.IsMainOrPerennialBranch(data.initialBranch) {
		return fmt.Errorf(messages.SetParentNoFeatureBranch, data.initialBranch)
	}
	return nil
}

func setParentProgram(outcome dialog.ParentOutcome, selectedBranch, currentBranch gitdomain.LocalBranchName, data setParentData) (prog program.Program, aborted bool) {
	connector, hasConnector := data.connector.Get()
	proposal, hasProposal := data.proposal.Get()
	switch outcome {
	case dialog.ParentOutcomeAborted:
		return prog, true
	case dialog.ParentOutcomePerennialBranch:
		prog.Add(&opcodes.AddToPerennialBranches{Branch: currentBranch})
		prog.Add(&opcodes.DeleteParentBranch{Branch: currentBranch})
	case dialog.ParentOutcomeSelectedParent:
		prog.Add(&opcodes.SetParent{Branch: currentBranch, Parent: selectedBranch})
		if hasConnector && hasProposal {
			prog.Add(&opcodes.UpdateProposalTarget{})
		}
	}
	return prog, false
}
