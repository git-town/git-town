package envconfig

import (
	"os"

	. "github.com/git-town/git-town/v15/pkg/prelude"
)

func RemoteURLOverride() Option[string] {
	return NewOption(os.Getenv("GIT_TOWN_REMOTE"))
}
