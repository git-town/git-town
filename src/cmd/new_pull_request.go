package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v14/src/cli/flags"
	"github.com/git-town/git-town/v14/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/spf13/cobra"
)

func newPullRequestCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addTitleFlag, readTitleFlag := flags.ProposalTitle()
	addBodyFlag, readBodyFlag := flags.ProposalBody()
	addBodyFileFlag, readBodyFileFlag := flags.ProposalBodyFile()
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Hidden:  true,
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, configdomain.KeyHostingPlatform, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			printDeprecationNotice()
			result := executePropose(readDryRunFlag(cmd), readVerboseFlag(cmd), readTitleFlag(cmd), readBodyFlag(cmd), readBodyFileFlag(cmd))
			printDeprecationNotice()
			return result
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	addTitleFlag(&cmd)
	addBodyFlag(&cmd)
	addBodyFileFlag(&cmd)
	return &cmd
}

func printDeprecationNotice() {
	fmt.Println(messages.PullRequestDeprecation)
	time.Sleep(2000 * time.Millisecond)
}
