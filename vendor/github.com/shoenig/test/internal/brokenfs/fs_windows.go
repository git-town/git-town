// Copyright (c) The Test Authors
// SPDX-License-Identifier: MPL-2.0

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
