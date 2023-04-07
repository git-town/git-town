package validate

import (
	"errors"

	"github.com/git-town/git-town/v8/src/git"
)

// IsOnline verifies that the given Git repository is online.
func IsOnline(config *git.RepoConfig) error {
	isOffline, err := config.IsOffline()
	if err != nil {
		return err
	}
	if isOffline {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
