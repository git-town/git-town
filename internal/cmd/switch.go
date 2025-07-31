package cmd

import (
	"cmp"
	"errors"
	"os"
	"os/exec"
	"regexp"

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
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: false,
				DryRun:      None[configdomain.DryRun](),
				Verbose:     verbose,
			})
			return executeSwitch(executeSwitchArgs{
				allBranches:  allBranches,
				argv:         args,
				branchTypes:  branchTypes,
				cliConfig:    cliConfig,
				displayTypes: displayTypes,
				merge:        merge,
			})
		},
	}
	addAllFlag(&cmd)
	addDisplayTypesFlag(&cmd)
	addMergeFlag(&cmd)
	addTypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeSwitchArgs struct {
	allBranches  configdomain.AllBranches
	argv         []string
	branchTypes  []configdomain.BranchType
	cliConfig    configdomain.PartialConfig
	displayTypes configdomain.DisplayTypes
	merge        configdomain.SwitchUsingMerge
}

func executeSwitch(args executeSwitchArgs) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, exit, err := determineSwitchData(args.argv, repo)
	if err != nil || exit {
		return err
	}
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(data.branchNames)
	unknownBranchType := repo.UnvalidatedConfig.NormalConfig.UnknownBranchType
	entries := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
		BranchInfos:       data.branchesSnapshot.Branches,
		BranchTypes:       args.branchTypes,
		BranchesAndTypes:  branchesAndTypes,
		ExcludeBranches:   gitdomain.LocalBranchNames{},
		Lineage:           data.config.NormalConfig.Lineage,
		MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
		Regexes:           data.regexes,
		ShowAllBranches:   args.allBranches,
		UnknownBranchType: unknownBranchType,
	})
	if len(entries) == 0 {
		return errors.New(messages.SwitchNoBranches)
	}
	cursor := entries.IndexOf(data.initialBranch)
	branchToCheckout, exit, err := dialog.SwitchBranch(dialog.SwitchBranchArgs{
		CurrentBranch:      Some(data.initialBranch),
		Cursor:             cursor,
		DisplayBranchTypes: args.displayTypes,
		Entries:            entries,
		InputName:          "switch-branch",
		Inputs:             data.inputs,
		Title:              None[string](),
		UncommittedChanges: data.uncommittedChanges,
	})
	if err != nil || exit {
		return err
	}
	if branchToCheckout == data.initialBranch {
		return nil
	}
	err = repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, args.merge)
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

func determineSwitchData(args []string, repo execute.OpenRepoResult) (data switchData, exit dialogdomain.Exit, err error) {
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
