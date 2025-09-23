package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type SnapshotInitialUpdateLocalSHA struct {
	Branch gitdomain.LocalBranchName
	SHA    gitdomain.SHA
}

func (self *SnapshotInitialUpdateLocalSHA) Run(args shared.RunArgs) error {
	return args.UpdateInitialSnapshotLocalSHA(self.Branch, self.SHA)
}
