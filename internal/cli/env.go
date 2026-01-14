package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/subshell"
)

// SystemUserConfigDir provides the UserConfigDir to use in production.
func SystemUserConfigDir() (configdomain.UserConfigDir, error) {
	if subshell.IsInTest() {
		home := os.Getenv("HOME")
		return configdomain.UserConfigDir(filepath.Join(home, ".config")), nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.ConfigDirUserCannotDetermine, err)
	}
	return configdomain.UserConfigDir(configDir), nil
}
