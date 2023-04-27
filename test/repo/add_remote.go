package repo

import (
	"fmt"
)

// AddRemote adds a Git remote with the given name and URL to the given repository clone.
func AddRemote(repo *Repo, name, url string) error {
	_, err := repo.Run("git", "remote", "add", name, url)
	if err != nil {
		return fmt.Errorf("cannot add remote %q --> %q: %w", name, url, err)
	}
	repo.Config().RemotesCache.Invalidate()
	return nil
}
