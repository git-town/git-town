package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

func FilePath(repoDir gitdomain.RepoRootDir, fileType FileType) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.RunstateCannotDetermineUserDir, err)
	}
	sanitizedRepo := SanitizePath(repoDir)
	return filepath.Join(configDir, "git-town", sanitizedRepo, fileType.String()+".json"), nil
}
