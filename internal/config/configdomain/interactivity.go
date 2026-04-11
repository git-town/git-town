package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// InteractiveEnabled is a sentinal value that indicates that interactive mode is enabled.
const InteractiveEnabled = Interactive("")

// Interactive indicates whether interactive features are enabled.
// If it is an empty string, interactive mode is enabled.
// If it is not an empty string, interactive mode is disabled and the string contains the reason why.
type Interactive string

// Check indicates via an error if interactive mode is enabled.
// No error: interactive mode is enabled.
// Error: interactive mode is disabled.
func (self Interactive) Check() error {
	if self == "" {
		return nil
	}
	return &InteractivityError{Reason: string(self)}
}

func (self Interactive) String() string {
	if err := self.Check(); err != nil {
		return "disabled: " + string(self)
	}
	return "enabled"
}

func NewInteractiveFromEnv(envTerm string) Option[Interactive] {
	if strings.ToLower(envTerm) == "dumb" {
		return Some(Interactive("only a dumb terminal available"))
	}
	return None[Interactive]()
}

func NewInteractiveFromTTY(tty HasTTY) Option[Interactive] {
	if tty {
		return None[Interactive]()
	}
	return Some(Interactive("no interactive terminal available"))
}

// ------------------------------------------------------------------------------

type InteractivityError struct {
	Reason string
}

func (self InteractivityError) Error() string {
	return self.Reason
}
