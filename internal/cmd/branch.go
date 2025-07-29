package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/state/runstate"
	"github.com/git-town/git-town/v21/internal/undo/undoconfig"
	"github.com/git-town/git-town/v21/internal/validate"
	"github.com/git-town/git-town/v21/internal/vm/interpreter/fullinterpreter"
	"github.com/git-town/git-town/v21/internal/vm/opcodes"
	"github.com/git-town/git-town/v21/internal/vm/optimizer"
	"github.com/git-town/git-town/v21/internal/vm/program"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	branchCommand = "branch"
	branchDesc    = `
Display hierarchy either for local branches or an existing proposal of stack.
	`
	branchHelp = `
Git Town's equivalent of the "git branch" command.`
)

func branchCmd() *cobra.Command {
	addProposalLineageFlag, readProposalFlag := flags.Proposal("Display the proposal stack lineage based on the current branch")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   branchCommand,
		Args:  cobra.NoArgs,
		Short: branchDesc,
		Long:  cmdhelpers.Long(branchDesc, branchHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			proposalLineageAction, err := readProposalFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  false,
				Verbose: verbose,
			}
			return executeBranch(cliConfig, proposalLineageAction)
		},
	}

	addProposalLineageFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeBranch(cliConfig cliconfig.CliConfig, proposal configdomain.Proposal) error {
	if proposal {
		return executeProposalDisplay(cliConfig)
	}

	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineBranchData(repo, cliConfig)
	if err != nil || exit {
		return err
	}
	entries := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
		BranchInfos:       data.branchInfos,
		BranchTypes:       []configdomain.BranchType{},
		BranchesAndTypes:  data.branchesAndTypes,
		ExcludeBranches:   gitdomain.LocalBranchNames{},
		Lineage:           repo.UnvalidatedConfig.NormalConfig.Lineage,
		MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
		Regexes:           []*regexp.Regexp{},
		ShowAllBranches:   false,
		UnknownBranchType: repo.UnvalidatedConfig.NormalConfig.UnknownBranchType,
	})
	fmt.Print(branchLayout(entries, data))
	return nil
}

func executeProposalDisplay(cliConfig cliconfig.CliConfig) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    false,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}

	inputs := dialogcomponents.LoadInputs(os.Environ())
	branchesSnapshot, err := repo.Git.BranchesSnapshot(repo.Backend)
	if err != nil {
		return err
	}

	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return err
	}

	currentBranch, err := repo.Git.CurrentBranchUncached(repo.Backend)
	if err != nil {
		return err
	}

	runProgram := NewMutable(&program.Program{})

	config := repo.UnvalidatedConfig.NormalConfig
	connector, err := forge.NewConnector(forge.NewConnectorArgs{
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
		return err
	}

	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	remotes, err := repo.Git.Remotes(repo.Backend)
	if err != nil {
		return err
	}

	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchInfos:        branchesSnapshot.Branches,
		BranchesAndTypes:   branchesAndTypes,
		BranchesToValidate: gitdomain.LocalBranchNames{currentBranch},
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
		return err
	}

	runProgram.Value.Add(&opcodes.ProposalLineageCreate{
		Branch:            currentBranch,
		ProposalLineageIn: configdomain.ProposalLineageInTerminal,
	})

	optimizedProgram := optimizer.Optimize(runProgram.Immutable())

	runState := runstate.RunState{
		BeginBranchesSnapshot: branchesSnapshot,
		BeginConfigSnapshot:   repo.ConfigSnapshot,
		BeginStashSize:        0,
		Command:               branchCommand,
		DryRun:                cliConfig.DryRun,
		EndBranchesSnapshot:   None[gitdomain.BranchesSnapshot](),
		EndConfigSnapshot:     None[undoconfig.ConfigSnapshot](),
		EndStashSize:          None[gitdomain.StashSize](),
		BranchInfosLastRun:    None[gitdomain.BranchInfos](),
		RunProgram:            optimizedProgram,
		TouchedBranches:       optimizedProgram.TouchedBranches(),
		UndoAPIProgram:        program.Program{},
	}
	return fullinterpreter.Execute(fullinterpreter.ExecuteArgs{
		Backend:                 repo.Backend,
		CommandsCounter:         repo.CommandsCounter,
		Config:                  validatedConfig,
		Connector:               connector,
		Detached:                false,
		FinalMessages:           repo.FinalMessages,
		Frontend:                repo.Frontend,
		Git:                     repo.Git,
		HasOpenChanges:          false,
		InitialBranch:           currentBranch,
		InitialBranchesSnapshot: branchesSnapshot,
		InitialConfigSnapshot:   repo.ConfigSnapshot,
		InitialStashSize:        0,
		Inputs:                  dialogcomponents.LoadInputs(os.Environ()),
		PendingCommand:          None[string](),
		RootDir:                 repo.RootDir,
		RunState:                runState,
		Verbose:                 cliConfig.Verbose,
	})
}

func determineBranchData(repo execute.OpenRepoResult, cliConfig cliconfig.CliConfig) (data branchData, exit dialogdomain.Exit, err error) {
	inputs := dialogcomponents.LoadInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		Connector:             None[forgedomain.Connector](),
		Detached:              true,
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
		Verbose:               cliConfig.Verbose,
	})
	if err != nil || exit {
		return data, exit, err
	}
	initialBranchOpt := branchesSnapshot.Active
	if initialBranchOpt.IsNone() {
		initialBranch, err := repo.Git.CurrentBranchUncached(repo.Backend)
		if err == nil {
			initialBranchOpt = Some(initialBranch)
		}
	}
	colors := colors.NewDialogColors()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.Names())
	return branchData{
		branchInfos:      branchesSnapshot.Branches,
		branchesAndTypes: branchesAndTypes,
		colors:           colors,
		initialBranchOpt: initialBranchOpt,
	}, false, err
}

type branchData struct {
	branchInfos      gitdomain.BranchInfos
	branchesAndTypes configdomain.BranchesAndTypes
	colors           colors.DialogColors
	initialBranchOpt Option[gitdomain.LocalBranchName]
}

func branchLayout(entries dialog.SwitchBranchEntries, data branchData) string {
	s := strings.Builder{}
	initialBranch, hasInitialBranch := data.initialBranchOpt.Get()
	for _, entry := range entries {
		isInitialBranch := entry.Branch == initialBranch
		switch {
		case hasInitialBranch && isInitialBranch:
			s.WriteString(data.colors.Initial.Styled("* " + entry.String()))
		case entry.OtherWorktree:
			s.WriteString("+ ")
			s.WriteString(colors.Cyan().Styled(entry.String()))
		default:
			s.WriteString("  ")
			s.WriteString(entry.String())
		}
		if dialog.ShouldDisplayBranchType(entry.Type) {
			s.WriteString("  ")
			s.WriteString(colors.Faint().Styled("(" + entry.Type.String() + ")"))
		}
		s.WriteRune('\n')
	}
	return s.String()
}
