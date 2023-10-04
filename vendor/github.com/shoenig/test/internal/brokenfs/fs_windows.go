//go:build windows

package brokenfs

import (
	"os"
)

var (
	Root = os.Getenv("HOMEDRIVE")
)

func init() {
	if Root == "" {
		Root = "C:"
	}
}
