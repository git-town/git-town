package cmd

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"slices"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/regexes"
	"github.com/git-town/git-town/v16/internal/validate"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/spf13/cobra"
)

const switchDesc = "Display the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("list both remote-tracking and local branches")
	addDisplayTypesFlag, readDisplayTypesFlag := flags.Displaytypes()
	addMergeFlag, readMergeFlag := flags.SwitchMerge()
	addTypeFlag, readTypeFlag := flags.BranchType()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.ArbitraryArgs,
		Short:   switchDesc,
		Long:    cmdhelpers.Long(switchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchTypes, err := readTypeFlag(cmd)
			if err != nil {
				return err
			}
			allBranches, err := readAllFlag(cmd)
			if err != nil {
				return err
			}
			displayTypes, err := readDisplayTypesFlag(cmd)
			if err != nil {
				return err
			}
			merge, err := readMergeFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			return executeSwitch(args, allBranches, verbose, merge, displayTypes, branchTypes)
		},
	}
	addAllFlag(&cmd)
	addDisplayTypesFlag(&cmd)
	addMergeFlag(&cmd)
	addTypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeSwitch(args []string, allBranches configdomain.AllBranches, verbose configdomain.Verbose, merge configdomain.SwitchUsingMerge, displayTypes configdomain.DisplayTypes, branchTypes []configdomain.BranchType) error {
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
	data, exit, err := determineSwitchData(args, repo, verbose)
	if err != nil || exit {
		return err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(data.branchNames)
	defaultBranchType := repo.UnvalidatedConfig.NormalConfig.DefaultBranchType
	entries := SwitchBranchEntries(data.branchesSnapshot.Branches, branchTypes, branchesAndTypes, data.config.NormalConfig.Lineage, defaultBranchType, allBranches, data.regexes)
	if len(entries) == 0 {
		return errors.New(messages.SwitchNoBranches)
	}
	cursor := SwitchBranchCursorPos(entries, data.initialBranch)
	branchToCheckout, exit, err := dialog.SwitchBranch(entries, cursor, data.uncommittedChanges, displayTypes, data.dialogInputs.Next())
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
	config             config.ValidatedConfig
	dialogInputs       components.TestInputs
	initialBranch      gitdomain.LocalBranchName
	lineage            configdomain.Lineage
	regexes            []*regexp.Regexp
	uncommittedChanges bool
}

func determineSwitchData(args []string, repo execute.OpenRepoResult, verbose configdomain.Verbose) (data switchData, exit bool, err error) {
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
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, exit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	connector, err := hosting.NewConnector(repo.UnvalidatedConfig, repo.UnvalidatedConfig.NormalConfig.DevRemote, print.Logger{})
	if err != nil {
		return data, false, err
	}
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(localBranches)
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: localBranches,
		Connector:          connector,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        NewMutable(&repo.UnvalidatedConfig),
	})
	if err != nil || exit {
		return data, exit, err
	}
	regexes, err := regexes.NewRegexes(args)
	if err != nil {
		return data, false, err
	}
	return switchData{
		branchNames:        branchesSnapshot.Branches.Names(),
		branchesSnapshot:   branchesSnapshot,
		config:             validatedConfig,
		dialogInputs:       dialogTestInputs,
		initialBranch:      initialBranch,
		lineage:            validatedConfig.NormalConfig.Lineage,
		regexes:            regexes,
		uncommittedChanges: repoStatus.OpenChanges,
	}, false, err
}

// SwitchBranchCursorPos provides the initial cursor position for the "switch branch" components.
func SwitchBranchCursorPos(entries []dialog.SwitchBranchEntry, initialBranch gitdomain.LocalBranchName) int {
	for e, entry := range entries {
		if entry.Branch == initialBranch {
			return e
		}
	}
	return 0
}

// SwitchBranchEntries provides the entries for the "switch branch" components.
func SwitchBranchEntries(branchInfos gitdomain.BranchInfos, branchTypes []configdomain.BranchType, branchesAndTypes configdomain.BranchesAndTypes, lineage configdomain.Lineage, defaultBranchType configdomain.BranchType, allBranches configdomain.AllBranches, regexes []*regexp.Regexp) []dialog.SwitchBranchEntry {
	entries := make([]dialog.SwitchBranchEntry, 0, lineage.Len())
	roots := lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(&entries, root, "", lineage, branchInfos, allBranches, branchTypes, branchesAndTypes, defaultBranchType, regexes)
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
		layoutBranches(&entries, localBranch, "", lineage, branchInfos, allBranches, branchTypes, branchesAndTypes, defaultBranchType, regexes)
	}
	return entries
}

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(result *[]dialog.SwitchBranchEntry, branch gitdomain.LocalBranchName, indentation string, lineage configdomain.Lineage, branchInfos gitdomain.BranchInfos, allBranches configdomain.AllBranches, branchTypes []configdomain.BranchType, branchesAndTypes configdomain.BranchesAndTypes, defaultBranchType configdomain.BranchType, regexes regexes.Regexes) {
	if branchInfos.HasLocalBranch(branch) || allBranches.Enabled() {
		var otherWorktree bool
		if branchInfo, hasBranchInfo := branchInfos.FindByLocalName(branch).Get(); hasBranchInfo {
			otherWorktree = branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		} else {
			otherWorktree = false
		}
		branchType, hasBranchType := branchesAndTypes[branch]
		if !hasBranchType && len(branchTypes) > 0 {
			branchType = defaultBranchType
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
		layoutBranches(result, child, indentation+"  ", lineage, branchInfos, allBranches, branchTypes, branchesAndTypes, defaultBranchType, regexes)
	}
}
