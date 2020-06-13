package git

import (
	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/util"
)

// Remotes are cached in order to minimize the number of git commands run.
var remotes []string
var remotesInitialized bool

// HasRemote returns whether the current repository contains a Git remote
// with the given name.
func HasRemote(name string) bool {
	if !remotesInitialized {
		remotes = command.MustRun("git", "remote").OutputLines()
		remotesInitialized = true
	}
	return util.DoesStringArrayContain(remotes, name)
}
