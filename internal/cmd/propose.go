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
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/git-town/git-town/v22/internal/state/runstate"
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
				AutoResolve:       autoResolve,
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          Some(configdomain.Detached(true)),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            dryRun,
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
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
		ConfigDir:               repo.ConfigDir,
		Connector:               data.connector,
		DryRun:                  data.config.NormalConfig.DryRun,
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
	branchType configdomain.BranchType
	name       gitdomain.LocalBranchName
	syncStatus gitdomain.SyncStatus
}

func determineProposeData(repo execute.OpenRepoResult, args proposeArgs) (proposeData, configdomain.ProgramFlow, error) {
	var emptyResult proposeData
	preFetchBranchSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	if preFetchBranchSnapshot.DetachedHead {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.ProposeDetached)
	}
	initialBranch, hasInitialBranch := preFetchBranchSnapshot.Active.Get()
	if !hasInitialBranch {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	config := repo.UnvalidatedConfig.NormalConfig
	connectorOpt, err := forge.NewConnector(forge.NewConnectorArgs{
		Backend:              repo.Backend,
		BitbucketAppPassword: config.BitbucketAppPassword,
		BitbucketUsername:    config.BitbucketUsername,
		Browser:              config.Browser,
		ConfigDir:            repo.ConfigDir,
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
		return emptyResult, configdomain.ProgramFlowExit, err
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
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return emptyResult, flow, nil
	}
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
		ConfigDir:          repo.ConfigDir,
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
		return emptyResult, configdomain.ProgramFlowExit, err
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
			return emptyResult, configdomain.ProgramFlowExit, err
		}
	}
	branchesToPropose := []branchToProposeData{}
	for _, branchNameToPropose := range branchNamesToPropose {
		branchType, has := branchesAndTypes[branchNameToPropose]
		if !has {
			continue
		}
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindByLocalName(branchNameToPropose).Get()
		if !hasBranchInfo {
			return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.BranchInfoNotFound, branchNameToPropose)
		}
		branchesToPropose = append(branchesToPropose, branchToProposeData{
			branchType: branchType,
			name:       branchNameToPropose,
			syncStatus: branchInfo.SyncStatus,
		})
	}
	branchInfosToSync, nonExistingBranches := branchesSnapshot.Branches.Select(branchNamesToSync...)
	branchesToSync, err := sync.BranchesToSync(branchInfosToSync, branchesSnapshot.Branches, repo, validatedConfig.ValidatedConfigData.MainBranch)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
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
	}, configdomain.ProgramFlowContinue, err
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
		if branchToPropose.syncStatus == gitdomain.SyncStatusDeletedAtRemote {
			repo.FinalMessages.Addf(messages.BranchDeletedAtRemote, branchToPropose.name)
			continue
		}
		switch branchToPropose.branchType {
		case configdomain.BranchTypePrototypeBranch:
			prog.Value.Add(&opcodes.BranchTypeOverrideRemove{Branch: branchToPropose.name})
			repo.FinalMessages.Addf(messages.PrototypeRemoved, branchToPropose.name)
		case configdomain.BranchTypeParkedBranch:
			prog.Value.Add(&opcodes.BranchTypeOverrideRemove{Branch: branchToPropose.name})
			repo.FinalMessages.Addf(messages.ParkedRemoved, branchToPropose.name)
		case configdomain.BranchTypeFeatureBranch:
		case configdomain.BranchTypeContributionBranch, configdomain.BranchTypeMainBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypePerennialBranch:
			continue
		}
		prog.Value.Add(&opcodes.BranchTrackingCreateIfNeeded{
			CurrentBranch: branchToPropose.name,
		})
		prog.Value.Add(&opcodes.CheckoutIfNeeded{Branch: branchToPropose.name})
		updateBreadcrumb := data.config.NormalConfig.ProposalBreadcrumb.Enabled()
		proposalBody := data.proposalBody
		if updateBreadcrumb {
			lineageSection := proposallineage.RenderSection(data.config.NormalConfig.Lineage, branchToPropose.name, data.config.NormalConfig.Order, data.config.NormalConfig.ProposalBreadcrumb, data.config.NormalConfig.ProposalBreadcrumbDirection, data.connector)
			if len(lineageSection) > 0 {
				proposalBody = Some(proposallineage.UpdateProposalBody(proposalBody.GetOrZero(), lineageSection))
			}
		}
		prog.Value.Add(&opcodes.ProposalCreate{
			Branch:        branchToPropose.name,
			MainBranch:    data.config.ValidatedConfigData.MainBranch,
			ProposalBody:  proposalBody,
			ProposalTitle: data.proposalTitle,
		})
		prog.Value.Add(&opcodes.ProgramEndOfBranch{})
	}
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
