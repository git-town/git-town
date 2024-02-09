package validate

import (
	"errors"

	"github.com/git-town/git-town/v12/src/messages"
)

func NoOpenChanges(hasOpenChanges bool) error {
	if hasOpenChanges {
		return errors.New(messages.ShipOpenChanges)
	}
	return nil
}
