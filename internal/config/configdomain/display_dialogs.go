package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type CannotDisplayDialogsError struct {
	Reason string
}

func (self CannotDisplayDialogsError) Error() string {
	return self.Reason
}

// ShowDisplayDialogs is a sentinal value for DisplayDialogs and indicates that dialogs can be displayed.
const ShowDisplayDialogs = DisplayDialogs("")

// DisplayDialogs indicates whether dialogs can be displayed.
// If it is an empty string, dialogs can be displayed.
// If it is not an empty string, the string contains the reason why dialogs cannot be displayed.
type DisplayDialogs string

// Check indicates via an error if dialogs cannot be displayed.
func (self DisplayDialogs) Check() error {
	if self == "" {
		return nil
	}
	return &CannotDisplayDialogsError{Reason: self.String()}
}

func (self DisplayDialogs) String() string {
	return string(self)
}

func NewDisplayDialogsFromEnv(envTerm string) Option[DisplayDialogs] {
	if strings.ToLower(envTerm) == "dumb" {
		return Some(DisplayDialogs("only a dumb terminal available"))
	}
	return None[DisplayDialogs]()
}

func NewDisplayDialogsFromTTY(tty HasTTY) Option[DisplayDialogs] {
	if tty {
		return None[DisplayDialogs]()
	}
	return Some(DisplayDialogs("no interactive terminal available"))
}
