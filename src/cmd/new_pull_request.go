package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v12/src/cli/flags"
	"github.com/git-town/git-town/v12/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v12/src/config/gitconfig"
	"github.com/git-town/git-town/v12/src/messages"
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
	fmt.Println(messages.PullRequestDeprecation)
	time.Sleep(2000 * time.Millisecond)
}
