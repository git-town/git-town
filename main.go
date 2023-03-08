// Git Town - a high-level CLI for Git
//
// Git Town adds Git commands that make software development more efficient
// by keeping Git branches better in sync with each other.
// This reduces merge conflicts and the number of Git commands you need to run.
package main

import (
	"os"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/cmd"
	"github.com/git-town/git-town/v7/src/dialog"
)

func main() {
	dialog.Initialize()
	err := cmd.Execute()
	if err != nil {
		cli.PrintError(err)
		os.Exit(1)
	}
}
