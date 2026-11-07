package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

type GhPrompt bool

func (self GhPrompt) Enabled() bool {
	return bool(self)
}

func (self GhPrompt) String() string {
	return strconv.FormatBool(bool(self))
}

func ParseGhPromptFromConfigFile(value stringss.Trimmed, _ string) (Option[GhPrompt], error) {
	if value == "" {
		return None[GhPrompt](), nil
	}
	return Some(GhPrompt(false)), nil
}

func ParseGhPromptOpt(value stringss.Trimmed, _ string) (Option[GhPrompt], error) {
	if value == "" {
		return None[GhPrompt](), nil
	}
	return Some(GhPrompt(false)), nil
}
