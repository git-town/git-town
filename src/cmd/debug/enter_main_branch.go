package debug

import (
	"fmt"

	"github.com/spf13/cobra"
)

func enterMainBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use: "enter-main-branch",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello")
		},
	}
}
