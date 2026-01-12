package state

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

func FilePath(repoDir gitdomain.RepoRootDir, homeDir configdomain.HomeDir, fileType FileType) (string, error) {
	sanitizedRepo := SanitizePath(repoDir)
	return homeDir.Join("git-town", sanitizedRepo, fileType.String()+".json"), nil
}
