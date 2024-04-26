package statefile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/runstate"
)

// Save stores the given run state for the given Git repo to disk.
func Save(runState runstate.RunState, repoDir gitdomain.RepoRootDir) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.RunstateSerializeProblem, err)
	}
	persistencePath, err := FilePath(repoDir)
	if err != nil {
		return err
	}
	persistenceDir := filepath.Dir(persistencePath)
	err = os.MkdirAll(persistenceDir, 0o700)
	if err != nil {
		return err
	}
	err = os.WriteFile(persistencePath, content, 0o600)
	if err != nil {
		return fmt.Errorf(messages.FileWriteProblem, persistencePath, err)
	}
	return nil
}
