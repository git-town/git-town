package script

import (
	"github.com/git-town/git-town/src/dryrun"
	"github.com/git-town/git-town/src/git"

	"github.com/fatih/color"
)

var dryRunMessage = `
In dry run mode. No commands will be run. When run in normal mode, the command
output will appear beneath the command. Some commands will only be run if
necessary. For example: 'git push' will run if and only if there are local
commits not on the remote.
`

// ActivateDryRun causes all commands to not be run.
func ActivateDryRun() {
	_, err := color.New(color.FgBlue).Print(dryRunMessage)
	if err != nil {
		panic(err)
	}
	dryrun.Activate(git.GetCurrentBranchName())
}
