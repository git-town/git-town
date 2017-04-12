package script

import "log"

// Fetch gets the local Git repo in sync with origin,
// without modifying the workspace.
func Fetch() {
	err := RunCommand("git", "fetch", "--prune")
	if err != nil {
		log.Fatal(err)
	}
}
