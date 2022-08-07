package cli

import (
	"fmt"

	"github.com/fatih/color"
)

// PrintDryRunMessage prints the dry-run message.
func PrintDryRunMessage() {
	_, err := color.New(color.FgBlue).Print(dryRunMessage)
	if err != nil {
		fmt.Print(dryRunMessage)
	}
}

const dryRunMessage = `
In dry run mode. No commands will be run. When run in normal mode, the command
output will appear beneath the command. Some commands will only be run if
necessary. For example: 'git push' will run if and only if there are local
commits not on origin.
`
