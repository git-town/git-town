package statefile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/messages"
)

func FilePath(repoDir domain.RepoRootDir) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.RunstatePathProblem, err)
	}
	persistenceDir := filepath.Join(configDir, "git-town", "runstate")
	filename := SanitizePath(repoDir)
	return filepath.Join(persistenceDir, filename+".json"), err
}
