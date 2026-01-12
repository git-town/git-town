package runstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/messages"
)

// Save stores the given run state for the given Git repo to disk.
func Save(runState RunState, runstatePath RunstatePath) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.RunstateSerializeProblem, err)
	}
	persistenceDir := filepath.Dir(runstatePath.String())
	if err = os.MkdirAll(persistenceDir, 0o700); err != nil {
		return err
	}
	if err = os.WriteFile(runstatePath.String(), content, 0o600); err != nil {
		return fmt.Errorf(messages.FileWriteProblem, runstatePath, err)
	}
	return nil
}
