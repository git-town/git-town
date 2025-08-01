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
	"github.com/git-town/git-town/v21/internal/cmd/sync"
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
	"github.com/git-town/git-town/v21/pkg/set"
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
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
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
			bodyFile, errBodyFile := readBodyFileFlag(cmd)
			bodyText, errBodyText := readBodyFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			stack, errStack := readStackFlag(cmd)
			title, errTitle := readTitleFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errBodyFile, errBodyText, errDryRun, errAutoResolve, errStack, errTitle, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: autoResolve,
				DryRun:      dryRun,
				Verbose:     verbose,
			})
			return executePropose(proposeArgs{
				body:      bodyText,
				bodyFile:  bodyFile,
				cliConfig: cliConfig,
				stack:     stack,
				title:     title,
			})
		},
	}
	addBodyFlag(&cmd)
	addBodyFileFlag(&cmd)
	addDryRunFlag(&cmd)
	addAutoResolveFlag(&cmd)
	addStackFlag(&cmd)
	addTitleFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type proposeArgs struct {
	body      Option[gitdomain.ProposalBody]
	bodyFile  Option[gitdomain.ProposalBodyFile]
	cliConfig configdomain.PartialConfig
	stack     configdomain.FullStack
	title     Option[gitdomain.ProposalTitle]
}

func executePropose(args proposeArgs) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineProposeData(repo, args)
	if err != nil || exit {
		return err
	}
	runProgram := proposeProgram(repo, data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               proposeCmd,
		DryRun:                data.config.NormalConfig.DryRun,
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
	})
}

type proposeData struct {
	branchInfos         gitdomain.BranchInfos
	branchInfosLastRun  Option[gitdomain.BranchInfos]
	branchesSnapshot    gitdomain.BranchesSnapshot
	branchesToPropose   []branchToProposeData
	branchesToSync      configdomain.BranchesToSync
	config              config.ValidatedConfig
	connector           Option[forgedomain.Connector]
	hasOpenChanges      bool
	initialBranch       gitdomain.LocalBranchName
	inputs              dialogcomponents.Inputs
	nonExistingBranches gitdomain.LocalBranchNames // branches that are listed in the lineage information, but don't exist in the repo, neither locally nor remotely
	preFetchBranchInfos gitdomain.BranchInfos
	previousBranch      Option[gitdomain.LocalBranchName]
	proposalBody        Option[gitdomain.ProposalBody]
	proposalTitle       Option[gitdomain.ProposalTitle]
	remotes             gitdomain.Remotes
	stashSize           gitdomain.StashSize
}

type branchToProposeData struct {
	branchType          configdomain.BranchType
	existingProposalURL Option[string]
	name                gitdomain.LocalBranchName
	syncStatus          gitdomain.SyncStatus
}

func determineProposeData(repo execute.OpenRepoResult, args proposeArgs) (data proposeData, exit dialogdomain.Exit, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, false, err
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
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
		Connector:             connectorOpt,
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
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		ConfigSnapshot:     repo.ConfigSnapshot,
		Connector:          connectorOpt,
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
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToPropose gitdomain.LocalBranchNames
	var branchNamesToSync gitdomain.LocalBranchNames
	if args.stack {
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialAndMain)
		branchNamesToPropose = make(gitdomain.LocalBranchNames, len(branchNamesToSync))
		copy(branchNamesToPropose, branchNamesToSync)
	} else {
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchAndAncestorsWithoutRoot(initialBranch)
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
					existingProposalURL = Some(existingProposal.Data.Data().URL)
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
	bodyText, err := ship.ReadFile(args.body, args.bodyFile)
	return proposeData{
		branchInfos:         branchesSnapshot.Branches,
		branchInfosLastRun:  branchInfosLastRun,
		branchesSnapshot:    branchesSnapshot,
		branchesToPropose:   branchesToPropose,
		branchesToSync:      branchesToSync,
		config:              validatedConfig,
		connector:           connectorOpt,
		hasOpenChanges:      repoStatus.OpenChanges,
		initialBranch:       initialBranch,
		inputs:              inputs,
		nonExistingBranches: nonExistingBranches,
		preFetchBranchInfos: preFetchBranchSnapshot.Branches,
		previousBranch:      previousBranch,
		proposalBody:        bodyText,
		proposalTitle:       args.title,
		remotes:             remotes,
		stashSize:           stashSize,
	}, false, err
}

func proposeProgram(repo execute.OpenRepoResult, data proposeData) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, repo.FinalMessages, repo.Backend)
	branchesToDelete := set.New[gitdomain.LocalBranchName]()
	sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
		BranchInfos:         data.branchInfos,
		BranchInfosLastRun:  data.branchInfosLastRun,
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
		switch branchToPropose.branchType {
		case configdomain.BranchTypePrototypeBranch:
			prog.Value.Add(&opcodes.BranchTypeOverrideRemove{Branch: branchToPropose.name})
			repo.FinalMessages.Add(fmt.Sprintf(messages.PrototypeRemoved, branchToPropose.name))
		case configdomain.BranchTypeParkedBranch:
			prog.Value.Add(&opcodes.BranchTypeOverrideRemove{Branch: branchToPropose.name})
			repo.FinalMessages.Add(fmt.Sprintf(messages.ParkedRemoved, branchToPropose.name))
		case configdomain.BranchTypeFeatureBranch:
		case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypePerennialBranch:
			continue
		}
		prog.Value.Add(&opcodes.PushCurrentBranchIfLocal{
			CurrentBranch: branchToPropose.name,
		})
		previousBranchCandidates := []Option[gitdomain.LocalBranchName]{data.previousBranch}
		cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
			DryRun:                   data.config.NormalConfig.DryRun,
			InitialStashSize:         data.stashSize,
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
