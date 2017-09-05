package script

import "github.com/Originate/git-town/src/exit"

// Fetch gets the local Git repo in sync with origin,
// without modifying the workspace.
func Fetch() {
	err := RunCommand("git", "fetch", "--prune")
	exit.On(err)
}
