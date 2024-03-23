package statefile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/messages"
)

func FilePath(repoDir gitdomain.RepoRootDir) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(messages.RunstatePathProblem, err)
	}
	persistenceDir := filepath.Join(configDir, "git-town", "runstate")
	filename := SanitizePath(repoDir)
	return filepath.Join(persistenceDir, filename+".json"), err
}
