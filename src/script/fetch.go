package script

import "fmt"

// Fetch gets the local Git repo in sync with origin
// without modifying the workspace.
func Fetch() error {
	err := RunCommand("git", "fetch", "--prune", "--tags")
	if err != nil {
		return fmt.Errorf("cannot fetch updates from origin: %w", err)
	}
	return nil
}
