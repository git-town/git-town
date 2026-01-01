package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/cmd/ship"
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
	"github.com/git-town/git-town/v22/internal/vm/optimizer"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
	addNoParentFlag, readNoParentFlag := flags.NoParent()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "set-parent [branch]",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.MaximumNArgs(1),
		Short:   setParentDesc,
		Long:    cmdhelpers.Long(setParentDesc, setParentHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			noParent, errNoParent := readNoParentFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if cmp.Or(errAutoResolve, errNoParent, errVerbose) != nil {
				return errVerbose
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       autoResolve,
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          Some(configdomain.Detached(true)),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeSetParent(args, cliConfig, noParent)
		},
	}
	addAutoResolveFlag(&cmd)
	addNoParentFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSetParent(args []string, cliConfig configdomain.PartialConfig, noParent configdomain.NoParent) error {
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
	data, flow, err := determineSetParentData(repo)
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
	err = verifySetParentData(data)
	if err != nil {
		return err
	}
	var selectedParent gitdomain.LocalBranchName
	var exit dialogdomain.Exit
	newParentOpt := None[gitdomain.LocalBranchName]()
	if !noParent {
		switch len(args) {
		case 0:
			// TODO: extract this logic into an "enterParent" function
			excludeBranches := append(
				gitdomain.LocalBranchNames{data.initialBranch},
				data.config.NormalConfig.Lineage.Children(data.initialBranch, data.config.NormalConfig.Order)...,
			)
			noneEntry := dialog.SwitchBranchEntry{
				Branch:        messages.SetParentNoneOption,
				Indentation:   "",
				OtherWorktree: false,
				Type:          configdomain.BranchTypeFeatureBranch,
			}
			args := dialog.NewSwitchBranchEntriesArgs{
				BranchInfos:       data.branchesSnapshot.Branches,
				BranchTypes:       []configdomain.BranchType{},
				BranchesAndTypes:  repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(data.branchesSnapshot.Branches.NamesAllBranches()),
				ExcludeBranches:   excludeBranches,
				Lineage:           repo.UnvalidatedConfig.NormalConfig.Lineage,
				MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
				Order:             data.config.NormalConfig.Order,
				Regexes:           []*regexp.Regexp{},
				ShowAllBranches:   false,
				UnknownBranchType: repo.UnvalidatedConfig.NormalConfig.UnknownBranchType,
			}
			entriesLocal := append(dialog.SwitchBranchEntries{noneEntry}, dialog.NewSwitchBranchEntries(args)...)
			args.ShowAllBranches = true
			entriesAll := append(dialog.SwitchBranchEntries{noneEntry}, dialog.NewSwitchBranchEntries(args)...)
			selectedParent, exit, err = dialog.SwitchBranch(dialog.SwitchBranchArgs{
				CurrentBranch:      None[gitdomain.LocalBranchName](),
				Cursor:             entriesLocal.IndexOf(data.defaultChoice),
				DisplayBranchTypes: data.config.NormalConfig.DisplayTypes,
				EntryData: dialog.EntryData{
					EntriesAll:      entriesAll,
					EntriesLocal:    entriesLocal,
					ShowAllBranches: false,
				},
				InputName:          fmt.Sprintf("parent-branch-for-%q", data.initialBranch),
				Inputs:             data.inputs,
				Title:              Some(fmt.Sprintf(messages.ParentBranchTitle, data.initialBranch)),
				UncommittedChanges: false,
			})
			if err != nil || exit {
				return err
			}
			if selectedParent != messages.SetParentNoneOption {
				newParentOpt = Some(selectedParent)
			}
		case 1:
			selectedParent = gitdomain.NewLocalBranchName(args[0])
			if !data.branchesSnapshot.Branches.HasLocalBranch(selectedParent) {
				return fmt.Errorf(messages.BranchDoesntExist, selectedParent)
			}
			newParentOpt = Some(selectedParent)
		}
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

func determineSetParentData(repo execute.OpenRepoResult) (data setParentData, flow configdomain.ProgramFlow, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
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
		GitLabConnectorType:  config.GitLabConnectorType,
		GitLabToken:          config.GitLabToken,
		GiteaToken:           config.GiteaToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
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
		return data, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return data, flow, nil
	}
	if branchesSnapshot.DetachedHead {
		return data, configdomain.ProgramFlowExit, errors.New(messages.SetParentRepoHasDetachedHead)
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
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
		return data, configdomain.ProgramFlowExit, err
	}
	mainBranch := validatedConfig.ValidatedConfigData.MainBranch
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
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
	}, configdomain.ProgramFlowContinue, nil
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
	oldParent, hasOldParent := data.config.NormalConfig.Lineage.Parent(data.initialBranch).Get()
	newParent, hasNewParent := newParentOpt.Get()
	if !hasNewParent {
		prog.Add(&opcodes.BranchTypeOverrideSet{Branch: data.initialBranch, BranchType: configdomain.BranchTypePerennialBranch})
		prog.Add(&opcodes.LineageParentRemove{Branch: data.initialBranch})
	} else {
		prog.Add(&opcodes.LineageParentSet{Branch: data.initialBranch, Parent: newParent})
		connector, hasConnector := data.connector.Get()
		_, canUpdateProposalTarget := connector.(forgedomain.ProposalTargetUpdater)
		if hasProposal && hasConnector && canUpdateProposalTarget {
			prog.Add(&opcodes.ProposalUpdateTarget{
				NewBranch: newParent,
				OldBranch: proposal.Data.Data().Target,
				Proposal:  proposal,
			})
		}
		// update commits
		switch data.config.NormalConfig.SyncFeatureStrategy {
		case configdomain.SyncFeatureStrategyMerge:
			// don't update commits when using the "merge" sync strategy
		case configdomain.SyncFeatureStrategyCompress, configdomain.SyncFeatureStrategyRebase:
			switch data.config.BranchType(data.initialBranch) {
			case
				configdomain.BranchTypeContributionBranch,
				configdomain.BranchTypeMainBranch,
				configdomain.BranchTypeObservedBranch,
				configdomain.BranchTypeParkedBranch,
				configdomain.BranchTypePerennialBranch:
				// don't update the commits of these branch types
			case
				configdomain.BranchTypePrototypeBranch,
				configdomain.BranchTypeFeatureBranch:
				initialBranchInfo, hasInitialBranchInfo := data.branchesSnapshot.Branches.FindByLocalName(data.initialBranch).Get()
				hasRemoteBranch := hasInitialBranchInfo && initialBranchInfo.HasTrackingBranch()
				if hasRemoteBranch {
					prog.Add(
						&opcodes.PullCurrentBranch{},
					)
				}
				// remove the old parent's changes from the moved branch
				if hasOldParent {
					prog.Add(
						&opcodes.RebaseOnto{
							BranchToRebaseOnto: newParent.BranchName(),
							CommitsToRemove:    oldParent.Location(),
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
			}
			// remove the old parent's changes from the descendents of the moved branch
			if hasOldParent {
				descendents := data.config.NormalConfig.Lineage.Descendants(data.initialBranch, data.config.NormalConfig.Order)
				for _, descendent := range descendents {
					switch data.config.BranchType(descendent) {
					case
						configdomain.BranchTypeContributionBranch,
						configdomain.BranchTypeMainBranch,
						configdomain.BranchTypeObservedBranch,
						configdomain.BranchTypeParkedBranch,
						configdomain.BranchTypePerennialBranch:
						// don't update the commits on thes branch types
					case
						configdomain.BranchTypePrototypeBranch,
						configdomain.BranchTypeFeatureBranch:
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
							&opcodes.RebaseOnto{
								BranchToRebaseOnto: data.initialBranch.BranchName(),
								CommitsToRemove:    oldParent.Location(),
							},
						)
						if hasDescendentBranchInfo && descendentBranchInfo.HasTrackingBranch() {
							prog.Add(
								&opcodes.PushCurrentBranchForce{ForceIfIncludes: true},
							)
						}
					}
				}
			}
			prog.Add(
				&opcodes.CheckoutIfNeeded{
					Branch: data.initialBranch,
				},
			)
		}
	}
	// update proposal lineages
	if data.config.NormalConfig.ProposalsShowLineage == forgedomain.ProposalsShowLineageCLI {
		parents := gitdomain.LocalBranchNames{}
		if hasOldParent {
			parents = append(parents, oldParent)
		}
		if hasNewParent {
			parents = append(parents, newParent)
		}
		sync.AddSyncProposalsProgram(sync.AddSyncProposalsProgramArgs{
			ChangedBranches: parents,
			Config:          data.config,
			Program:         NewMutable(&prog),
		})
	}
	return optimizer.Optimize(prog), false
}
