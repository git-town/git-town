package git

import (
	"errors"

	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/util"
)

// IsOffline returns whether Git Town is currently in offline mode
func IsOffline() bool {
	return util.StringToBool(getConfigurationValueWithDefault("git-town.offline", "false"))
}

// ValidateIsOnline asserts that Git Town is not in offline mode
func ValidateIsOnline() error {
	if IsOffline() {
		return errors.New("This command requires an active internet connection")
	}
	return nil
}

// ValidateIsRepository asserts that the current directory is in a repository
func ValidateIsRepository() error {
	if IsRepository() {
		return nil
	}
	return errors.New("This is not a Git repository")
}

// isRepository is cached in order to minimize the number of git commands run
var isRepository bool
var isRepositoryInitialized bool

// IsRepository returns whether or not the current directory is in a repository
func IsRepository() bool {
	if !isRepositoryInitialized {
		isRepository = command.New("git", "rev-parse").Err() == nil
		isRepositoryInitialized = true
	}
	return isRepository
}
