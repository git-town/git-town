package git

import (
	"errors"

	"github.com/Originate/git-town/src/util"
)

// ValidateIsRepository asserts that the current directory is in a repository
func ValidateIsRepository() error {
	if IsRepository() {
		return nil
	}
	return errors.New("This is not a Git repository")
}

// IsRepository returns whether or not the current directory is in a repository
func IsRepository() bool {
	_, err := util.GetFullCommandOutput("git", "rev-parse")
	return err == nil
}
