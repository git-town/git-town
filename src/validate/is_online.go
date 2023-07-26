package validate

import (
	"errors"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
)

// IsOnline verifies that the given Git repository is online.
func IsOnline(config *git.RepoConfig) error {
	isOffline, err := config.IsOffline()
	if err != nil {
		return err
	}
	if isOffline {
		return errors.New(messages.OfflineNotAllowed)
	}
	return nil
}
