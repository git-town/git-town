package configdomain

import (
	"fmt"
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// InteractivityEnabled is a sentinal value for DisplayDialogs and indicates that dialogs can be displayed.
const InteractivityEnabled = Interactivity("")

// Interactivity indicates whether interactive features are enabled.
// If it is an empty string, interactivity is enabled.
// If it is not an empty string, the string contains the reason why interactivity is disabled.
type Interactivity string

// Check indicates via an error if interactivity is enabled.
func (self Interactivity) Check() error {
	if self == "" {
		return nil
	}
	return &InteractivityError{Reason: self.String()}
}

func (self Interactivity) String() string {
	if err := self.Check(); err != nil {
		return fmt.Sprintf("disabled: %s", err.Error())
	}
	return "enabled"
}

func NewInteractivityFromEnv(envTerm string) Option[Interactivity] {
	if strings.ToLower(envTerm) == "dumb" {
		return Some(Interactivity("only a dumb terminal available"))
	}
	return None[Interactivity]()
}

func NewInteractivityFromTTY(tty HasTTY) Option[Interactivity] {
	if tty {
		return None[Interactivity]()
	}
	return Some(Interactivity("no interactive terminal available"))
}

// ------------------------------------------------------------------------------

type InteractivityError struct {
	Reason string
}

func (self InteractivityError) Error() string {
	return self.Reason
}
