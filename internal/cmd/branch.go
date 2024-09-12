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
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
)

const branchDesc = "Displays all branches, their hierarchy and type"

const branchHelp = `
Git Town's version of Git's "branch" command.`

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
	output := branchLayout(entries, data)
	fmt.Println(output)
	print.Footer(verbose, repo.CommandsCounter.Get(), repo.FinalMessages.Result())
	return nil
}

func determineBranchData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data branchData, exit bool, err error) {
	dialogTestInputs := components.LoadTestInputs(os.Environ())
	repoStatus, err := repo.Git.RepoStatus(repo.Backend)
	branchesSnapshot, _, exit, err := execute.LoadRepoSnapshot(execute.LoadRepoSnapshotArgs{
		Backend:               repo.Backend,
		CommandsCounter:       repo.CommandsCounter,
		ConfigSnapshot:        repo.ConfigSnapshot,
		DialogTestInputs:      dialogTestInputs,
		Fetch:                 false,
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
	defaultBranchType := repo.UnvalidatedConfig.Config.Value.DefaultBranchType
	colors := colors.NewDialogColors()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.Names())
	return branchData{
		branchInfos:       branchesSnapshot.Branches,
		branchesAndTypes:  branchesAndTypes,
		colors:            colors,
		defaultBranchType: defaultBranchType,
		initialBranchOpt:  branchesSnapshot.Active,
		lineage:           repo.UnvalidatedConfig.Config.Value.Lineage,
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
		isInitial := entry.Branch == initialBranch
		switch {
		case hasInitialBranch && isInitial:
			color := data.colors.Initial
			if entry.OtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("* " + entry.String()))
		case entry.OtherWorktree:
			s.WriteString(colors.Faint().Styled("+ " + entry.String()))
		default:
			color := termenv.String()
			if entry.OtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("  " + entry.String()))
		}
		s.WriteString("  ")
		s.WriteString(colors.Faint().Styled("(" + entry.Type.String() + ")"))
		s.WriteRune('\n')
	}
	return s.String()
}
