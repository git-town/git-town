package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
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
	setParentDesc = "Set the parent branch for the current branch"
	setParentHelp = `
Consider this stack:

main
 \
  feature-1
   \
*   feature-B
 \
  feature-A

After running "git town set-parent"
and selecting "feature-A" in the dialog,
we end up with this stack:

main
 \
  feature-1
 \
  feature-A
   \
*   feature-B
`
)

func setParentCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "set-parent [branch]",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.MaximumNArgs(1),
		Short:   setParentDesc,
		Long:    cmdhelpers.Long(setParentDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  false,
				Verbose: verbose,
			}
			return executeSetParent(args, cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(args []string, cliConfig cliconfig.CliConfig) error {
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
	data, exit, err := determineSetParentData(repo, cliConfig)
	if err != nil || exit {
		return err
	}
	err = verifySetParentData(data)
	if err != nil {
		return err
	}
	var selectedBranch gitdomain.LocalBranchName
	newParentOpt := None[gitdomain.LocalBranchName]()
	switch len(args) {
	case 0:
		excludeBranches := append(
			gitdomain.LocalBranchNames{data.initialBranch},
			data.config.NormalConfig.Lineage.Children(data.initialBranch)...,
		)
		entries := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
			BranchInfos:       data.branchesSnapshot.Branches,
			BranchTypes:       []configdomain.BranchType{},
			BranchesAndTypes:  repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(data.branchesSnapshot.Branches.Names()),
			ExcludeBranches:   excludeBranches,
			Lineage:           repo.UnvalidatedConfig.NormalConfig.Lineage,
			MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
			Regexes:           []*regexp.Regexp{},
			ShowAllBranches:   false,
			UnknownBranchType: repo.UnvalidatedConfig.NormalConfig.UnknownBranchType,
		})
		noneEntry := dialog.SwitchBranchEntry{
			Branch:        messages.SetParentNoneOption,
			Indentation:   "",
			OtherWorktree: false,
			Type:          configdomain.BranchTypeFeatureBranch,
		}
		entries = append(dialog.SwitchBranchEntries{noneEntry}, entries...)
		selectedBranch, exit, err = dialog.SwitchBranch(dialog.SwitchBranchArgs{
			CurrentBranch:      None[gitdomain.LocalBranchName](),
			Cursor:             entries.IndexOf(data.defaultChoice),
			DisplayBranchTypes: true,
			Entries:            entries,
			InputName:          fmt.Sprintf("parent-branch-for-%q", data.initialBranch),
			Inputs:             data.inputs,
			Title:              Some(fmt.Sprintf(messages.ParentBranchTitle, data.initialBranch)),
			UncommittedChanges: false,
		})
		if err != nil || exit {
			return err
		}
		if selectedBranch != messages.SetParentNoneOption {
			newParentOpt = Some(selectedBranch)
		}
	case 1:
		selectedBranch = gitdomain.NewLocalBranchName(args[0])
		if !data.branchesSnapshot.Branches.HasLocalBranch(selectedBranch) {
			return fmt.Errorf(messages.BranchDoesntExist, selectedBranch)
		}
		newParentOpt = Some(selectedBranch)
	}
	runProgram, exit := setParentProgram(newParentOpt, data)
	if exit {
		return nil
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               "set-parent",
		DryRun:                false,
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

type setParentData struct {
	branchInfosLastRun Option[gitdomain.BranchInfos]
	branchesSnapshot   gitdomain.BranchesSnapshot
	config             config.ValidatedConfig
	connector          Option[forgedomain.Connector]
	defaultChoice      gitdomain.LocalBranchName
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	inputs             dialogcomponents.Inputs
	proposal           Option[forgedomain.Proposal]
	stashSize          gitdomain.StashSize
}

func determineSetParentData(repo execute.OpenRepoResult, cliConfig cliconfig.CliConfig) (data setParentData, exit dialogdomain.Exit, err error) {
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
		Verbose:               cliConfig.Verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{},
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
	mainBranch := validatedConfig.ValidatedConfigData.MainBranch
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	parentOpt := validatedConfig.NormalConfig.Lineage.Parent(initialBranch)
	existingParent, hasParent := parentOpt.Get()
	var defaultChoice gitdomain.LocalBranchName
	if hasParent {
		defaultChoice = existingParent
	} else {
		defaultChoice = mainBranch
	}
	proposalOpt := None[forgedomain.Proposal]()
	if !repo.IsOffline {
		proposalOpt = ship.FindProposal(connector, initialBranch, parentOpt)
	}
	return setParentData{
		branchInfosLastRun: branchInfosLastRun,
		branchesSnapshot:   branchesSnapshot,
		config:             validatedConfig,
		connector:          connector,
		defaultChoice:      defaultChoice,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		inputs:             inputs,
		proposal:           proposalOpt,
		stashSize:          stashSize,
	}, false, nil
}

func verifySetParentData(data setParentData) error {
	if data.config.IsMainOrPerennialBranch(data.initialBranch) {
		return fmt.Errorf(messages.SetParentNoFeatureBranch, data.initialBranch)
	}
	return nil
}

func setParentProgram(newParentOpt Option[gitdomain.LocalBranchName], data setParentData) (prog program.Program, exit dialogdomain.Exit) {
	proposal, hasProposal := data.proposal.Get()
	// update lineage
	newParent, hasNewParent := newParentOpt.Get()
	if !hasNewParent {
		prog.Add(&opcodes.BranchTypeOverrideSet{Branch: data.initialBranch, BranchType: configdomain.BranchTypePerennialBranch})
		prog.Add(&opcodes.LineageParentRemove{Branch: data.initialBranch})
	} else {
		prog.Add(&opcodes.LineageParentSet{Branch: data.initialBranch, Parent: newParent})
		connector, hasConnector := data.connector.Get()
		connectorCanUpdateProposalTarget := hasConnector && connector.UpdateProposalTargetFn().IsSome()
		if hasProposal && hasConnector && connectorCanUpdateProposalTarget {
			prog.Add(&opcodes.ProposalUpdateTarget{
				NewBranch: newParent,
				OldBranch: proposal.Data.Data().Target,
				Proposal:  proposal,
			})
		}
		// update commits
		switch data.config.NormalConfig.SyncFeatureStrategy {
		case configdomain.SyncFeatureStrategyMerge:
		case configdomain.SyncFeatureStrategyCompress, configdomain.SyncFeatureStrategyRebase:
			initialBranchInfo, hasInitialBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(data.initialBranch).Get()
			hasRemoteBranch := hasInitialBranchInfo && initialBranchInfo.HasTrackingBranch()
			if hasRemoteBranch {
				prog.Add(
					&opcodes.PullCurrentBranch{},
				)
			}
			parentOpt := data.config.NormalConfig.Lineage.Parent(data.initialBranch)
			if parent, hasParent := parentOpt.Get(); hasParent {
				prog.Add(
					&opcodes.RebaseOntoKeepDeleted{
						BranchToRebaseOnto: newParent.BranchName(),
						CommitsToRemove:    parent.Location(),
						Upstream:           None[gitdomain.LocalBranchName](),
					},
				)
			} else {
				prog.Add(
					&opcodes.RebaseBranch{
						Branch: newParent.BranchName(),
					},
				)
			}
			if hasRemoteBranch {
				prog.Add(
					&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
				)
			}
			// remove commits from descendents
			descendents := data.config.NormalConfig.Lineage.Descendants(data.initialBranch)
			for _, descendent := range descendents {
				prog.Add(
					&opcodes.CheckoutIfNeeded{
						Branch: descendent,
					},
				)
				descendentBranchInfo, hasDescendentBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(descendent).Get()
				if hasDescendentBranchInfo && descendentBranchInfo.HasTrackingBranch() {
					prog.Add(
						&opcodes.PullCurrentBranch{},
					)
				}
				prog.Add(
					&opcodes.RebaseOntoRemoveDeleted{
						BranchToRebaseOnto: data.initialBranch,
						CommitsToRemove:    descendent.BranchName(),
						Upstream:           parentOpt,
					},
				)
				if hasDescendentBranchInfo && descendentBranchInfo.HasTrackingBranch() {
					prog.Add(
						&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
					)
				}
			}
			prog.Add(
				&opcodes.CheckoutIfNeeded{
					Branch: data.initialBranch,
				},
			)
		}
	}
	return optimizer.Optimize(prog), false
}
