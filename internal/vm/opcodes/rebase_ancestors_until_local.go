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
		ancestor, hasAncestor := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasAncestor {
			break
		}
		ancestorIsPerennial := args.Config.Value.IsMainOrPerennialBranch(ancestor)
		if ancestorIsPerennial && args.Detached.IsTrue() {
			break
		}
		ancestorIsLocal := branchInfos.HasLocalBranch(ancestor)
		if !ancestorIsLocal {
			// here the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
			program = append(program, &RebaseAncestorRemote{
				Ancestor: ancestor.AtRemote(args.Config.Value.NormalConfig.DevRemote),
				Branch:   self.Branch,
			})
			branch = ancestor
			continue
		}
		// here we found a local parent
		program = append(program, &RebaseAncestorLocal{
			Ancestor:        ancestor,
			Branch:          self.Branch,
			CommitsToRemove: self.CommitsToRemove,
		})
		break
	}
	args.PrependOpcodes(program...)
	return nil
}
