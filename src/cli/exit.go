package cli

import "os"

// Exit prints the given error message and terminates the application.
func Exit(messages ...interface{}) {
	PrintError(messages...)
	os.Exit(1)
}
