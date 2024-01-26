package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/gitconfig"
	"github.com/spf13/cobra"
)

func newPullRequestCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Hidden:  true,
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, gitconfig.KeyHostingPlatform, gitconfig.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, args []string) error {
			printDeprecationNotice()
			result := executePropose(readDryRunFlag(cmd), readVerboseFlag(cmd))
			printDeprecationNotice()
			return result
		},
	}
	addDryRunFlag(&cmd)
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
