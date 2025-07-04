package gitconfig

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/git-town/git-town/v21/internal/cli/colors"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// IO provides low-level access to the Git configuration on disk.
type IO struct {
	Shell subshelldomain.RunnerQuerier
}

func (self *IO) RemoteURL(remote gitdomain.Remote) Option[string] {
	output, err := self.Shell.Query("git", "remote", "get-url", remote.String())
	if err != nil {
		// NOTE: it's okay to ignore the error here.
		// If we get an error here, we simply don't use the origin remote.
		return None[string]()
	}
	return NewOption(strings.TrimSpace(output))
}

func (self *IO) RemoveConfigValue(scope configdomain.ConfigScope, key configdomain.Key) error {
	args := []string{"config"}
	if scope == configdomain.ConfigScopeGlobal {
		args = append(args, "--global")
	}
	args = append(args, "--unset", key.String())
	return self.Shell.Run("git", args...)
}

// RemoveLocalGitConfiguration removes all Git Town configuration.
func (self *IO) RemoveLocalGitConfiguration(localSnapshot configdomain.SingleSnapshot) error {
	if err := self.Shell.Run("git", "config", "--remove-section", "git-town"); err != nil {
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
	for key := range localSnapshot {
		if strings.HasPrefix(key.String(), "git-town-branch.") {
			if err := self.Shell.Run("git", "config", "--unset", key.String()); err != nil {
				return fmt.Errorf(messages.ConfigRemoveError, err)
			}
		}
	}
	return nil
}

// SetConfigValue sets the given configuration setting in the global Git configuration.
func (self *IO) SetConfigValue(scope configdomain.ConfigScope, key configdomain.Key, value string) error {
	args := []string{"config"}
	if scope == configdomain.ConfigScopeGlobal {
		args = append(args, "--global")
	}
	args = append(args, key.String(), value)
	return self.Shell.Run("git", args...)
}

// updates a custom Git alias (not set up by Git Town)
func (self *IO) UpdateExternalGitTownAlias(scope configdomain.ConfigScope, key configdomain.Key, oldValue, newValue string) {
	fmt.Println(colors.Cyan().Styled(fmt.Sprintf(messages.SettingDeprecatedValueMessage, scope, key, oldValue, newValue)))
	if err := self.SetConfigValue(scope, key, newValue); err != nil {
		fmt.Printf(messages.SettingCannotWrite, scope, key, err)
	}
}
