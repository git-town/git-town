package opcodes

import "github.com/git-town/git-town/v21/internal/vm/shared"

type FileStage struct {
	FilePath                string
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *FileStage) Run(args shared.RunArgs) error {
	return args.Git.StageFiles(args.Frontend, self.FilePath)
}
