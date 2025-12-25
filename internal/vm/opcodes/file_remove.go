package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

type FileRemove struct {
	FilePath string
}

func (self *FileRemove) Run(args shared.RunArgs) error {
	return args.Git.RemoveFile(args.Frontend, self.FilePath)
}
