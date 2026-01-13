package runlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Append writes the given entry to the end of the runlog for this repo.
func Write(event Event, branchInfos gitdomain.BranchInfos, pendingCommand Option[string], runlogPath FilePath) error {
	entry := NewEntry(event, branchInfos, pendingCommand)
	content, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf(messages.RunLogSerializeProblem, err)
	}
	content = append(content, []byte("\n\n")...)
	persistenceDir := filepath.Dir(runlogPath.String())
	if err = os.MkdirAll(persistenceDir, 0o700); err != nil {
		return err
	}
	file, err := os.OpenFile(runlogPath.String(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf(messages.RunLogCannotOpen, runlogPath, err)
	}
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf(messages.RunLogCannotWrite, runlogPath, err)
	}
	return nil
}
