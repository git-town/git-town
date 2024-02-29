package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/git-town/git-town/v12/src/cli/dialog/components"
	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cli/print"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git/gitdomain"
	"github.com/git-town/git-town/v12/src/hosting"
	"github.com/git-town/git-town/v12/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/git-town/git-town/v12/src/sync"
	"github.com/git-town/git-town/v12/src/undo/undoconfig"
	fullInterpreter "github.com/git-town/git-town/v12/src/vm/interpreter/full"
	"github.com/git-town/git-town/v12/src/vm/opcodes"
	"github.com/git-town/git-town/v12/src/vm/program"
	"github.com/git-town/git-town/v12/src/vm/runstate"
	"github.com/spf13/cobra"
)

const proposeDesc = "Creates a proposal to merge a feature branch"

const proposeHelp = `
Syncs the current branch and opens a browser window to the new proposal page of your repository.

The form is pre-populated for the current branch so that the proposal only shows the changes made against the immediate parent branch.

Supported only for repositories hosted on GitHub, GitLab, Gitea and Bitbucket. When using self-hosted versions this command needs to be configured with "git config %s <driver>" where driver is "github", "gitlab", "gitea", or "bitbucket". When using SSH identities, this command needs to be configured with "git config %s <hostname>" where hostname matches what is in your ssh config file.`

func proposeCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "propose",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, gitconfig.KeyHostingPlatform, gitconfig.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
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
	config, initialBranchesSnapshot, initialStashSize, exit, err := determineProposeConfig(repo, dryRun, verbose)
	if err != nil || exit {
		return err
	}
	if err = validateProposeConfig(config); err != nil {
		return err
	}
	runState := runstate.RunState{
		BeginBranchesSnapshot: initialBranchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        initialStashSize,
		Command:               "propose",
		DryRun:                dryRun,
		EndBranchesSnapshot:   gitdomain.EmptyBranchesSnapshot(),
		EndConfigSnapshot:     undoconfig.EmptyConfigSnapshot(),
		EndStashSize:          0,
		RunProgram:            proposeProgram(config),
	}
	return fullInterpreter.Execute(fullInterpreter.ExecuteArgs{
		Connector:               config.connector,
		DialogTestInputs:        &config.dialogTestInputs,
		FullConfig:              config.FullConfig,
		HasOpenChanges:          config.hasOpenChanges,
		InitialBranchesSnapshot: initialBranchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        initialStashSize,
		RootDir:                 repo.RootDir,
		Run:                     repo.Runner,
		RunState:                &runState,
		Verbose:                 verbose,
	})
}

type proposeConfig struct {
	*configdomain.FullConfig
	allBranches      gitdomain.BranchInfos
	branchesToSync   gitdomain.BranchInfos
	connector        hostingdomain.Connector
	dialogTestInputs components.TestInputs
	dryRun           bool
	hasOpenChanges   bool
	initialBranch    gitdomain.LocalBranchName
	previousBranch   gitdomain.LocalBranchName
	remotes          gitdomain.Remotes
}

func determineProposeConfig(repo *execute.OpenRepoResult, dryRun, verbose bool) (*proposeConfig, gitdomain.BranchesSnapshot, gitdomain.StashSize, bool, error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	branchesSnapshot, stashSize, repoStatus, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 true,
		FullConfig:            &repo.Runner.Config.FullConfig,
		HandleUnfinishedState: true,
		Repo:                  repo,
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
	err = execute.EnsureKnownBranchAncestry(branchesSnapshot.Active, execute.EnsureKnownBranchAncestryArgs{
		Config:           &repo.Runner.Config.FullConfig,
		AllBranches:      branchesSnapshot.Branches,
		DefaultBranch:    repo.Runner.Config.FullConfig.MainBranch,
		DialogTestInputs: &dialogTestInputs,
		Runner:           repo.Runner,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	originURL := repo.Runner.Config.OriginURL()
	connector, err := hosting.NewConnector(hosting.NewConnectorArgs{
		FullConfig:      &repo.Runner.Config.FullConfig,
		HostingPlatform: repo.Runner.Config.FullConfig.HostingPlatform,
		Log:             print.Logger{},
		OriginURL:       originURL,
	})
	if err != nil {
		return nil, branchesSnapshot, stashSize, false, err
	}
	if connector == nil {
		return nil, branchesSnapshot, stashSize, false, hostingdomain.UnsupportedServiceError()
	}
	branchNamesToSync := repo.Runner.Config.FullConfig.Lineage.BranchAndAncestors(branchesSnapshot.Active)
	branchesToSync, err := branchesSnapshot.Branches.Select(branchNamesToSync)
	return &proposeConfig{
		FullConfig:       &repo.Runner.Config.FullConfig,
		allBranches:      branchesSnapshot.Branches,
		branchesToSync:   branchesToSync,
		connector:        connector,
		dialogTestInputs: dialogTestInputs,
		dryRun:           dryRun,
		hasOpenChanges:   repoStatus.OpenChanges,
		initialBranch:    branchesSnapshot.Active,
		previousBranch:   previousBranch,
		remotes:          remotes,
	}, branchesSnapshot, stashSize, false, err
}

func proposeProgram(config *proposeConfig) program.Program {
	prog := program.Program{}
	for _, branch := range config.branchesToSync {
		sync.BranchProgram(branch, sync.BranchProgramArgs{
			Config:        config.FullConfig,
			BranchInfos:   config.allBranches,
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

func validateProposeConfig(config *proposeConfig) error {
	initialBranchType := config.FullConfig.BranchType(config.initialBranch)
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
