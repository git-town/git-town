package cmd

import (
	"cmp"
	"errors"
	"os"
	"os/exec"
	"regexp"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/regexes"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const switchDesc = "Display the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addAllFlag, readAllFlag := flags.All("list both remote-tracking and local branches")
	addDisplayTypesFlag, readDisplayTypesFlag := flags.Displaytypes()
	addMergeFlag, readMergeFlag := flags.Merge()
	addOrderFlag, readOrderFlag := flags.Order()
	addStashFlag, readStashFlag := flags.Stash()
	addTypeFlag, readTypeFlag := flags.BranchType()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: cmdhelpers.GroupIDBasic,
		Args:    cobra.ArbitraryArgs,
		Short:   switchDesc,
		Long:    cmdhelpers.Long(switchDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchTypes, errBranchTypes := readTypeFlag(cmd)
			allBranches, errAllBranches := readAllFlag(cmd)
			displayTypes, errDisplayTypes := readDisplayTypesFlag(cmd)
			merge, errMerge := readMergeFlag(cmd)
			order, errOrder := readOrderFlag(cmd)
			stash, errStash := readStashFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errBranchTypes, errAllBranches, errDisplayTypes, errMerge, errOrder, errStash, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:  None[configdomain.AutoResolve](),
				AutoSync:     None[configdomain.AutoSync](),
				Detached:     Some(configdomain.Detached(true)),
				DisplayTypes: displayTypes,
				DryRun:       None[configdomain.DryRun](),
				Order:        order,
				PushBranches: None[configdomain.PushBranches](),
				Stash:        stash,
				Verbose:      verbose,
			})
			return executeSwitch(executeSwitchArgs{
				allBranches: allBranches,
				argv:        args,
				branchTypes: branchTypes,
				cliConfig:   cliConfig,
				merge:       merge,
			})
		},
	}
	addAllFlag(&cmd)
	addDisplayTypesFlag(&cmd)
	addMergeFlag(&cmd)
	addOrderFlag(&cmd)
	addStashFlag(&cmd)
	addTypeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeSwitchArgs struct {
	allBranches configdomain.AllBranches
	argv        []string
	branchTypes []configdomain.BranchType
	cliConfig   configdomain.PartialConfig
	merge       configdomain.SwitchUsingMerge
}

func executeSwitch(args executeSwitchArgs) error {
Start:
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        args.cliConfig,
		IgnoreUnknown:    false,
		PrintBranchNames: true,
		PrintCommands:    true,
		ValidateGitRepo:  true,
		ValidateIsOnline: false,
	})
	if err != nil {
		return err
	}
	data, flow, err := determineSwitchData(args.argv, repo)
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
	branchesAndTypes := repo.UnvalidatedConfig.UnvalidatedBranchesAndTypes(data.branchesSnapshot.Branches.NamesAllBranches())
	unknownBranchType := repo.UnvalidatedConfig.NormalConfig.UnknownBranchType
	entriesArgs := dialog.NewSwitchBranchEntriesArgs{
		BranchInfos:       data.branchesSnapshot.Branches,
		BranchTypes:       args.branchTypes,
		BranchesAndTypes:  branchesAndTypes,
		ExcludeBranches:   gitdomain.LocalBranchNames{},
		Lineage:           data.config.NormalConfig.Lineage,
		MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
		Order:             data.config.NormalConfig.Order,
		Regexes:           data.regexes,
		ShowAllBranches:   false,
		UnknownBranchType: unknownBranchType,
	}
	entriesLocal := dialog.NewSwitchBranchEntries(entriesArgs)
	entriesArgs.ShowAllBranches = true
	entriesAll := dialog.NewSwitchBranchEntries(entriesArgs)
	if args.allBranches && len(entriesAll) == 0 {
		return errors.New(messages.SwitchNoBranches)
	}
	if !args.allBranches && len(entriesLocal) == 0 {
		return errors.New(messages.SwitchNoBranches)
	}
	cursor := entriesLocal.IndexOf(data.initialBranch)
	branchToCheckout, exit, err := dialog.SwitchBranch(dialog.SwitchBranchArgs{
		CurrentBranch:      Some(data.initialBranch),
		Cursor:             cursor,
		DisplayBranchTypes: repo.UnvalidatedConfig.NormalConfig.DisplayTypes,
		EntryData: dialog.EntryData{
			EntriesAll:      entriesAll,
			EntriesLocal:    entriesLocal,
			ShowAllBranches: args.allBranches,
		},
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
	if err := performSwitch(branchToCheckout, args.cliConfig.Stash.GetOr(false), data.hasOpenChanges, args.merge, repo); err != nil {
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
	branchesSnapshot   gitdomain.BranchesSnapshot
	config             config.UnvalidatedConfig
	hasOpenChanges     bool
	initialBranch      gitdomain.LocalBranchName
	inputs             dialogcomponents.Inputs
	lineage            configdomain.Lineage
	regexes            []*regexp.Regexp
	uncommittedChanges bool
}

func performSwitch(branchToCheckout gitdomain.LocalBranchName, stash configdomain.Stash, hasOpenChanges bool, merge configdomain.SwitchUsingMerge, repo execute.OpenRepoResult) error {
	if stash.ShouldStash() && hasOpenChanges {
		return switchWithStash(branchToCheckout, merge, repo)
	}
	return repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, merge)
}

func switchWithStash(branchToCheckout gitdomain.LocalBranchName, merge configdomain.SwitchUsingMerge, repo execute.OpenRepoResult) error {
	if err := repo.Git.Stash(repo.Frontend); err != nil {
		return err
	}
	errCheckout := repo.Git.CheckoutBranch(repo.Frontend, branchToCheckout, merge)
	errPop := repo.Git.PopStash(repo.Frontend)
	return cmp.Or(errCheckout, errPop)
}

func determineSwitchData(args []string, repo execute.OpenRepoResult) (data switchData, flow configdomain.ProgramFlow, err error) {
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
		HandleUnfinishedState: true,
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
	initialBranch, hasInitialBranch := branchesSnapshot.Active.Get()
	if !hasInitialBranch {
		return data, configdomain.ProgramFlowExit, errors.New(messages.CurrentBranchCannotDetermine)
	}
	regexes, err := regexes.NewRegexes(args)
	if err != nil {
		return data, configdomain.ProgramFlowExit, err
	}
	return switchData{
		branchesSnapshot:   branchesSnapshot,
		config:             repo.UnvalidatedConfig,
		hasOpenChanges:     repoStatus.OpenChanges,
		initialBranch:      initialBranch,
		inputs:             inputs,
		lineage:            repo.UnvalidatedConfig.NormalConfig.Lineage,
		regexes:            regexes,
		uncommittedChanges: repoStatus.OpenChanges,
	}, configdomain.ProgramFlowContinue, err
}
