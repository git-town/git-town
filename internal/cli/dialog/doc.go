// Package dialog provides high-level screens through which the user can enter data into Git Town.
package dialog

import (
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type DialogArgs[T any] struct {
	Global Option[T]
	Inputs dialogcomponents.TestInputs
	Local  Option[T]
}
