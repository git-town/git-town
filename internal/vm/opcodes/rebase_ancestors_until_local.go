package opcodes

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type RebaseAncestorsUntilLocal struct {
	Branch gitdomain.LocalBranchName
	// SHA of the direct parent at the previous run.
	// These are the commits we need to remove from this branch.
	CommitsToRemove Option[gitdomain.SHA]
}

func (self *RebaseAncestorsUntilLocal) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	branch := self.Branch
	detached := args.Config.Value.NormalConfig.Detached.ShouldWorkDetached()
	for {
		ancestor, hasAncestor := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasAncestor {
			break
		}
		ancestorIsPerennial := args.Config.Value.IsMainOrPerennialBranch(ancestor)
		if ancestorIsPerennial && detached {
			break
		}
		ancestorInfo, hasAncestorInfo := args.BranchInfos.FindLocalOrRemote(ancestor).Get()
		if !hasAncestorInfo {
			branch = ancestor
			continue
		}
		if localAncestor, ancestorIsLocal := ancestorInfo.Local.Get(); ancestorIsLocal {
			program = append(program, &RebaseAncestorLocal{
				Ancestor:        localAncestor.Name,
				Branch:          self.Branch,
				CommitsToRemove: self.CommitsToRemove,
			})
			break
		}
		// the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
		ancestorTracking, ancestorIsRemote := ancestorInfo.RemoteName.Get()
		if !ancestorIsRemote {
			return errors.New(messages.BranchInfoNoContent)
		}
		program = append(program, &RebaseAncestorRemote{
			Ancestor: ancestorTracking,
			Branch:   self.Branch,
		})
		branch = ancestor
	}
	args.PrependOpcodes(program...)
	return nil
}
