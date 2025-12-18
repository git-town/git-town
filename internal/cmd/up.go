package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v22/internal/cli/dialog"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/flags"
	"github.com/git-town/git-town/v22/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v22/internal/config/cliconfig"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/execute"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	upShort = "Switch to the child branch"
	upLong  = `Moves "up" in the stack by switching to the child of the current branch.`
)

func upCmd() *cobra.Command {
	addDisplayTypesFlag, readDisplayTypesFlag := flags.Displaytypes()
	addMergeFlag, readMergeFlag := flags.Merge()
	addOrderFlag, readOrderFlag := flags.Order()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "up",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.NoArgs,
		Short:   upShort,
		Long:    cmdhelpers.Long(upShort, upLong),
		RunE: func(cmd *cobra.Command, _ []string) error {
			displayTypes, errDisplayTypes := readDisplayTypesFlag(cmd)
			merge, errMerge := readMergeFlag(cmd)
			order, errOrder := readOrderFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errDisplayTypes, errMerge, errOrder, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:       None[configdomain.AutoResolve](),
				AutoSync:          None[configdomain.AutoSync](),
				Detached:          None[configdomain.Detached](),
				DisplayTypes:      displayTypes,
				DryRun:            None[configdomain.DryRun](),
				IgnoreUncommitted: None[configdomain.IgnoreUncommitted](),
				Order:             order,
				PushBranches:      None[configdomain.PushBranches](),
				Stash:             None[configdomain.Stash](),
				Verbose:           verbose,
			})
			return executeUp(executeUpArgs{
				cliConfig: cliConfig,
				merge:     merge,
			})
		},
	}
	addDisplayTypesFlag(&cmd)
	addMergeFlag(&cmd)
	addOrderFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeUpArgs struct {
	cliConfig configdomain.PartialConfig
	merge     configdomain.SwitchUsingMerge
}

func executeUp(args executeUpArgs) error {
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

	// Get the current branch
	currentBranchOpt, err := repo.Git.CurrentBranch(repo.Backend)
	if err != nil {
		return err
	}
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if !hasCurrentBranch {
		return errors.New(messages.UpNoCurrentBranch)
	}

	// Get the child branches from lineage
	children := repo.UnvalidatedConfig.NormalConfig.Lineage.Children(currentBranch, repo.UnvalidatedConfig.NormalConfig.Order)
	var child gitdomain.LocalBranchName
	switch len(children) {
	case 0:
		return fmt.Errorf(messages.UpNoChild, currentBranch)
	case 1:
		child = children[0]
	default:
		// more than one child --> let the user choose
		inputs := dialogcomponents.LoadInputs(os.Environ())
		selectedChild, exit, err := dialog.ChildBranch(dialog.ChildBranchArgs{
			ChildBranches: children,
			Inputs:        inputs,
		})
		if err != nil || exit {
			return err
		}
		child = selectedChild
	}

	// check out the child branch
	err = repo.Git.CheckoutBranch(repo.Frontend, child, args.merge)
	if err != nil {
		return err
	}

	// Display the branch hierarchy
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
	fmt.Println()
	fmt.Print(branchLayout(entries, data, repo.UnvalidatedConfig.NormalConfig.DisplayTypes))

	return nil
}
