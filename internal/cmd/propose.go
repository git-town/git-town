package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/flags"
	"github.com/git-town/git-town/v19/internal/cli/print"
	"github.com/git-town/git-town/v19/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v19/internal/cmd/sync"
	"github.com/git-town/git-town/v19/internal/config"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/execute"
	"github.com/git-town/git-town/v19/internal/forge"
	"github.com/git-town/git-town/v19/internal/forge/forgedomain"
	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/messages"
	"github.com/git-town/git-town/v19/internal/undo/undoconfig"
	"github.com/git-town/git-town/v19/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v19/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v19/internal/vm/opcodes"
	"github.com/git-town/git-town/v19/internal/vm/optimizer"
	"github.com/git-town/git-town/v19/internal/vm/program"
	"github.com/git-town/git-town/v19/internal/vm/runstate"
	. "github.com/git-town/git-town/v19/pkg/prelude"
	"github.com/git-town/git-town/v19/pkg/set"
	"github.com/spf13/cobra"
)

const (
	proposeCmd  = "propose"
	proposeDesc = "Create a proposal to merge a feature branch"
	proposeHelp = `
Syncs the current branch and opens a browser window
to the new proposal page of your repository.

The form is pre-populated for the current branch
so that the proposal only shows the changes
made against the immediate parent branch.

Supported only for repositories hosted on
GitHub, GitLab, Gitea, Bitbucket, and Codeberg.
When using self-hosted versions
this command needs to be configured with "git config %s <driver>"
where driver is "github", "gitlab", "gitea", or "bitbucket".
When using SSH identities,
this command needs to be configured with
"git config %s <hostname>"
where hostname matches what is in your ssh config file.`
)

func proposeCommand() *cobra.Command {
	addBodyFlag, readBodyFlag := flags.ProposalBody("b")
	addBodyFileFlag, readBodyFileFlag := flags.ProposalBodyFile()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addStackFlag, readStackFlag := flags.Stack("propose the entire stack")
	addTitleFlag, readTitleFlag := flags.ProposalTitle()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     proposeCmd,
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, configdomain.KeyForgeType, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			bodyFile, err := readBodyFileFlag(cmd)
			if err != nil {
				return err
			}
			bodyText, err := readBodyFlag(cmd)
			if err != nil {
				return err
			}
			stack, err := readStackFlag(cmd)
			if err != nil {
				return err
			}
			title, err := readTitleFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executePropose(dryRun, verbose, title, bodyText, bodyFile, stack)
		},
	}
	addBodyFlag(&cmd)
	addBodyFileFlag(&cmd)
	addDryRunFlag(&cmd)
	addStackFlag(&cmd)
	addTitleFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executePropose(dryRun configdomain.DryRun, verbose configdomain.Verbose, title gitdomain.ProposalTitle, body gitdomain.ProposalBody, bodyFile gitdomain.ProposalBodyFile, fullStack configdomain.FullStack) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineProposeData(repo, dryRun, fullStack, verbose, title, body, bodyFile)
	if err != nil || exit {
		return err
	}
	runProgram := proposeProgram(repo, data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		Command:               proposeCmd,
		DryRun:                dryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		RunProgram:            runProgram,
		TouchedBranches:       runProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
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

type proposeData struct {
	branchInfos         gitdomain.BranchInfos
	branchesSnapshot    gitdomain.BranchesSnapshot
	branchesToPropose   []branchToProposeData
	branchesToSync      configdomain.BranchesToSync
	config              config.ValidatedConfig
	connector           Option[forgedomain.Connector]
	dialogTestInputs    components.TestInputs
	dryRun              configdomain.DryRun
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos gitdomain.BranchInfos
	previousBranch      Option[gitdomain.LocalBranchName]
	proposalBody        gitdomain.ProposalBody
	proposalTitle       gitdomain.ProposalTitle
	remotes             gitdomain.Remotes
	stashSize           gitdomain.StashSize
}

type branchToProposeData struct {
	branchType          configdomain.BranchType
	existingProposalURL Option[string]
	name                gitdomain.LocalBranchName
	syncStatus          gitdomain.SyncStatus
}

func determineProposeData(repo execute.OpenRepoResult, dryRun configdomain.DryRun, fullStack configdomain.FullStack, verbose configdomain.Verbose, title gitdomain.ProposalTitle, body gitdomain.ProposalBody, bodyFile gitdomain.ProposalBodyFile) (data proposeData, exit bool, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
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
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, false, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connectorOpt, err := forge.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		Connector:          connectorOpt,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToPropose gitdomain.LocalBranchNames
	var branchNamesToSync gitdomain.LocalBranchNames
	if fullStack {
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialAndMain)
		branchNamesToPropose = make(gitdomain.LocalBranchNames, len(branchNamesToSync))
		copy(branchNamesToPropose, branchNamesToSync)
	} else {
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchAndAncestors(initialBranch)
		branchNamesToPropose = gitdomain.LocalBranchNames{initialBranch}
		if err = validateBranchTypeToPropose(branchesAndTypes[initialBranch]); err != nil {
			return data, false, err
		}
		if validatedConfig.NormalConfig.Lineage.Parent(initialBranch).IsNone() {
			return data, false, fmt.Errorf(messages.ProposalNoParent, initialBranch)
		}
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return data, false, forgedomain.UnsupportedServiceError()
	}
	findProposal, canFindProposals := connector.FindProposalFn().Get()
	branchesToPropose := make([]branchToProposeData, len(branchNamesToPropose))
	for b, branchNameToPropose := range branchNamesToPropose {
		branchType, has := branchesAndTypes[branchNameToPropose]
		if !has {
			return data, false, fmt.Errorf(messages.BranchTypeCannotDetermine, branchNameToPropose)
		}
		existingProposalURL := None[string]()
		if canFindProposals {
			if parent, hasParent := validatedConfig.NormalConfig.Lineage.Parent(branchNameToPropose).Get(); hasParent {
				existingProposalOpt, err := findProposal(branchNameToPropose, parent)
				if err != nil {
					print.Error(err)
				}
				if existingProposal, has := existingProposalOpt.Get(); has {
					existingProposalURL = Some(existingProposal.URL)
				}
			}
		}
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToPropose).Get()
		if !hasBranchInfo {
			return data, false, fmt.Errorf(messages.BranchInfoNotFound, branchNameToPropose)
		}
		branchesToPropose[b] = branchToProposeData{
			branchType:          branchType,
			existingProposalURL: existingProposalURL,
			name:                branchNameToPropose,
			syncStatus:          branchInfo.SyncStatus,
		}
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(repo.UnvalidatedConfig.NormalConfig.DevRemote, branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return data, false, err
	}
	var bodyText gitdomain.ProposalBody
	if len(body) > 0 {
		bodyText = body
	} else if len(bodyFile) > 0 {
		if bodyFile.ShouldReadStdin() {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return data, false, fmt.Errorf("cannot read STDIN: %w", err)
			}
			bodyText = gitdomain.ProposalBody(content)
		} else {
			fileData, err := os.ReadFile(bodyFile.String())
			if err != nil {
				return data, false, err
			}
			bodyText = gitdomain.ProposalBody(fileData)
		}
	}
	return proposeData{
		branchInfos:         branchesSnapshot.Branches,
		branchesSnapshot:    branchesSnapshot,
		branchesToPropose:   branchesToPropose,
		branchesToSync:      branchesToSync,
		config:              validatedConfig,
		connector:           connectorOpt,
		dialogTestInputs:    dialogTestInputs,
		dryRun:              dryRun,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		nonExistingBranches: nonExistingBranches,
		preFetchBranchInfos: preFetchBranchSnapshot.Branches,
		previousBranch:      previousBranch,
		proposalBody:        bodyText,
		proposalTitle:       title,
		remotes:             remotes,
		stashSize:           stashSize,
	}, false, err
}

