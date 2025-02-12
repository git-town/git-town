package status

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v18/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v18/internal/git/gitdomain"
	"github.com/git-town/git-town/v18/internal/messages"
	"github.com/git-town/git-town/v18/internal/vm/statefile"
	"github.com/spf13/cobra"
)

const statusResetDesc = "Resets the current suspended Git Town command"

func resetRunstateCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "reset",
		Args:  cobra.NoArgs,
		Short: statusResetDesc,
		Long:  cmdhelpers.Long(statusResetDesc),
		RunE: func(_ *cobra.Command, _ []string) error {
			return executeStatusReset()
		},
	}
	return &cmd
}

func executeStatusReset() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = statefile.Delete(gitdomain.RepoRootDir(cwd))
	if err != nil {
		return err
	}
	fmt.Println(messages.RunstateDeleted)
	return nil
}
