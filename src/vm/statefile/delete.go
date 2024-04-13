package statefile

import (
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

// Delete removes the stored run state from disk.
func Delete(repoDir gitdomain.RepoRootDir) error {
	filename, err := FilePath(repoDir)
	if err != nil {
		return err
	}
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf(messages.FileStatProblem, filename, err)
	}
	err = os.Remove(filename)
	if err != nil {
		return fmt.Errorf(messages.FileDeleteProblem, filename, err)
	}
	return nil
}
