package runstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/state"
)

// Save stores the given run state for the given Git repo to disk.
func Save(runState RunState, repoDir gitdomain.RepoRootDir) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.RunstateSerializeProblem, err)
	}
	persistencePath, err := state.FilePath(repoDir, state.FileTypeRunstate)
	if err != nil {
		return err
	}
	persistenceDir := filepath.Dir(persistencePath)
	if err = os.MkdirAll(persistenceDir, 0o700); err != nil {
		return err
	}
	if err = os.WriteFile(persistencePath, content, 0o600); err != nil {
		return fmt.Errorf(messages.FileWriteProblem, persistencePath, err)
	}
	return nil
}
