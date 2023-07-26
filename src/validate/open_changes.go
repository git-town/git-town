package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/git"
)

func NoOpenChanges(backend git.BackendCommands) error {
	hasOpenChanges, err := backend.HasOpenChanges()
	if err != nil {
		return err
	}
	if hasOpenChanges {
		return fmt.Errorf("you have uncommitted changes. Did you mean to commit them before shipping?")
	}
	return nil
}
