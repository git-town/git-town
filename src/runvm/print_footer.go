package runvm

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

func PrintFooter(debug bool, commandsCount int, finalMessages []string) {
	fmt.Println()
	if debug {
		fmt.Printf(messages.CommandsRun, commandsCount)
	}
	for _, message := range finalMessages {
		fmt.Println("\n" + message)
	}
}

// NoFinalMessages indicates that commands have no final messages to print.
var NoFinalMessages = []string{} //nolint:gochecknoglobals
