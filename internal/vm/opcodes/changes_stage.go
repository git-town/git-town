package opcodes

import "github.com/git-town/git-town/v21/internal/vm/shared"

type ChangesStage struct{}

func (self *ChangesStage) Run(args shared.RunArgs) error {
	return args.Git.StageFiles(args.Frontend, "-A")
}
