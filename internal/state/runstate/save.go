package runstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
)

// Save stores the given run state for the given Git repo to disk.
func Save(runState RunState, configDir configdomain.ConfigDir) error {
	content, err := json.MarshalIndent(runState, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.RunstateSerializeProblem, err)
	}
	persistencePath := configDir.RunstatePath()
	persistenceDir := filepath.Dir(persistencePath)
	if err = os.MkdirAll(persistenceDir, 0o700); err != nil {
		return err
	}
	if err = os.WriteFile(persistencePath, content, 0o600); err != nil {
		return fmt.Errorf(messages.FileWriteProblem, persistencePath, err)
	}
	return nil
}
