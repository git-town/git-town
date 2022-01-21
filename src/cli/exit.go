package cli

import "os"

// Exit prints the given error message and terminates the application.
func Exit(err error) {
	PrintError(err)
	os.Exit(1)
}
