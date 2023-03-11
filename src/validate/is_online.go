package validate

import (
	"errors"

	"github.com/git-town/git-town/v7/src/git"
)

// IsOnline is a validationCondition that verifies that the given Git repository is online.
func IsOnline(repo *git.ProdRepo) error {
	isOffline, err := repo.Config.IsOffline()
	if err != nil {
		return err
	}
	if isOffline {
		return errors.New("this command requires an active internet connection")
	}
	return nil
}
