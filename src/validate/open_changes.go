package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
)

func NoOpenChanges(backend git.BackendCommands) error {
	hasOpenChanges, err := backend.HasOpenChanges()
	if err != nil {
		return err
	}
	if hasOpenChanges {
		return fmt.Errorf(messages.ShipOpenChanges)
	}
	return nil
}
