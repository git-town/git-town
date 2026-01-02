package cmd

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cli/print"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/undo"
	"github.com/git-town/git-town/v22/internal/validate"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const undoDesc = "Undo the most recent Git Town command"

func undoCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "undo",
		GroupID: cmdhelpers.GroupIDErrors,
		Args:    cobra.NoArgs,
		Short:   undoDesc,
		Long:    cmdhelpers.Long(undoDesc),
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
			return executeUndo(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUndo(cliConfig configdomain.PartialConfig) error {
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
	data, flow, err := determineUndoData(repo)
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
	runStateOpt, err := runstate.Load(repo.RootDir)
	if err != nil {
		return fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	runState, hasRunState := runStateOpt.Get()
	if !hasRunState {
		fmt.Println(messages.UndoNothingToDo)
		return nil
	}
	return undo.Execute(undo.ExecuteArgs{
		Backend:          repo.Backend,
		CommandsCounter:  repo.CommandsCounter,
		Config:           data.config,
		Connector:        data.connector,
		FinalMessages:    repo.FinalMessages,
		Frontend:         repo.Frontend,
		Git:              repo.Git,
		HasOpenChanges:   data.hasOpenChanges,
		InitialStashSize: data.stashSize,
		RootDir:          repo.RootDir,
		RunState:         runState,
	})
}

type undoData struct {
	config                  config.ValidatedConfig
	connector               Option[forgedomain.Connector]
	hasOpenChanges          bool
	initialBranchesSnapshot gitdomain.BranchesSnapshot
	inputs                  dialogcomponents.Inputs
	previousBranch          Option[gitdomain.LocalBranchName]
	stashSize               gitdomain.StashSize
}

func determineUndoData(repo execute.OpenRepoResult) (data undoData, flow configdomain.ProgramFlow, err error) {
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
		GiteaToken:           config.GiteaToken,
		GithubConnectorType:  config.GithubConnectorType,
		GithubToken:          config.GithubToken,
		GitlabConnectorType:  config.GitlabConnectorType,
		GitlabToken:          config.GitlabToken,
		Log:                  print.Logger{},
		RemoteURL:            config.DevURL(repo.Backend),
		TestHome:             config.TestHome,
	})
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, stashSize, _, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             connector,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: false,
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
	previousBranch := repo.Git.PreviouslyCheckedOutBranch(repo.Backend)
	return undoData{
		config:                  validatedConfig,
		connector:               connector,
		hasOpenChanges:          repoStatus.OpenChanges,
		initialBranchesSnapshot: branchesSnapshot,
		inputs:                  inputs,
		previousBranch:          previousBranch,
		stashSize:               stashSize,
	}, configdomain.ProgramFlowContinue, nil
}
