// Code generated via scripts/generate.sh. DO NOT EDIT.

//go:build windows

package must

import (
	"os"
)

var (
	fsRoot = os.Getenv("HOMEDRIVE")
)

func init() {
	if fsRoot == "" {
		fsRoot = "C:"
	}
}
