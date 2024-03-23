package statefile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/messages"
	"github.com/git-town/git-town/v13/src/vm/runstate"
)

// Load loads the run state for the given Git repo from disk. Can return nil if there is no saved runstate.
func Load(repoDir gitdomain.RepoRootDir) (*runstate.RunState, error) {
	filename, err := FilePath(repoDir)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil
		}
		return nil, fmt.Errorf(messages.FileStatProblem, filename, err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf(messages.FileReadProblem, filename, err)
	}
	var runState runstate.RunState
	err = json.Unmarshal(content, &runState)
	if err != nil {
		return nil, fmt.Errorf(messages.FileContentInvalidJSON, filename, err)
	}
	return &runState, nil
}
