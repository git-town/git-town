package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

type ExitToShellIfUncommittedChanges struct{}

func (self *ExitToShellIfUncommittedChanges) Run(args shared.RunArgs) error {
	hasUncommittedFiles, err := args.Git.HasUncommittedFiles(args.Backend)
	if err != nil {
		return err
	}
	if hasUncommittedFiles {
		args.PrependOpcodes(&ExitToShell{})
	}
	return nil
}
