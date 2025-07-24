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
	entries := SwitchBranchEntries(SwitchBranchArgs{
		ShowAllBranches:   allBranches,
		BranchInfos:       data.branchesSnapshot.Branches,
		BranchTypes:       branchTypes,
		BranchesAndTypes:  branchesAndTypes,
		ExcludeBranches:   gitdomain.LocalBranchNames{},
		Lineage:           data.config.NormalConfig.Lineage,
		Regexes:           data.regexes,
		UnknownBranchType: unknownBranchType,
	})
	if len(entries) == 0 {
		return errors.New(messages.SwitchNoBranches)
	}
	cursor := entries.IndexOf(data.initialBranch)
	branchToCheckout, exit, err := dialog.SwitchBranch(dialog.SwitchBranchArgs{
		CurrentBranch:      Some(data.initialBranch),
		Cursor:             cursor,
		DisplayBranchTypes: displayTypes,
		Entries:            entries,
		InputName:          "switch-branch",
		Inputs:             data.inputs,
		UncommittedChanges: data.uncommittedChanges,
	})
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

// SwitchBranchEntries provides the entries for the "switch branch" components.
func SwitchBranchEntries(args SwitchBranchArgs) dialog.SwitchBranchEntries {
	entries := make(dialog.SwitchBranchEntries, 0, args.Lineage.Len())
	roots := args.Lineage.Roots()
	// add all entries from the lineage
	for _, root := range roots {
		layoutBranches(layoutBranchesArgs{
			branch:            root,
			branchInfos:       args.BranchInfos,
			branchTypes:       args.BranchTypes,
			branchesAndTypes:  args.BranchesAndTypes,
			excludeBranches:   args.ExcludeBranches,
			indentation:       "",
			lineage:           args.Lineage,
			regexes:           args.Regexes,
			result:            &entries,
			showAllBranches:   args.ShowAllBranches,
			unknownBranchType: args.UnknownBranchType,
		})
	}
	// add branches not in the lineage
	branchesInLineage := args.Lineage.BranchesWithParents()
	for _, branchInfo := range args.BranchInfos {
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
		layoutBranches(layoutBranchesArgs{
			branch:            localBranch,
			branchInfos:       args.BranchInfos,
			branchTypes:       args.BranchTypes,
			branchesAndTypes:  args.BranchesAndTypes,
			excludeBranches:   args.ExcludeBranches,
			indentation:       "",
			lineage:           args.Lineage,
			regexes:           args.Regexes,
			result:            &entries,
			showAllBranches:   args.ShowAllBranches,
			unknownBranchType: args.UnknownBranchType,
		})
	}
	return entries
}

type SwitchBranchArgs struct {
	BranchInfos       gitdomain.BranchInfos
	BranchTypes       []configdomain.BranchType
	BranchesAndTypes  configdomain.BranchesAndTypes
	ExcludeBranches   gitdomain.LocalBranchNames
	Lineage           configdomain.Lineage
	Regexes           []*regexp.Regexp
	ShowAllBranches   configdomain.AllBranches
	UnknownBranchType configdomain.UnknownBranchType
}

// layoutBranches adds entries for the given branch and its children to the given entry list.
// The entries are indented according to their position in the given lineage.
func layoutBranches(args layoutBranchesArgs) {
	if args.excludeBranches.Contains(args.branch) {
		return
	}
	if args.branchInfos.HasLocalBranch(args.branch) || args.showAllBranches.Enabled() {
		var otherWorktree bool
		if branchInfo, hasBranchInfo := args.branchInfos.FindByLocalName(args.branch).Get(); hasBranchInfo {
			otherWorktree = branchInfo.SyncStatus == gitdomain.SyncStatusOtherWorktree
		} else {
			otherWorktree = false
		}
		branchType, hasBranchType := args.branchesAndTypes[args.branch]
		if !hasBranchType && len(args.branchTypes) > 0 {
			branchType = args.unknownBranchType.BranchType()
		}
		var hasCorrectBranchType bool
		if len(args.branchTypes) == 0 || slices.Contains(args.branchTypes, branchType) {
			hasCorrectBranchType = true
		}
		matchesRegex := args.regexes.Matches(args.branch.String())
		if hasCorrectBranchType && matchesRegex {
			*args.result = append(*args.result, dialog.SwitchBranchEntry{
				Branch:        args.branch,
				Indentation:   args.indentation,
				OtherWorktree: otherWorktree,
				Type:          branchType,
			})
		}
	}
	for _, child := range args.lineage.Children(args.branch) {
		layoutBranches(layoutBranchesArgs{
			branch:            child,
			branchInfos:       args.branchInfos,
			branchTypes:       args.branchTypes,
			branchesAndTypes:  args.branchesAndTypes,
			excludeBranches:   args.excludeBranches,
			indentation:       args.indentation + "  ",
			lineage:           args.lineage,
			regexes:           args.regexes,
			result:            args.result,
			showAllBranches:   args.showAllBranches,
			unknownBranchType: args.unknownBranchType,
		})
	}
}

type layoutBranchesArgs struct {
	branch            gitdomain.LocalBranchName
	branchInfos       gitdomain.BranchInfos
	branchTypes       []configdomain.BranchType
	branchesAndTypes  configdomain.BranchesAndTypes
	excludeBranches   gitdomain.LocalBranchNames
	indentation       string
	lineage           configdomain.Lineage
	regexes           regexes.Regexes
	result            *dialog.SwitchBranchEntries
	showAllBranches   configdomain.AllBranches
	unknownBranchType configdomain.UnknownBranchType
}
