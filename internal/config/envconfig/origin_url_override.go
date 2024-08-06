package envconfig

import "os"

func OriginURLOverride() string {
	return os.Getenv("GIT_TOWN_REMOTE")
}
