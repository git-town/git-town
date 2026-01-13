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

// UserConfigDir is the directory that contains the user-specific configuration on this machine,
// i.e. ~/.config.
type UserConfigDir string

func (self UserConfigDir) RepoConfigDir(repoDir gitdomain.RepoRootDir) RepoConfigDir {
	return RepoConfigDir(filepath.Join(self.String(), "git-town", SanitizePath(repoDir.String())))
}

func (self UserConfigDir) String() string {
	return string(self)
}

func SanitizePath[T ~string](dir T) T { //nolint:ireturn
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

// SystemUserConfigDir provides the UserConfigDir to use in production.
func SystemUserConfigDir() (UserConfigDir, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.ConfigDirUserCannotDetermine, err)
	}
	return UserConfigDir(configDir), nil
}
