package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/flags"
	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/execute"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/git-town/git-town/v15/internal/sync"
	"github.com/git-town/git-town/v15/internal/undo/undoconfig"
	"github.com/git-town/git-town/v15/internal/validate"
	fullInterpreter "github.com/git-town/git-town/v15/internal/vm/interpreter/full"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/git-town/git-town/v15/internal/vm/runstate"
	"github.com/spf13/cobra"
)

const proposeCmd = "propose"

const proposeDesc = "Create a proposal to merge a feature branch"

const proposeHelp = `
Syncs the current branch and opens a browser window to the new proposal page of your repository.

The form is pre-populated for the current branch so that the proposal only shows the changes made against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, Gitea and Bitbucket. When using self-hosted versions this command needs to be configured with "git config %s <driver>" where driver is "github", "gitlab", "gitea", or "bitbucket". When using SSH identities, this command needs to be configured with "git config %s <hostname>" where hostname matches what is in your ssh config file.`

func proposeCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addTitleFlag, readTitleFlag := flags.ProposalTitle()
	addBodyFlag, readBodyFlag := flags.ProposalBody()
	addBodyFileFlag, readBodyFileFlag := flags.ProposalBodyFile()
	cmd := cobra.Command{
		Use:     proposeCmd,
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, configdomain.KeyHostingPlatform, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executePropose(readDryRunFlag(cmd), readVerboseFlag(cmd), readTitleFlag(cmd), readBodyFlag(cmd), readBodyFileFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addTitleFlag(&cmd)
	addBodyFlag(&cmd)
	addBodyFileFlag(&cmd)
	return &cmd
}

func executePropose(dryRun configdomain.DryRun, verbose configdomain.Verbose, title gitdomain.ProposalTitle, body gitdomain.ProposalBody, bodyFile gitdomain.ProposalBodyFile) error {
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
	data, exit, err := determineProposeData(repo, dryRun, verbose, title, body, bodyFile)
	if err != nil || exit {
		return err
	}
	if err = validateProposeData(data); err != nil {
		return err
	}
	runProgram := proposeProgram(data)
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
	allBranches       gitdomain.BranchInfos
	branchesSnapshot  gitdomain.BranchesSnapshot
	branchesToSync    gitdomain.BranchInfos
	config            config.ValidatedConfig
	connector         Option[hostingdomain.Connector]
	dialogTestInputs  components.TestInputs
	dryRun            configdomain.DryRun
	hasOpenChanges    bool
	initialBranch     gitdomain.LocalBranchName
	initialBranchType configdomain.BranchType
	previousBranch    Option[gitdomain.LocalBranchName]
	proposalBody      gitdomain.ProposalBody
	proposalTitle     gitdomain.ProposalTitle
	remotes           gitdomain.Remotes
	stashSize         gitdomain.StashSize
}

func determineProposeData(repo execute.OpenRepoResult, dryRun configdomain.DryRun, verbose configdomain.Verbose, title gitdomain.ProposalTitle, body gitdomain.ProposalBody, bodyFile gitdomain.ProposalBodyFile) (data proposeData, exit bool, err error) {
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
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: gitdomain.LocalBranchNames{initialBranch},
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
	initialBranchType := validatedConfig.Config.BranchType(initialBranch)
	var connector Option[hostingdomain.Connector]
	if originURL, hasOriginURL := validatedConfig.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			Config:          *validatedConfig.Config.UnvalidatedConfig,
			HostingPlatform: validatedConfig.Config.HostingPlatform,
			Log:             print.Logger{},
			RemoteURL:       originURL,
		})
		if err != nil {
			return data, false, err
		}
	}
	if connector.IsNone() {
		return data, false, hostingdomain.UnsupportedServiceError()
	}
	branchNamesToSync := validatedConfig.Config.Lineage.BranchAndAncestors(initialBranch)
	branchesToSync, err := branchesSnapshot.Branches.Select(branchNamesToSync...)
	var bodyText gitdomain.ProposalBody
	if body != "" {
		bodyText = body
	} else if bodyFile != "" {
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
		allBranches:       branchesSnapshot.Branches,
		branchesSnapshot:  branchesSnapshot,
		branchesToSync:    branchesToSync,
		config:            validatedConfig,
		connector:         connector,
		dialogTestInputs:  dialogTestInputs,
		dryRun:            dryRun,
		hasOpenChanges:    repoStatus.OpenChanges,
		initialBranch:     initialBranch,
		initialBranchType: initialBranchType,
		previousBranch:    previousBranch,
		proposalBody:      bodyText,
		proposalTitle:     title,
		remotes:           remotes,
		stashSize:         stashSize,
	}, false, err
}

func proposeProgram(data proposeData) program.Program {
	prog := NewMutable(&program.Program{})
	for _, branch := range data.branchesToSync {
		sync.BranchProgram(branch, sync.BranchProgramArgs{
			BranchInfos:   data.allBranches,
			Config:        data.config.Config,
			InitialBranch: data.initialBranch,
			Remotes:       data.remotes,
			Program:       prog,
			PushBranches:  true,
		})
	}
	if data.initialBranchType == configdomain.BranchTypePrototypeBranch {
		prog.Value.Add(&opcodes.RemoveFromPrototypeBranches{Branch: data.initialBranch})
	}
	previousBranchCandidates := gitdomain.LocalBranchNames{}
	if previousBranch, hasPreviousBranch := data.previousBranch.Get(); hasPreviousBranch {
		previousBranchCandidates = append(previousBranchCandidates, previousBranch)
	}
	cmdhelpers.Wrap(prog, cmdhelpers.WrapOptions{
		DryRun:                   data.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         data.hasOpenChanges,
		PreviousBranchCandidates: previousBranchCandidates,
	})
	prog.Value.Add(&opcodes.CreateProposal{
		Branch:        data.initialBranch,
		MainBranch:    data.config.Config.MainBranch,
		ProposalBody:  data.proposalBody,
		ProposalTitle: data.proposalTitle,
	})
	return prog.Get()
}

func validateProposeData(data proposeData) error {
	initialBranchType := data.config.Config.BranchType(data.initialBranch)
	switch initialBranchType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch, configdomain.BranchTypePrototypeBranch:
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
	panic(fmt.Sprintf("unhandled branch type: %v", initialBranchType))
}
