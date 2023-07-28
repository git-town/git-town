package validate

import (
	"fmt"

	"github.com/git-town/git-town/v9/src/messages"
)

func NoOpenChanges(hasOpenChanges bool) error {
	if hasOpenChanges {
		return fmt.Errorf(messages.ShipOpenChanges)
	}
	return nil
}
