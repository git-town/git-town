package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/git-town/git-town/v16/internal/cli/dialog"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/flags"
	"github.com/git-town/git-town/v16/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/execute"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/validate"
	"github.com/spf13/cobra"
)

const switchDesc = "Display the local branches visually and allows switching between them"

func switchCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addMergeFlag, readMergeFlag := flags.SwitchMerge()
	addTypeFlag, readTypeFlag := flags.BranchType()
	cmd := cobra.Command{
		Use:     "switch",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   switchDesc,
		Long:    cmdhelpers.Long(switchDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			branchTypes, err := readTypeFlag(cmd)
			if err != nil {
				return err
			}
			return executeSwitch(readVerboseFlag(cmd), readMergeFlag(cmd), branchTypes)
		},
	}
	addMergeFlag(&cmd)
	addVerboseFlag(&cmd)
	addTypeFlag(&cmd)
	return &cmd
}

func executeSwitch(verbose configdomain.Verbose, merge configdomain.SwitchUsingMerge, branchTypes []configdomain.BranchType) error {
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
	data, exit, err := determineSwitchData(repo, verbose)
	if err != nil || exit {
		return err
	}
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(data.branchNames)
	branchToCheckout, exit, err := dialog.SwitchBranch(data.branchNames, branchTypes, branchesAndTypes, data.initialBranch, data.config.Config.Lineage, data.branchesSnapshot.Branches, data.uncommittedChanges, data.dialogInputs.Next())
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
	uncommittedChanges bool
}

func determineSwitchData(repo execute.OpenRepoResult, verbose configdomain.Verbose) (data switchData, exit bool, err error) {
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
	localBranches := branchesSnapshot.Branches.LocalBranches().Names()
	branchesAndTypes := repo.UnvalidatedConfig.Config.Value.BranchesAndTypes(branchesSnapshot.Branches.LocalBranches().Names())
	validatedConfig, exit, err := validate.Config(validate.ConfigArgs{
		Backend:            repo.Backend,
		BranchesAndTypes:   branchesAndTypes,
		BranchesSnapshot:   branchesSnapshot,
		BranchesToValidate: localBranches,
		DialogTestInputs:   dialogTestInputs,
		Frontend:           repo.Frontend,
		Git:                repo.Git,
		LocalBranches:      localBranches,
		RepoStatus:         repoStatus,
		TestInputs:         dialogTestInputs,
		Unvalidated:        repo.UnvalidatedConfig,
	})
	if err != nil || exit {
		return data, exit, err
	}
	return switchData{
		branchNames:        branchesSnapshot.Branches.Names(),
		branchesSnapshot:   branchesSnapshot,
		config:             validatedConfig,
		dialogInputs:       dialogTestInputs,
		initialBranch:      initialBranch,
		lineage:            validatedConfig.Config.Lineage,
		uncommittedChanges: repoStatus.UntrackedChanges,
	}, false, err
}
