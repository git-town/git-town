package envconfig

import (
	"os"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// TODO: reveive the pre-loaded env vars as an argument
func RemoteURLOverride() Option[string] {
	return NewOptionIfExists(os.LookupEnv("GIT_TOWN_REMOTE"))
}
