package configdomain

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/messages"
)

// CachedAccess provides access to the local and global configuration data stored in Git metadata
// made efficient through an in-memory cache.
type CachedAccess struct {
	Access
	FullCache
}

// NewConfiguration provides a Configuration instance reflecting the configuration values in the given directory.
func NewCachedAccess(fullCache FullCache, runner Runner) CachedAccess {
	return CachedAccess{
		FullCache: fullCache,
		Access: Access{
			Runner: runner,
		},
	}
}

// Reload refreshes the cached configuration information.
func (self *CachedAccess) Reload() {
	self.FullCache, _ = LoadFullCache(&self.Access)
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (self *CachedAccess) RemoveLocalGitConfiguration() error {
	err := self.Run("git", "config", "--remove-section", "git-town")
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if exitErr.ExitCode() == 128 {
				// Git returns exit code 128 when trying to delete a non-existing config section.
				// This is not an error condition in this workflow so we can ignore it here.
				return nil
			}
		}
		return fmt.Errorf(messages.ConfigRemoveError, err)
	}
	for child := range *self.LocalConfig.Lineage {
		key := fmt.Sprintf("git-town-branch.%s.parent", child)
		err = self.Run("git", "config", "--unset", key)
		if err != nil {
			return fmt.Errorf(messages.ConfigRemoveError, err)
		}
	}
	return nil
}

// RemoveParent removes the parent branch entry for the given branch
// from the Git configuration.
func (self *CachedAccess) RemoveParent(branch gitdomain.LocalBranchName) {
	self.LocalConfig.Lineage.RemoveBranch(branch)
	// ignoring errors here because the entry might not exist
	_ = self.RemoveLocalConfigValue(NewParentKey(branch))
}

// SetGlobalConfigValue sets the given configuration setting in the global Git configuration.
func (self *CachedAccess) SetGlobalConfigValue(key Key, value string) error {
	self.GlobalCache[key] = value
	return self.Access.SetGlobalConfigValue(key, value)
}

// SetLocalConfigValue sets the local configuration with the given key to the given value.
func (self *CachedAccess) SetLocalConfigValue(key Key, value string) error {
	self.LocalCache[key] = value
	return self.Access.SetLocalConfigValue(key, value)
}
