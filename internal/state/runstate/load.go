package runstate

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Load loads the run state for the given Git repo from disk.
// Returns None if there is no saved runstate.
func Load(configDir configdomain.ConfigDirRepo) (Option[RunState], error) {
	filename := configDir.RunstatePath()
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return None[RunState](), nil
		}
		return None[RunState](), fmt.Errorf(messages.FileStatProblem, filename, err)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return None[RunState](), fmt.Errorf(messages.FileReadProblem, filename, err)
	}
	var runState RunState
	if err = json.Unmarshal(content, &runState); err != nil {
		return None[RunState](), fmt.Errorf(messages.FileContentInvalidJSON, filename, err)
	}
	return Some(runState), nil
}
