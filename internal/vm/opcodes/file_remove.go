package opcodes

import "github.com/git-town/git-town/v21/internal/vm/shared"

type FileRemove struct {
	FilePath                string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *FileRemove) Run(args shared.RunArgs) error {
	return args.Git.RemoveFile(args.Frontend, self.FilePath)
}
