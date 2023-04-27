package commands

import (
	"fmt"
)

// AddRemote adds a Git remote with the given name and URL to this repository.
func AddRemote(cmds TestCommands, name, url string) error {
	_, err := cmds.Run("git", "remote", "add", name, url)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w", name, url, err)
	}
	cmds.Config.RemotesCache.Invalidate()
	return nil
}
