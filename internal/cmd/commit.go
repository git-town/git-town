package cmd

import (
	"cmp"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcolors"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/pkg/colors"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	commitDesc = "Commit into another branch"
	commitHelp = `
Allows you to commit the currently staged changes
into another branch without needing to change branches.`
)

func commitCmd() *cobra.Command {
	addMessageFlag, readMessageFlag := flags.CommitMessage("specify the commit message")
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "commit",
		Args:    cobra.ArbitraryArgs,
		GroupID: cmdhelpers.GroupIDStack,
		Short:   commitDesc,
		Long:    cmdhelpers.Long(commitDesc, commitHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			message, errMessage := readMessageFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errMessage, errVerbose); err != nil {
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
			return executeCommit(cliConfig, message)
		},
	}
	addMessageFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeCommit(cliConfig configdomain.PartialConfig, message Option[gitdomain.CommitMessage]) error {
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
	data, flow, err := determineBranchData(repo)
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
	entries := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
		BranchInfos:       data.branchInfos,
		BranchTypes:       []configdomain.BranchType{},
		BranchesAndTypes:  data.branchesAndTypes,
		ExcludeBranches:   gitdomain.LocalBranchNames{},
		Lineage:           repo.UnvalidatedConfig.NormalConfig.Lineage,
		MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
		Order:             repo.UnvalidatedConfig.NormalConfig.Order,
		Regexes:           []*regexp.Regexp{},
		ShowAllBranches:   false,
		UnknownBranchType: repo.UnvalidatedConfig.NormalConfig.UnknownBranchType,
	})
	fmt.Print(branchLayout(entries, data, repo.UnvalidatedConfig.NormalConfig.DisplayTypes))
	return nil
}

func determineBranchData(repo execute.OpenRepoResult) (data branchData, flow configdomain.ProgramFlow, err error) {
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
	initialBranchOpt := branchesSnapshot.Active
	if branchesSnapshot.DetachedHead {
		initialBranch, hasInitialBranch := initialBranchOpt.Get()
		if hasInitialBranch {
			branchesSnapshot.Branches = branchesSnapshot.Branches.Remove(initialBranch)
		}
		initialBranchOpt = None[gitdomain.LocalBranchName]()
	} else if initialBranchOpt.IsNone() {
		if currentBranchOpt, err := repo.Git.CurrentBranchUncached(repo.Backend); err == nil {
			initialBranchOpt = currentBranchOpt
		}
	}
	colors := dialogcolors.NewDialogColors()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.NamesLocalBranches())
	return branchData{
		branchInfos:      branchesSnapshot.Branches,
		branchesAndTypes: branchesAndTypes,
		colors:           colors,
		initialBranchOpt: initialBranchOpt,
	}, configdomain.ProgramFlowContinue, err
}

type branchData struct {
	branchInfos      gitdomain.BranchInfos
	branchesAndTypes configdomain.BranchesAndTypes
	colors           dialogcolors.DialogColors
	initialBranchOpt Option[gitdomain.LocalBranchName]
}

func branchLayout(entries dialog.SwitchBranchEntries, data branchData, displayTypes configdomain.DisplayTypes) string {
	s := strings.Builder{}
	initialBranch, hasInitialBranch := data.initialBranchOpt.Get()
	for _, entry := range entries {
		isInitialBranch := entry.Branch == initialBranch
		switch {
		case hasInitialBranch && isInitialBranch:
			s.WriteString(data.colors.Initial.Styled("* " + entry.String()))
		case entry.OtherWorktree:
			s.WriteString("+ ")
			s.WriteString(colors.Faint().Styled(entry.String()))
		default:
			s.WriteString("  ")
			s.WriteString(entry.String())
		}
		if displayTypes.ShouldDisplayType(entry.Type) {
			s.WriteString("  ")
			s.WriteString(colors.Faint().Styled("(" + entry.Type.String() + ")"))
		}
		s.WriteRune('\n')
	}
	return s.String()
}
