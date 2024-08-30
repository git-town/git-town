package validate

import (
	"errors"

	"github.com/git-town/git-town/v16/internal/messages"
)

func NoOpenChanges(hasOpenChanges bool) error {
	if hasOpenChanges {
		return errors.New(messages.ShipOpenChanges)
	}
	return nil
}
