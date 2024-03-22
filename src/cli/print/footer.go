package print

import (
	"fmt"

	"github.com/git-town/git-town/v12/src/messages"
)

func Footer(verbose bool, commandsCount int, finalMessages []string) {
	fmt.Println()
	if verbose {
		fmt.Printf(messages.CommandsRun, commandsCount)
	}
	FinalMessages(finalMessages)
}

// NoFinalMessages can be used by callers of PrintFooter to indicate
// that the command has no final messages to print.
var NoFinalMessages = []string{} //nolint:gochecknoglobals

func FinalMessages(messages []string) {
	for _, message := range messages {
		fmt.Println("\n" + message)
	}
}
