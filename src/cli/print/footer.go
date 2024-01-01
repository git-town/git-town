package print

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/messages"
)

func Footer(verbose bool, commandsCount int, finalMessages []string) {
	fmt.Println()
	if verbose {
		fmt.Printf(messages.CommandsRun, commandsCount)
	}
	for _, message := range finalMessages {
		fmt.Println("\n" + message)
	}
}

// NoFinalMessages can be used by callers of PrintFooter to indicate
// that the command has no final messages to print.
var NoFinalMessages = []string{} //nolint:gochecknoglobals
