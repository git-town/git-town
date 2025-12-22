package cmd

import (
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/gitconfig"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/setup"
	"github.com/git-town/git-town/v22/internal/vm/interpreter/configinterpreter"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	initCmd  = "init"
	initDesc = "Prompts to setup your Git Town configuration"
)

func initCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     initCmd,
		Args:    cobra.NoArgs,
		GroupID: cmdhelpers.GroupIDConfig,
		Short:   initDesc,
		Long:    cmdhelpers.Long(initDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      None[configdomain.DisplayTypes](),
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             None[configdomain.Order](),
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeConfigSetup(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeConfigSetup(cliConfig configdomain.PartialConfig) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		IgnoreUnknown:    true,
		PrintBranchNames: false,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := LoadData(repo)
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
	userInput, exit, enterAll, err := setup.Enter(data)
	if err != nil || exit {
		return err
	}
	if err = setup.Save(userInput, repo.UnvalidatedConfig, data, enterAll, repo.Frontend); err != nil {
		return err
	}
	return configinterpreter.Finished(configinterpreter.FinishedArgs{
		Backend:               repo.Backend,
		BeginBranchesSnapshot: None[gitdomain.BranchesSnapshot](),
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		Command:               initCmd,
		CommandsCounter:       repo.CommandsCounter,
		FinalMessages:         repo.FinalMessages,
		Git:                   repo.Git,
		RootDir:               repo.RootDir,
		TouchedBranches:       []gitdomain.BranchName{},
		Verbose:               repo.UnvalidatedConfig.NormalConfig.Verbose,
	})
}

func LoadData(repo execute.OpenRepoResult) (data setup.Data, flow configdomain.ProgramFlow, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, _, _, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             None[forgedomain.Connector](),
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
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	if len(remotes) == 0 {
		remotes = gitdomain.Remotes{gitconfig.DefaultRemote(repo.Backend)}
	}
	return setup.Data{
		Backend:       repo.Backend,
		Config:        repo.UnvalidatedConfig,
		Git:           repo.Git,
		Inputs:        inputs,
		LocalBranches: branchesSnapshot.Branches.LocalBranches().NamesLocalBranches(),
		Remotes:       remotes,
		Snapshot:      repo.ConfigSnapshot,
	}, configdomain.ProgramFlowContinue, nil
}
