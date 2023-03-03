package cli

import "os"

// Exit prints the given error message and terminates the application.
// TODO: delete this file.
func Exit(err error) {
	PrintError(err)
	os.Exit(1)
}
