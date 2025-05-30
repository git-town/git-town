package runlog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/state"
)

// Append writes the given entry to the end of the runlog for this repo.
// TODO:
// - fullinterpreter.Execute receives a receipt that the initial runlog was written
// - it writes the final runlog on exit
func Write(entry Entry, repoDir gitdomain.RepoRootDir) (WriteReceipt, error) {
	content, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return WriteReceipt{}, fmt.Errorf(messages.RunLogSerializeProblem, err)
	}
	persistencePath, err := state.FilePath(repoDir, state.FileTypeRunlog)
	if err != nil {
		return WriteReceipt{}, err
	}
	persistenceDir := filepath.Dir(persistencePath)
	err = os.MkdirAll(persistenceDir, 0o700)
	if err != nil {
		return WriteReceipt{}, err
	}
	file, err := os.OpenFile(persistencePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return WriteReceipt{}, fmt.Errorf(messages.RunLogCannotOpen, persistencePath, err)
	}
	_, err = file.Write(content)
	if err != nil {
		return WriteReceipt{}, fmt.Errorf(messages.RunLogCannotWrite, persistencePath, err)
	}
	return WriteReceipt{}, nil
}

type WriteReceipt struct{}
