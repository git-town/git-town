package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

func FilePath(repoDir gitdomain.RepoRootDir, fileType FileType) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.RunstatePathProblem, err)
	}
	persistenceDir := filepath.Join(configDir, "git-town", fileType.String())
	filePath := SanitizePath(repoDir)
	return filepath.Join(persistenceDir, filePath+".json"), err
}
