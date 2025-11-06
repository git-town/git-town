package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
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
	"github.com/git-town/git-town/v22/internal/git"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v22/internal/validate"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/optimizer"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
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
GitHub, GitLab, Gitea, Bitbucket, and Forgejo.
When using self-hosted versions
this command needs to be configured with "git config %s <driver>"
where driver is "github", "gitlab", "gitea", or "bitbucket".
When using SSH identities,
this command needs to be configured with
"git config %s <hostname>"
where hostname matches what is in your ssh config file.`
)

func proposeCommand() *cobra.Command {
	addAutoResolveFlag, readAutoResolveFlag := flags.AutoResolve()
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
			autoResolve, errAutoResolve := readAutoResolveFlag(cmd)
			bodyFile, errBodyFile := readBodyFileFlag(cmd)
			bodyText, errBodyText := readBodyFlag(cmd)
			dryRun, errDryRun := readDryRunFlag(cmd)
			stack, errStack := readStackFlag(cmd)
			title, errTitle := readTitleFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errBodyFile, errBodyText, errDryRun, errAutoResolve, errStack, errTitle, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:  autoResolve,
				AutoSync:     None[configdomain.AutoSync](),
				Detached:     Some(configdomain.Detached(true)),
				DisplayTypes: None[configdomain.DisplayTypes](),
				DryRun:       dryRun,
				Order:        None[configdomain.Order](),
				PushBranches: None[configdomain.PushBranches](),
				Stash:        None[configdomain.Stash](),
				Verbose:      verbose,
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
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineProposeData(repo, args)
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
	runProgram := proposeProgram(repo, data)
	runState := runstate.RunState{
		BeginBranchesSnapshot: data.branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        data.stashSize,
		BranchInfosLastRun:    data.branchInfosLastRun,
		Command:               proposeCmd,
		DryRun:                data.config.NormalConfig.DryRun,
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

type determineProposalTitleArgs struct {
	Title   Option[gitdomain.ProposalTitle]
	Branch  gitdomain.LocalBranchName
	Parent  gitdomain.LocalBranchName
	Config  config.NormalConfig
	Git     git.Commands
	Inputs  dialogcomponents.Inputs
	Backend subshelldomain.Querier
}

// determineProposalTitle determines the proposal title based on CLI args, commits, and config.
func determineProposalTitle(args determineProposalTitleArgs) (Option[gitdomain.ProposalTitle], error) {
	if args.Title.IsSome() {
		return args.Title, nil
	}
	commits, err := args.Git.CommitsInFeatureBranch(args.Backend, args.Branch, args.Parent.BranchName())
	if err != nil {
		return None[gitdomain.ProposalTitle](), err
	}
	if len(commits) == 0 {
		// No commits - fallback to native
		return None[gitdomain.ProposalTitle](), nil
	}
	if len(commits) == 1 {
		parts := commits[0].Message.Parts()
		return Some(gitdomain.ProposalTitle(parts.Subject)), nil
	}
	proposeTitle := args.Config.ProposeTitle.GetOr(configdomain.ProposeTitleNative)
	switch proposeTitle {
	case configdomain.ProposeTitleFirst:
		parts := commits[0].Message.Parts()
		return Some(gitdomain.ProposalTitle(parts.Subject)), nil
	case configdomain.ProposeTitleSelect:
		commitsWithMessages, err := args.Git.CommitsInFeatureBranchWithMessages(args.Backend, args.Branch, args.Parent.BranchName())
		if err != nil {
			return None[gitdomain.ProposalTitle](), err
		}
		selectedTitle, exit, err := dialog.CommitTitle(commitsWithMessages, args.Inputs)
		if err != nil {
			return None[gitdomain.ProposalTitle](), err
		}
		if exit {
			return None[gitdomain.ProposalTitle](), errors.New("commit title selection aborted")
		}
		if selected, has := selectedTitle.Get(); has {
			parts := selected.Parts()
			return Some(gitdomain.ProposalTitle(parts.Subject)), nil
		}
		// User selected "none" - fallback to native
		return None[gitdomain.ProposalTitle](), nil
	case configdomain.ProposeTitleNative:
		// Fallback to native behavior
		return None[gitdomain.ProposalTitle](), nil
	default:
		// Unknown value - fallback to native
		return None[gitdomain.ProposalTitle](), nil
	}
}

func determineProposeData(repo execute.OpenRepoResult, args proposeArgs) (data proposeData, flow configdomain.ProgramFlow, err error) {
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	if preFetchBranchSnapshot.DetachedHead {
		return data, configdomain.ProgramFlowExit, errors.New(messages.ProposeDetached)
	}
	initialBranch, hasInitialBranch := preFetchBranchSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		ForgeType:            config.ForgeType,
		ForgejoToken:         config.ForgejoToken,
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
		return data, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, branchInfosLastRun, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connectorOpt,
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
		return data, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return data, flow, nil
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
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
		return data, configdomain.ProgramFlowExit, err
	}
	perennialAndMain := branchesAndTypes.BranchesOfTypes(configdomain.BranchTypePerennialBranch, configdomain.BranchTypeMainBranch)
	var branchNamesToPropose gitdomain.LocalBranchNames
	var branchNamesToSync gitdomain.LocalBranchNames
	if args.stack {
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchLineageWithoutRoot(initialBranch, perennialAndMain, validatedConfig.NormalConfig.Order)
		branchNamesToPropose = make(gitdomain.LocalBranchNames, len(branchNamesToSync))
		copy(branchNamesToPropose, branchNamesToSync)
	} else {
		branchNamesToSync = validatedConfig.NormalConfig.Lineage.BranchAndAncestorsWithoutRoot(initialBranch)
		branchNamesToPropose = gitdomain.LocalBranchNames{initialBranch}
		if err = validateBranchTypeToPropose(branchesAndTypes[initialBranch]); err != nil {
			return data, configdomain.ProgramFlowExit, err
		}
		if validatedConfig.NormalConfig.Lineage.Parent(initialBranch).IsNone() {
			return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.ProposalNoParent, initialBranch)
		}
	}
	connector, hasConnector := connectorOpt.Get()
	if !hasConnector {
		return data, configdomain.ProgramFlowExit, forgedomain.UnsupportedServiceError()
	}
	proposalFinder, canFindProposals := connector.(forgedomain.ProposalFinder)
	branchesToPropose := make([]branchToProposeData, len(branchNamesToPropose))
	for b, branchNameToPropose := range branchNamesToPropose {
		branchType, has := branchesAndTypes[branchNameToPropose]
		if !has {
			return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchTypeCannotDetermine, branchNameToPropose)
		}
		existingProposalURL := None[string]()
		if canFindProposals {
			if parent, hasParent := validatedConfig.NormalConfig.Lineage.Parent(branchNameToPropose).Get(); hasParent {
				existingProposalOpt, err := proposalFinder.FindProposal(branchNameToPropose, parent)
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
			return data, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchInfoNotFound, branchNameToPropose)
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
		return data, configdomain.ProgramFlowExit, err
	}
	bodyText, err := ship.ReadFile(args.body, args.bodyFile)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}

	// Determine proposal title (only for single branch proposals)
	var proposalTitle Option[gitdomain.ProposalTitle] = args.title
	if len(branchNamesToPropose) == 1 {
		branchToPropose := branchNamesToPropose[0]
		if parent, hasParent := validatedConfig.NormalConfig.Lineage.Parent(branchToPropose).Get(); hasParent {
			determinedTitle, err := determineProposalTitle(determineProposalTitleArgs{
				Title:   args.title,
				Branch:  branchToPropose,
				Parent:  parent,
				Config:  validatedConfig.NormalConfig,
				Git:     repo.Git,
				Inputs:  inputs,
				Backend: repo.Backend,
			})
			if err != nil {
				return data, configdomain.ProgramFlowExit, err
			}
			// Only override if we determined a title (not if it's None)
			if determinedTitle.IsSome() {
				proposalTitle = determinedTitle
			}
		}
	}

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
		proposalTitle:       proposalTitle,
		remotes:             remotes,
		stashSize:           stashSize,
	}, configdomain.ProgramFlowContinue, nil
}

func proposeProgram(repo execute.OpenRepoResult, data proposeData) program.Program {
	prog := NewMutable(&program.Program{})
	data.config.CleanupLineage(data.branchInfos, data.nonExistingBranches, repo.FinalMessages, repo.Backend, data.config.NormalConfig.Order)
	branchesToDelete := set.New[gitdomain.LocalBranchName]()
	sync.BranchesProgram(data.branchesToSync, sync.BranchProgramArgs{
		BranchInfos:         data.branchInfos,
		BranchInfosPrevious: data.branchInfosLastRun,
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
