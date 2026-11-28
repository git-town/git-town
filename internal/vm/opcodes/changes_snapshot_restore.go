package opcodes

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
)

// ChangesSnapshotRestore restores a previously created snapshot of the exact index and worktree state.
type ChangesSnapshotRestore struct {
	SHA gitdomain.SHA
}

func (self *ChangesSnapshotRestore) Run(args shared.RunArgs) error {
	return args.Git.RestoreOpenChangesSnapshot(args.Backend, args.Frontend, self.SHA)
}
