package cmd

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"

	"github.com/git-town/git-town/v21/internal/cli/dialog"
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
	downShort = "Switch to the parent branch"
	downLong  = `Moves "down" in the stack by switching to the parent of the current branch.`
)

func downCmd() *cobra.Command {
	addMergeFlag, readMergeFlag := flags.Merge()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "down",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.NoArgs,
		Short:   downShort,
		Long:    cmdhelpers.Long(downShort, downLong),
		RunE: func(cmd *cobra.Command, _ []string) error {
			merge, errMerge := readMergeFlag(cmd)
			verbose, errVerbose := readVerboseFlag(cmd)
			if err := cmp.Or(errMerge, errVerbose); err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve:  None[configdomain.AutoResolve](),
				AutoSync:     None[configdomain.AutoSync](),
				Detached:     Some(configdomain.Detached(true)),
				DryRun:       None[configdomain.DryRun](),
				PushBranches: None[configdomain.PushBranches](),
				Stash:        None[configdomain.Stash](),
				Verbose:      verbose,
			})
			return executeDown(executeDownArgs{
				cliConfig: cliConfig,
				merge:     merge,
			})
		},
	}
	addMergeFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

type executeDownArgs struct {
	cliConfig configdomain.PartialConfig
	merge     configdomain.SwitchUsingMerge
}

func executeDown(args executeDownArgs) error {
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
	currentBranchOpt, err := repo.Git.CurrentBranch(repo.Backend)
	if err != nil {
		return err
	}
	currentBranch, hasCurrentBranch := currentBranchOpt.Get()
	if !hasCurrentBranch {
		return errors.New(messages.DownNoCurrentBranch)
	}

	// Get the parent branch from lineage
	parent, hasParent := repo.UnvalidatedConfig.NormalConfig.Lineage.Parent(currentBranch).Get()
	if !hasParent {
		return fmt.Errorf(messages.DownNoParent, currentBranch)
	}

	// Check out the parent branch
	err = repo.Git.CheckoutBranch(repo.Frontend, parent, args.merge)
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
