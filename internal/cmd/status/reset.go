package status

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/git-town/git-town/v17/internal/vm/statefile"
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
	err := statefile.Delete(".")
	if err != nil {
		return err
	}
	fmt.Println(messages.RunstateDeleted)
	return nil
}
