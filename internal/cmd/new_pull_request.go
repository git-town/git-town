package cmd

import (
	"fmt"
	"time"

	"github.com/git-town/git-town/v17/internal/cli/flags"
	"github.com/git-town/git-town/v17/internal/cmd/cmdhelpers"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	"github.com/spf13/cobra"
)

func newPullRequestCommand() *cobra.Command {
	addBodyFlag, readBodyFlag := flags.ProposalBody("b")
	addBodyFileFlag, readBodyFileFlag := flags.ProposalBodyFile()
	addDetachedFlag, readDetachedFlag := flags.Detached()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	addTitleFlag, readTitleFlag := flags.ProposalTitle()
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	cmd := cobra.Command{
		Use:     "new-pull-request",
		GroupID: "basic",
		Hidden:  true,
		Args:    cobra.NoArgs,
		Short:   proposeDesc,
		Long:    cmdhelpers.Long(proposeDesc, fmt.Sprintf(proposeHelp, configdomain.KeyHostingPlatform, configdomain.KeyHostingOriginHostname)),
		RunE: func(cmd *cobra.Command, _ []string) error {
			printDeprecationNotice()
			detached, err := readDetachedFlag(cmd)
			if err != nil {
				return err
			}
			dryRun, err := readDryRunFlag(cmd)
			if err != nil {
				return err
			}
			bodyFile, err := readBodyFileFlag(cmd)
			if err != nil {
				return err
			}
			bodyText, err := readBodyFlag(cmd)
			if err != nil {
				return err
			}
			title, err := readTitleFlag(cmd)
			if err != nil {
				return err
			}
			verbose, err := readVerboseFlag(cmd)
			if err != nil {
				return err
			}
			result := executePropose(detached, dryRun, verbose, title, bodyText, bodyFile)
			printDeprecationNotice()
			return result
		},
	}
	addBodyFlag(&cmd)
	addBodyFileFlag(&cmd)
	addDetachedFlag(&cmd)
	addDryRunFlag(&cmd)
	addTitleFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func printDeprecationNotice() {
	fmt.Println(messages.PullRequestDeprecation)
	time.Sleep(2000 * time.Millisecond)
}
