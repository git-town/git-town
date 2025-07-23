package cmd

import (
	"cmp"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"slices"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/regexes"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const switchDesc = "Display the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("list both remote-tracking and local branches")
	addDisplayTypesFlag, readDisplayTypesFlag := flags.Displaytypes()
	addMergeFlag, readMergeFlag := flags.Merge()
	addTypeFlag, readTypeFlag := flags.BranchType()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.ArbitraryArgs,
		Short:   switchDesc,
		Long:    cmdhelpers.Long(switchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchTypes, err1 := readTypeFlag(cmd)
			allBranches, err2 := readAllFlag(cmd)
			displayTypes, err3 := readDisplayTypesFlag(cmd)
			merge, err4 := readMergeFlag(cmd)
			verbose, err5 := readVerboseFlag(cmd)
			if err := cmp.Or(err1, err2, err3, err4, err5); err != nil {
				return err
			}
			cliConfig := cliconfig.CliConfig{
				DryRun:  false,
				Verbose: verbose,
			}
			return executeSwitch(args, cliConfig, allBranches, merge, displayTypes, branchTypes)
		},
	}
	addAllFlag(&cmd)
	addDisplayTypesFlag(&cmd)
	addMergeFlag(&cmd)
	addTypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwitch(args []string, cliConfig cliconfig.CliConfig, allBranches configdomain.AllBranches, merge configdomain.SwitchUsingMerge, displayTypes configdomain.DisplayTypes, branchTypes []configdomain.BranchType) error {
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
	data, exit, err := determineSwitchData(args, repo, cliConfig)
	if err != nil || exit {
		return err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(data.branchNames)
	unknownBranchType := repo.UnvalidatedConfig.NormalConfig.UnknownBranchType
	entries := SwitchBranchEntries(data.branchesSnapshot.Branches, branchTypes, branchesAndTypes, data.config.NormalConfig.Lineage, unknownBranchType, allBranches, data.regexes)
	if len(entries) == 0 {
		return errors.New(messages.SwitchNoBranches)
	}
	cursor := SwitchBranchCursorPos(entries, data.initialBranch)
	branchToCheckout, exit, err := dialog.SwitchBranch(entries, cursor, data.uncommittedChanges, displayTypes, data.inputs)
	if err != nil || exit {
		return err
	}
	if branchToCheckout == data.initialBranch {
		return nil
	}
	err = repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, merge)
	if err != nil {
		exitCode := 1
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
		os.Exit(exitCode)
	}
	return nil
}

type switchData struct {
	branchNames        gitdomain.LocalBranchNames
	branchesSnapshot   gitdomain.BranchesSnapshot
	config             config.UnvalidatedConfig
	initialBranch      gitdomain.LocalBranchName
	inputs             dialogcomponents.Inputs
	lineage            configdomain.Lineage
	regexes            []*regexp.Regexp
	uncommittedChanges bool
}

func determineSwitchData(args []string, repo execute.OpenRepoResult, cliConfig cliconfig.CliConfig) (data switchData, exit dialogdomain.Exit, err error) {
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
		HandleUnfinishedState: true,
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
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	regexes, err := regexes.NewRegexes(args)
	if err != nil {
		return data, false, err
	}
	return switchData{
		branchNames:        branchesSnapshot.Branches.Names(),
		branchesSnapshot:   branchesSnapshot,
		config:             repo.UnvalidatedConfig,
		initialBranch:      initialBranch,
		inputs:             inputs,
		lineage:            repo.UnvalidatedConfig.NormalConfig.Lineage,
		regexes:            regexes,
		uncommittedChanges: repoStatus.OpenChanges,
	}, false, err
}

// SwitchBranchCursorPos provides the initial cursor position for the "switch branch" components.
func SwitchBranchCursorPos(entries dialog.SwitchBranchEntries, initialBranch gitdomain.LocalBranchName) int {
	for e, entry := range entries {
		if entry.Branch == initialBranch {
			return e
		}
	}
	return 0
}

// SwitchBranchEntries provides the entries for the "switch branch" components.
func SwitchBranchEntries(branchInfos gitdomain.BranchInfos, branchTypes []configdomain.BranchType, branchesAndTypes configdomain.BranchesAndTypes, lineage configdomain.Lineage, unknownBranchType configdomain.UnknownBranchType, allBranches configdomain.AllBranches, regexes []*regexp.Regexp) dialog.SwitchBranchEntries {
	entries := make(dialog.SwitchBranchEntries, 0, lineage.Len())
	roots := lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(&entries, root, "", lineage, branchInfos, allBranches, branchTypes, branchesAndTypes, unknownBranchType, regexes)
	}
	// add branches not in the lineage
	branchesInLineage := lineage.BranchesWithParents()
	for _, branchInfo := range branchInfos {
		localBranch := branchInfo.LocalBranchName()
		if slices.Contains(roots, localBranch) {
			continue
		}
		if slices.Contains(branchesInLineage, localBranch) {
			continue
		}
		if entries.ContainsBranch(localBranch) {
			continue
		}
		layoutBranches(&entries, localBranch, "", lineage, branchInfos, allBranches, branchTypes, branchesAndTypes, unknownBranchType, regexes)
	}
	return entries
}

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(result *dialog.SwitchBranchEntries, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage, branchInfos gitdomain.BranchInfos, allBranches configdomain.AllBranches, branchTypes []configdomain.BranchType, branchesAndTypes configdomain.BranchesAndTypes, unknownBranchType configdomain.UnknownBranchType, regexes regexes.Regexes) {
	if branchInfos.HasLocalBranch(branch) || allBranches.Enabled() {
		var otherWorktree bool
		if branchInfo, hasBranchInfo := branchInfos.FindByLocalName(branch).Get(); hasBranchInfo {
			otherWorktree = branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		} else {
			otherWorktree = false
		}
		branchType, hasBranchType := branchesAndTypes[branch]
		if !hasBranchType && len(branchTypes) > 0 {
			branchType = unknownBranchType.BranchType()
		}
		var hasCorrectBranchType bool
		if len(branchTypes) == 0 || slices.Contains(branchTypes, branchType) {
			hasCorrectBranchType = true
		}
		matchesRegex := regexes.Matches(branch.String())
		if hasCorrectBranchType && matchesRegex {
			*result = append(*result, dialog.SwitchBranchEntry{
				Branch:        branch,
				Indentation:   indentation,
				OtherWorktree: otherWorktree,
				Type:          branchType,
			})
		}
	}
	for _, child := range lineage.Children(branch) {
		layoutBranches(result, child, indentation+"  ", lineage, branchInfos, allBranches, branchTypes, branchesAndTypes, unknownBranchType, regexes)
	}
}
