package runlog

import (
	"os"
	"strings"
	"time"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Entry is an entry in the runlog.
type Entry struct {
	Branches       map[gitdomain.BranchName]gitdomain.SHA // branches at this state
	Command        string                                 // the command through which the user called Git Town via the CLI
	Event          Event
	PendingCommand Option[string] // the currently pending Git Town command
	Time           time.Time
}

func NewEntryFromBranchInfos(branchInfos gitdomain.BranchInfos, pendingCommand Option[string]) Entry {
	branches := map[gitdomain.BranchName]gitdomain.SHA{}
	for _, branchInfo := range branchInfos {
		if hasLocalBranch, localName, localSHA := branchInfo.GetLocal(); hasLocalBranch {
			branches[localName.BranchName()] = localSHA
		}
		if hasRemoteBranch, remoteName, remoteSHA := branchInfo.GetRemoteBranch(); hasRemoteBranch {
			branches[remoteName.BranchName()] = remoteSHA
		}
	}
	return Entry{
		Branches:       branches,
		Command:        strings.Join(os.Args, " "),
		Event:          EventStart,
		PendingCommand: pendingCommand,
		Time:           time.Now(),
	}
}
