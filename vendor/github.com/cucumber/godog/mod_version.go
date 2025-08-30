//go:build go1.12
// +build go1.12

package godog

import (
	"runtime/debug"
)

func init() {
	if info, available := debug.ReadBuildInfo(); available {
		if Version == "v0.0.0-dev" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}
}
