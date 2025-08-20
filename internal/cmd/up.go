package cmd

import (
	"cmp"
	"fmt"
	"os"
	"regexp"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/flags"
	"github.com/git-town/git-town/v21/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/spf13/cobra"
)

const (
	upShort = "Switch to the child branch"
	upLong  = `Moves "up" in the stack by switching to the child of the current branch.`
)

func upCmd() *cobra.Command {
	addMergeFlag, readMergeFlag := flags.Merge()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "up",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.NoArgs,
		Short:   upShort,
		Long:    cmdhelpers.Long(upShort, upLong),
		RunE: func(cmd *cobra.Command, _ []string) error {
			merge, errMerge := readMergeFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errMerge, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				Detached:    None[configdomain.Detached](),
				DryRun:      None[configdomain.DryRun](),
				Verbose:     verbose,
			})
			return executeUp(executeUpArgs{
				cliConfig: cliConfig,
				merge:     merge,
			})
		},
	}
	addMergeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeUpArgs struct {
	cliConfig configdomain.PartialConfig
	merge     configdomain.SwitchUsingMerge
}

func executeUp(args executeUpArgs) error {
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

	// Get the current branch
	currentBranch, err := repo.Git.CurrentBranch(repo.Backend)
	if err != nil {
		return err
	}

	// Get the child branches from lineage
	children := repo.UnvalidatedConfig.NormalConfig.Lineage.Children(currentBranch)
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
	data, exit, err := determineBranchData(repo)
	if err != nil || exit {
		return err
	}
	entries := dialog.NewSwitchBranchEntries(dialog.NewSwitchBranchEntriesArgs{
		BranchInfos:       data.branchInfos,
		BranchTypes:       []configdomain.BranchType{},
		BranchesAndTypes:  data.branchesAndTypes,
		ExcludeBranches:   gitdomain.LocalBranchNames{},
		Lineage:           repo.UnvalidatedConfig.NormalConfig.Lineage,
		MainBranch:        repo.UnvalidatedConfig.UnvalidatedConfig.MainBranch,
		Regexes:           []*regexp.Regexp{},
		ShowAllBranches:   false,
		UnknownBranchType: repo.UnvalidatedConfig.NormalConfig.UnknownBranchType,
	})
	fmt.Println()
	fmt.Print(branchLayout(entries, data))

	return nil
}
