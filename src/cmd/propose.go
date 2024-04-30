package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/execute"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/hosting"
	"github.com/git-town/git-town/v14/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/sync"
	"github.com/git-town/git-town/v14/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v14/src/vm/interpreter/full"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/runstate"
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
	cmd := cobra.Command{
		Use:     proposeCmd,
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, gitconfig.KeyHostingPlatform, gitconfig.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executePropose(readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executePropose(dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           dryRun,
		OmitBranchNames:  false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: true,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, initialBranchesSnapshot, initialStashSize, exit, err := determineProposeData(repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	if err = validateProposeData(data); err != nil {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               proposeCmd,
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            proposeProgram(data),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Config:                  data.config,
		Connector:               data.connector,
		DialogTestInputs:        &data.dialogTestInputs,
		HasOpenChanges:          data.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                runState,
		Verbose:                 verbose,
	})
}

type proposeData struct {
	allBranches      gitdomain.BranchInfos
	branchesToSync   gitdomain.BranchInfos
	config           configdomain.FullConfig
	connector        hostingdomain.Connector
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   gitdomain.LocalBranchName
	remotes          gitdomain.Remotes
}

func determineProposeData(repo *execute.OpenRepoResult, dryRun, verbose bool) (*proposeData, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Runner.Backend.RepoStatus()
	if err != nil {
		return nil, gitdomain.EmptyBranchesSnapshot(), 0, false, err
	}
	branchesSnapshot, stashSize, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Config:                repo.Runner.Config,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		HandleUnfinishedState: true,
		Repo:                  repo,
		RepoStatus:            repoStatus,
		ValidateIsConfigured:  true,
		ValidateNoOpenChanges: false,
		Verbose:               verbose,
	})
	if err != nil || exit {
		return nil, branchesSnapshot, stashSize, exit, err
	}
	previousBranch := repo.Runner.Backend.PreviouslyCheckedOutBranch()
	remotes, err := repo.Runner.Backend.Remotes()
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	err = execute.EnsureKnownBranchesAncestry(execute.EnsureKnownBranchesAncestryArgs{
		BranchesToVerify: gitdomain.LocalBranchNames{branchesSnapshot.Active},
		Config:           repo.Runner.Config,
		DefaultChoice:    repo.Runner.Config.FullConfig.MainBranch,
		DialogTestInputs: &dialogTestInputs,
		LocalBranches:    branchesSnapshot.Branches,
		MainBranch:       repo.Runner.Config.FullConfig.MainBranch,
		Runner:           repo.Runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	var connector hostingdomain.Connector
	if originURL, hasOriginURL := repo.Runner.Config.OriginURL().Get(); hasOriginURL {
		connector, err = hosting.NewConnector(hosting.NewConnectorArgs{
			FullConfig:      &repo.Runner.Config.FullConfig,
			HostingPlatform: repo.Runner.Config.FullConfig.HostingPlatform,
			Log:             print.Logger{},
			OriginURL:       originURL,
		})
		if err != nil {
			return nil, branchesSnapshot, stashSize, false, err
		}
	}
	if connector == nil {
		return nil, branchesSnapshot, stashSize, false, hostingdomain.UnsupportedServiceError()
	}
	branchNamesToSync := repo.Runner.Config.FullConfig.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	branchesToSync, err := branchesSnapshot.Branches.Select(branchNamesToSync...)
	return &proposeData{
		allBranches:      branchesSnapshot.Branches,
		branchesToSync:   branchesToSync,
		config:           repo.Runner.Config.FullConfig,
		connector:        connector,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		previousBranch:   previousBranch,
		remotes:          remotes,
	}, branchesSnapshot, stashSize, false, err
}

func proposeProgram(config *proposeData) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		sync.BranchProgram(branch, sync.BranchProgramArgs{
			BranchInfos:   config.allBranches,
			Config:        config.config,
			InitialBranch: config.initialBranch,
			Remotes:       config.remotes,
			Program:       &prog,
			PushBranch:    true,
		})
	}
	cmdhelpers.Wrap(&prog, cmdhelpers.WrapOptions{
		DryRun:                   config.dryRun,
		RunInGitRoot:             true,
		StashOpenChanges:         config.hasOpenChanges,
		PreviousBranchCandidates: gitdomain.LocalBranchNames{config.previousBranch},
	})
	prog.Add(&opcodes.CreateProposal{Branch: config.initialBranch})
	return prog
}

func validateProposeData(data *proposeData) error {
	initialBranchType := data.config.BranchType(data.initialBranch)
	switch initialBranchType {
	case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeParkedBranch:
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