func proposeProgram(repo execute.OpenRepoResult, data proposeData) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, repo.FinalMessages)
	branchesToDelete := set.New[gitdomain.LocalBranchName]()
	sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
		BranchInfos:         data.branchInfos,
		BranchesToDelete:    NewMutable(&branchesToDelete),
		Config:              data.config,
		InitialBranch:       data.initialBranch,
		PrefetchBranchInfos: data.preFetchBranchInfos,
		Remotes:             data.remotes,
		Program:             prog,
		Prune:               false,
		PushBranches:        true,
	})
	for _, branchToPropose := range data.branchesToPropose {
		if branchToPropose.branchType == configdomain.BranchTypePrototypeBranch {
			prog.Value.Add(&opcodes.BranchTypeOverrideRemove{Branch: branchToPropose.name})
			repo.FinalMessages.Add(fmt.Sprintf(messages.PrototypeRemoved, branchToPropose.name))
		}
		prog.Value.Add(&opcodes.PushCurrentBranchIfLocal{
			CurrentBranch: branchToPropose.name,
		})
		previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
		cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
			DryRun:                   data.dryRun,
			RunInGitRoot:             true,
			StashOpenChanges:         data.hasOpenChanges,
			PreviousBranchCandidates: previousBranchCandidates,
		})
		if branchToPropose.syncStatus == gitdomain.SyncStatusDeletedAtRemote {
			repo.FinalMessages.Add(fmt.Sprintf(messages.BranchDeletedAtRemote, branchToPropose.name))
			return prog.Immutable()
		}
		if existingProposalURL, hasExistingProposal := branchToPropose.existingProposalURL.Get(); hasExistingProposal {
			prog.Value.Add(
				&opcodes.BrowserOpen{
					URL: existingProposalURL,
				},
			)
		} else {
			prog.Value.Add(&opcodes.ProposalCreate{
				Branch:        branchToPropose.name,
				MainBranch:    data.config.ValidatedConfigData.MainBranch,
				ProposalBody:  data.proposalBody,
				ProposalTitle: data.proposalTitle,
			})
		}
	}
	return optimizer.Optimize(prog.Immutable())
}

func validateBranchTypeToPropose(branchType configdomain.BranchType) error {
	switch branchType {
	case
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypeParkedBranch,
		configdomain.BranchTypePrototypeBranch:
		return nil
	case configdomain.BranchTypeMainBranch:
		return errors.New(messages.MainBranchCannotPropose)
	case configdomain.BranchTypeContributionBranch:
		return errors.New(messages.ContributionBranchCannotPropose)
	case configdomain.BranchTypeObservedBranch:
		return errors.New(messages.ObservedBranchCannotPropose)
	case configdomain.BranchTypePerennialBranch:
		return errors.New(messages.PerennialBranchCannotPropose)
	}
	return nil
}
