package configdomain

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

// ConfigDirUser is the directory that contains the user configuration on this machine,
// i.e. ~/.config.
type ConfigDirUser string

func (self ConfigDirUser) RepoConfigDir(repoDir gitdomain.RepoRootDir) ConfigDirRepo {
	return ConfigDirRepo(filepath.Join(self.String(), "git-town", SanitizePath(repoDir.String())))
}

func (self ConfigDirUser) String() string {
	return string(self)
}

// SystemUserConfigDir provides the UserConfigDir to use in production.
func SystemUserConfigDir() (ConfigDirUser, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.RunstateCannotDetermineUserDir, err)
	}
	return ConfigDirUser(configDir), nil
}

func SanitizePath[T ~string](dir T) T {
	replaceCharacterRE := regexp.MustCompile("[[:^alnum:]]")
	sanitized := replaceCharacterRE.ReplaceAllString(string(dir), "-")
	sanitized = strings.ToLower(sanitized)
	replaceDoubleMinusRE := regexp.MustCompile("--+") // two or more dashes
	sanitized = replaceDoubleMinusRE.ReplaceAllString(sanitized, "-")
	for strings.HasPrefix(sanitized, "-") {
		sanitized = sanitized[1:]
	}
	return T(sanitized)
}
