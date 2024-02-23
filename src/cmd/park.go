package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cli/format"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/execute"
	"github.com/git-town/git-town/v12/src/git"
	"github.com/git-town/git-town/v12/src/gohacks"
	"github.com/git-town/git-town/v12/src/messages"
	"github.com/spf13/cobra"
)

const parkDesc = "Parks the current branch"

const parkHelp = `
Git Town does not sync parked branches until they are currently checked out.`

func parkCmd() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "park",
		Args:    cobra.NoArgs,
		GroupID: "types",
		Short:   parkDesc,
		Long:    cmdhelpers.Long(parkDesc, parkHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePark(args, readVerboseFlag(cmd))
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func executePark(args []string, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		DryRun:           false,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateGitRepo:  false,
		ValidateIsOnline: false,
		Verbose:          verbose,
	})
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return setOfflineStatus(args[0], repo.Runner)
	}
	displayOfflineStatus(repo.Runner)
	return nil
}

func displayOfflineStatus(run *git.ProdRunner) {
	fmt.Println(format.Bool(run.Config.FullConfig.Offline.Bool()))
}

func setOfflineStatus(text string, run *git.ProdRunner) error {
	value, err := gohacks.ParseBool(text)
	if err != nil {
		return fmt.Errorf(messages.ValueInvalid, gitconfig.KeyOffline, text)
	}
	return run.Config.SetOffline(configdomain.Offline(value))
}
