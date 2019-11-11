package script

import "github.com/pkg/errors"

// Fetch gets the local Git repo in sync with origin
// without modifying the workspace.
func Fetch() error {
	err := RunCommand("git", "fetch", "--prune", "--tags")
	if err != nil {
		return errors.Wrap(err, "cannot fetch updates from origin")
	}
	return nil
}
