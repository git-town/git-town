package runstate

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Load loads the run state for the given Git repo from disk.
// Returns None if there is no saved runstate.
func Load(runstatePath RunstatePath) (Option[RunState], error) {
	_, err := os.Stat(runstatePath.String())
	if err != nil {
		if os.IsNotExist(err) {
			return None[RunState](), nil
		}
		return None[RunState](), fmt.Errorf(messages.FileStatProblem, runstatePath, err)
	}
	content, err := os.ReadFile(runstatePath.String())
	if err != nil {
		return None[RunState](), fmt.Errorf(messages.FileReadProblem, runstatePath, err)
	}
	var runState RunState
	if err = json.Unmarshal(content, &runState); err != nil {
		return None[RunState](), fmt.Errorf(messages.FileContentInvalidJSON, runstatePath, err)
	}
	return Some(runState), nil
}
