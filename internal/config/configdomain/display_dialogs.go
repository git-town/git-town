package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type DisplayDialogs bool

func LoadDisplayDialogsFromEnv(envTerm string) Option[DisplayDialogs] {
	if strings.ToLower(envTerm) == "dumb" {
		return Some(DisplayDialogs(false))
	}
	return None[DisplayDialogs]()
}
