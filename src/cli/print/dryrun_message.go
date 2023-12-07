package print

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/messages"
)

// DryRunMessage prints the dry-run message.
func DryRunMessage() {
	_, err := color.New(color.FgBlue).Print(messages.DryRun)
	if err != nil {
		fmt.Print(messages.DryRun)
	}
}
