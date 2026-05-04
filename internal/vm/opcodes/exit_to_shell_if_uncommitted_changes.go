package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/vm/shared"
)

type ExitToShellIfUncommittedChanges struct{}

func (self *ExitToShellIfUncommittedChanges) Run(args shared.RunArgs) error {
	uncommittedFiles, err := args.Git.UncommittedFiles(args.Backend)
	if err != nil {
		return err
	}
	fmt.Println("1111111111111111111111111111111111111111111", uncommittedFiles)
	if len(uncommittedFiles) > 0 {
		args.PrependOpcodes(&ExitToShell{})
	}
	return nil
}
