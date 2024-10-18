package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v16/internal/cli/colors"
	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const branchDesc = "Display the local branch hierarchy and types"

const branchHelp = `
Git Town's equivalent of the "git branch" command.`

func branchCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:   "branch",
		Args:  cobra.NoArgs,
		Short: branchDesc,
		Long:  cmdhelpers.Long(branchDesc, branchHelp),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return executeBranch(readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeBranch(verbose configdomain.Verbose) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineBranchData(repo, verbose)
	if err != nil || exit {
		return err
	}
	entries := SwitchBranchEntries(data.branchInfos, []configdomain.BranchType{}, data.branchesAndTypes, data.lineage, data.defaultBranchType, false, []*regexp.Regexp{})
	fmt.Print(branchLayout(entries, data))
	return nil
}

func determineBranchData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data branchData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	if err != nil {
		return data, false, err
	}
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
		FinalMessages:         repo.FinalMessages,
		Frontend:              repo.Frontend,
		Git:                   repo.Git,
		HandleUnfinishedState: false,
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
	initialBranchOpt := branchesSnapshot.Active
	if initialBranchOpt.IsNone() {
		initialBranch, err := repo.Git.CurrentBranchUncached(repo.Backend)
		if err == nil {
			initialBranchOpt = Some(initialBranch)
		}
	}
	defaultBranchType := repo.UnvalidatedConfig.NormalConfig.Config.DefaultBranchType
	colors := colors.NewDialogColors()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(branchesSnapshot.Branches.Names())
	return branchData{
		branchInfos:       branchesSnapshot.Branches,
		branchesAndTypes:  branchesAndTypes,
		colors:            colors,
		defaultBranchType: defaultBranchType,
		initialBranchOpt:  initialBranchOpt,
		lineage:           repo.UnvalidatedConfig.NormalConfig.Config.Lineage,
	}, false, err
}

type branchData struct {
	branchInfos       gitdomain.BranchInfos
	branchesAndTypes  configdomain.BranchesAndTypes
	colors            colors.DialogColors
	defaultBranchType configdomain.DefaultBranchType
	initialBranchOpt  Option[gitdomain.LocalBranchName]
	lineage           configdomain.Lineage
}

func branchLayout(entries []dialog.SwitchBranchEntry, data branchData) string {
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
