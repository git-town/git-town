package cmd

import (
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

const upDesc = "Move one position up in the current stack"

func upCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "up",
		GroupID: cmdhelpers.GroupIDStack,
		Args:    cobra.NoArgs,
		Short:   upDesc,
		Long:    cmdhelpers.Long(upDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			cliConfig := cliconfig.New(cliconfig.NewArgs{
				AutoResolve: None[configdomain.AutoResolve](),
				DryRun:      None[configdomain.DryRun](),
				Verbose:     verbose,
			})
			return executeUp(cliConfig)
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executeUp(cliConfig configdomain.PartialConfig) error {
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

	// Get the current branch
	currentBranch, err := repo.Git.CurrentBranch(repo.Backend)
	if err != nil {
		return err
	}

	// Get the parent branch from lineage
	parent, hasParent := repo.UnvalidatedConfig.NormalConfig.Lineage.Parent(currentBranch).Get()
	if !hasParent {
		return fmt.Errorf(messages.UpNoParent, currentBranch)
	}

	// Check out the parent branch
	err = repo.Git.CheckoutBranch(repo.Frontend, parent, false)
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
	fmt.Print(branchLayout(entries, data))

	return nil
}
