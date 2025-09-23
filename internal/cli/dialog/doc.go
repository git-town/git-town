// Package dialog provides high-level screens through which the user can enter data into Git Town.
package dialog

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Args are arguments for dialogs that allow to enter a textual configuration value.
type Args[T any] struct {
	Global Option[T]
	Inputs dialogcomponents.Inputs
	Local  Option[T]
}
