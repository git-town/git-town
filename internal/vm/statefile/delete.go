package statefile

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
)

// Delete removes the stored run state from disk.
func Delete(repoDir gitdomain.RepoRootDir) (existed bool, err error) {
	filename, err := FilePath(repoDir)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf(messages.FileStatProblem, filename, err)
	}
	err = os.Remove(filename)
	if err != nil {
		return false, fmt.Errorf(messages.FileDeleteProblem, filename, err)
	}
	return true, nil
}
