package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/vm/shared"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type RebaseAncestorsUntilLocal struct {
	Branch                  gitdomain.LocalBranchName
	CommitsToRemove         Option[gitdomain.Location]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *RebaseAncestorsUntilLocal) Run(args shared.RunArgs) error {
	fmt.Println("333333333333333333333333333333333333 RebaseAncestorsUntilLocal for", self.Branch)
	program := []shared.Opcode{}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	branch := self.Branch
	for {
		parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		parentIsPerennial := args.Config.Value.IsMainOrPerennialBranch(parent)
		if args.Detached.IsTrue() && parentIsPerennial {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if !parentIsLocal {
			// here the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
			program = append(program, &RebaseAncestorRemote{
				Ancestor: parent.AtRemote(args.Config.Value.NormalConfig.DevRemote),
				Branch:   self.Branch,
			})
			branch = parent
			continue
		}
		// here we found a local parent
		program = append(program, &RebaseAncestorLocal{
			Ancestor:        parent,
			Branch:          self.Branch,
			CommitsToRemove: self.CommitsToRemove,
		})
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
