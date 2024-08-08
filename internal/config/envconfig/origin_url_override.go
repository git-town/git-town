package envconfig

import (
	"os"

	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

func OriginURLOverride() Option[string] {
	override := os.Getenv("GIT_TOWN_REMOTE")
	if override == "" {
		return None[string]()
	}
	return Some(override)
}
