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

const ShowDisplayDialogs = DisplayDialogs("")

// DisplayDialogs indicates whether dialogs can be displayed.
// If it is an empty string, dialogs can be displayed.
// If it is not an empty string, the string contains the reason why dialogs cannot be displayed.
type DisplayDialogs string

func (self DisplayDialogs) String() string {
	return string(self)
}

// Verify indicates via an error if dialogs cannot be displayed.
func (self DisplayDialogs) Verify() error {
	if self == "" {
		return nil
	}
	return &CannotDisplayDialogsError{Reason: self.String()}
}

func NewDisplayDialogsFromEnv(envTerm string) Option[DisplayDialogs] {
	if strings.ToLower(envTerm) == "dumb" {
		return Some(DisplayDialogs("dumb terminal configured"))
	}
	return None[DisplayDialogs]()
}

func NewDisplayDialogsFromTTY(tty bool) Option[DisplayDialogs] {
	if tty {
		return None[DisplayDialogs]()
	}
	return Some(DisplayDialogs("no TTY detected"))
}
