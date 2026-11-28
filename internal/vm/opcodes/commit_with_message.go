package opcodes

import (
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// CommitWithMessage commits all open changes using the given commit message.
type CommitWithMessage struct {
	AuthorOverride  Option[gitdomain.Author]
	CleanupVerbatim bool
	CommitHook      configdomain.CommitHook
	Message         gitdomain.CommitMessage
}

func (self *CommitWithMessage) Run(args shared.RunArgs) error {
	useMessage := configdomain.UseCustomMessage(self.Message)
	if self.CleanupVerbatim {
		useMessage = configdomain.UseCustomMessageVerbatim(self.Message)
	}
	return args.Git.Commit(args.Frontend, useMessage, self.AuthorOverride, self.CommitHook)
}
