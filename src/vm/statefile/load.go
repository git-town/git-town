package statefile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
	"github.com/git-town/git-town/v14/src/vm/runstate"
)

// Load loads the run state for the given Git repo from disk.
// Returns None if there is no saved runstate.
func Load(repoDir gitdomain.RepoRootDir) (Option[runstate.RunState], error) {
	filename, err := FilePath(repoDir)
	if err != nil {
		return None[runstate.RunState](), err
	}
	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return None[runstate.RunState](), nil
		}
		return None[runstate.RunState](), fmt.Errorf(messages.FileStatProblem, filename, err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return None[runstate.RunState](), fmt.Errorf(messages.FileReadProblem, filename, err)
	}
	var runState runstate.RunState
	err = json.Unmarshal(content, &runState)
	if err != nil {
		return None[runstate.RunState](), fmt.Errorf(messages.FileContentInvalidJSON, filename, err)
	}
	return Some(runState), nil
}
