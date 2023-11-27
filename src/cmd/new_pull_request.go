package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v10/src/cli/flags"
	"github.com/git-town/git-town/v10/src/config"
	"github.com/spf13/cobra"
)

func newPullRequestCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    long(proposeDesc, fmt.Sprintf(proposeHelp, config.KeyCodeHostingDriver, config.KeyCodeHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			printDeprecationNotice()
			result := executePropose(readVerboseFlag(cmd))
			printDeprecationNotice()
			return result
		},
	}
	addVerboseFlag(&cmd)
	return &cmd
}

func printDeprecationNotice() {
	fmt.Println("DEPRECATION NOTICE")
	fmt.Println("")
	fmt.Println("This command has been renamed to \"git town propose\"")
	fmt.Println("and will be removed in future versions of Git Town.")
	time.Sleep(2000 * time.Millisecond)

}
