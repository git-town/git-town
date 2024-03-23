// Git Town - a high-level CLI for Git
//
// Git Town adds Git commands that make software development more efficient
// by keeping Git branches better in sync with each other.
// This reduces merge conflicts and the number of Git commands you need to run.
package main

import (
	"os"
	"runtime/debug"

	"github.com/git-town/git-town/v13/src/cli/print"
	"github.com/git-town/git-town/v13/src/cmd"
)

func main() {
	debug.SetGCPercent(-1)
	err := cmd.Execute()
	if err != nil {
		print.Error(err)
		os.Exit(1)
	}
}
