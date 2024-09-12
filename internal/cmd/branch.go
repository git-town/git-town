package cmd

import (
	"os"
	"strings"

	"github.com/git-town/git-town/v16/internal/cli/colors"
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
	printBranches(data)
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
	colors := colors.NewDialogColors()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.Names())
	return branchData{
		branchInfos:      branchesSnapshot.Branches,
		branchesAndTypes: branchesAndTypes,
		colors:           colors,
		initialBranchOpt: branchesSnapshot.Active,
	}, false, err
}

type branchData struct {
	branchInfos       gitdomain.BranchInfos
	branchesAndTypes  configdomain.BranchesAndTypes
	colors            colors.DialogColors
	initialBranchOpt  Option[gitdomain.LocalBranchName]
	defaultBranchType configdomain.DefaultBranchType
}

func printBranches(data branchData) {
	s := strings.Builder{}
	for _, branchInfo := range data.branchInfos {
		branchName, hasLocalBranch := branchInfo.LocalName.Get()
		if !hasLocalBranch {
			continue
		}
		initialBranch, hasInitialBranch := data.initialBranchOpt.Get()
		var isInitial bool
		if hasInitialBranch {
			isInitial = branchName == initialBranch
		}
		switch {
		case isInitial:
			color := data.colors.Initial
			if branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("* " + branchName.String()))
		case branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree:
			s.WriteString(colors.Faint().Styled("+ " + branchName.String()))
		default:
			color := termenv.String()
			if branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree {
				color = color.Faint()
			}
			s.WriteString(color.Styled("  " + branchName.String()))
		}
		branchType, hasBranchType := data.branchesAndTypes[branchName]
		if !hasBranchType {
			branchType = data.defaultBranchType.BranchType
		}
		s.WriteString("  ")
		s.WriteString(colors.Faint().Styled("(" + branchType.String() + ")"))
		s.WriteRune('\n')
	}
}
