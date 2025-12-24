package runlog

import (
	"os"
	"time"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// Entry is an entry in the runlog.
type Entry struct {
	Branches       map[gitdomain.BranchName]gitdomain.SHA // branches at this state
	Command        string                                 // the command through which the user called Git Town via the CLI
	Event          Event                                  // whether this event happens at the beginning or end of the Git Town command
	PendingCommand Option[string]                         // the currently pending Git Town command
	Time           time.Time                              // the time when this event happened
}

func NewEntry(event Event, branchInfos gitdomain.BranchInfos, pendingCommand Option[string]) Entry {
	branches := map[gitdomain.BranchName]gitdomain.SHA{}
	for _, branchInfo := range branchInfos {
		if hasLocalBranch, localName, localSHA := branchInfo.GetLocal(); hasLocalBranch {
			branches[localName.BranchName()] = localSHA
		}
		if hasRemoteBranch, remoteName, remoteSHA := branchInfo.GetRemote(); hasRemoteBranch {
			branches[remoteName.BranchName()] = remoteSHA
		}
	}
	return Entry{
		Branches:       branches,
		Command:        stringslice.JoinArgs(os.Args),
		Event:          event,
		PendingCommand: pendingCommand,
		Time:           time.Now(),
	}
}
