package opcodes

import "github.com/git-town/git-town/v22/internal/vm/shared"

type FileStage struct {
	FilePath string
}

func (self *FileStage) Run(args shared.RunArgs) error {
	return args.Git.StageFiles(args.Frontend, self.FilePath)
}
