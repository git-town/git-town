package cmd

import (
	"errors"
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
	"github.com/git-town/git-town/v22/internal/skip"
	"github.com/git-town/git-town/v22/internal/state/runstate"
	"github.com/git-town/git-town/v22/internal/validate"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const skipDesc = "Resume the last run Git Town command by skipping the current branch"

func skipCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addParkFlag, readParkFlag := flags.Park()
	cmd := cobra.Command{
		Use:     "skip",
		GroupID: cmdhelpers.GroupIDErrors,
		Args:    cobra.NoArgs,
		Short:   skipDesc,
		Long:    cmdhelpers.Long(skipDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			park, err := readParkFlag(cmd)
			if err != nil {
				return err
			}
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
			return executeSkip(cliConfig, park)
		},
	}
	addParkFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSkip(cliConfig configdomain.PartialConfig, park configdomain.Park) error {
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
	data, flow, err := loadSkipData(repo)
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
	if unfinishedDetails, hasUnfinishedDetails := data.runState.UnfinishedDetails.Get(); hasUnfinishedDetails {
		if !unfinishedDetails.CanSkip {
			return errors.New(messages.SkipBranchHasConflicts)
		}
	}
	if park {
		activeBranchType := data.config.BranchType(data.activeBranch)
		if err = canParkBranchType(activeBranchType, data.activeBranch, repo.FinalMessages); err != nil {
			return err
		}
	}
	return skip.Execute(skip.ExecuteArgs{
		Backend:         repo.Backend,
		CommandsCounter: repo.CommandsCounter,
		Config:          data.config,
		Connector:       data.connector,
		FinalMessages:   repo.FinalMessages,
		Frontend:        repo.Frontend,
		Git:             repo.Git,
		HasOpenChanges:  data.hasOpenChanges,
		InitialBranch:   data.activeBranch,
		Inputs:          data.inputs,
		Park:            park,
		RootDir:         repo.RootDir,
		RunState:        data.runState,
	})
}

func loadSkipData(repo execute.OpenRepoResult) (skipData, configdomain.ProgramFlow, error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	var emptyResult skipData
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
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
	})
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	branchesSnapshot, _, _, flow, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
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
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	switch flow {
	case configdomain.ProgramFlowContinue:
	case configdomain.ProgramFlowExit, configdomain.ProgramFlowRestart:
		return emptyResult, flow, nil
	}
	activeBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
	if !hasCurrentBranch {
		currentBranchOpt, err := repo.Git.CurrentBranch(repo.Backend)
		if err != nil {
			return emptyResult, configdomain.ProgramFlowExit, err
		}
		if currentBranch, has := currentBranchOpt.Get(); has {
			activeBranch = currentBranch
		}
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().NamesLocalBranches()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().NamesLocalBranches())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: localBranches,
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
		return emptyResult, configdomain.ProgramFlowExit, err
	}
	runStateOpt, err := runstate.Load(repo.RootDir)
	if err != nil {
		return emptyResult, configdomain.ProgramFlowExit, fmt.Errorf(messages.RunstateLoadProblem, err)
	}
	runState, hasRunState := runStateOpt.Get()
	if !hasRunState || runState.IsFinished() {
		return emptyResult, configdomain.ProgramFlowExit, errors.New(messages.SkipNothingToDo)
	}
	return skipData{
		activeBranch:   activeBranch,
		config:         validatedConfig,
		connector:      connector,
		hasOpenChanges: repoStatus.OpenChanges,
		inputs:         inputs,
		runState:       runState,
	}, configdomain.ProgramFlowContinue, nil
}

type skipData struct {
	activeBranch   gitdomain.LocalBranchName
	config         config.ValidatedConfig
	connector      Option[forgedomain.Connector]
	hasOpenChanges bool
	inputs         dialogcomponents.Inputs
	runState       runstate.RunState
}
