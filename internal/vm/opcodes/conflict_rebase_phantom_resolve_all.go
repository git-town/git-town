package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/vm/shared"
)

type ConflictRebasePhantomResolveAll struct {
	BranchToRebaseOnto      gitdomain.BranchName
	CurrentBranch           gitdomain.LocalBranchName
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *ConflictRebasePhantomResolveAll) Run(args shared.RunArgs) error {
	quickInfos, err := args.Git.FileConflictQuickInfos(args.Backend)
	if err != nil {
		return err
	}
	fmt.Println("111111111111111111111111111 QUICK INFOS")
	fmt.Println(quickInfos)
	rootBranch := args.Config.Value.NormalConfig.Lineage.Root(self.CurrentBranch)
	fullInfos, err := args.Git.FileConflictFullInfos(args.Backend, quickInfos, self.BranchToRebaseOnto.Location(), rootBranch)
	if err != nil {
		return err
	}
	fmt.Println("111111111111111111111111111 FULL INFOS")
	fmt.Println(fullInfos)
	phantomConflicts := git.DetectPhantomRebaseConflicts(fullInfos, self.BranchToRebaseOnto, rootBranch)
	newOpcodes := []shared.Opcode{}
	for _, phantomConflict := range phantomConflicts {
		newOpcodes = append(newOpcodes, &ConflictPhantomResolve{
			FilePath:   phantomConflict.FilePath,
			Resolution: gitdomain.ConflictResolutionTheirs,
		})
	}
	newOpcodes = append(newOpcodes, &ConflictRebasePhantomFinalize{})
	args.PrependOpcodes(newOpcodes...)
	return nil
}
