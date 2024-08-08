package envconfig

import (
	"os"

	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

func OriginURLOverride() Option[string] {
	return NewOption(os.Getenv("GIT_TOWN_REMOTE"))
}
