package light

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

func Execute(args ExecuteArgs) {
	for _, opcode := range args.Prog {
		err := opcode.Run(shared.RunArgs{
			Backend:                         args.Backend,
			Config:                          &args.Config,
			Connector:                       nil,
			DialogTestInputs:                nil,
			FinalMessages:                   args.FinalMessages,
			Frontend:                        args.Frontend,
			Lineage:                         args.Lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			fmt.Println(colors.Red().Styled("NOTICE: " + err.Error()))
		}
	}
}

type ExecuteArgs struct {
	Backend       git.BackendCommands
	Config        config.Config
	FinalMessages *stringslice.Collector
	Frontend      git.FrontendCommands
	Lineage       configdomain.Lineage
	Prog          program.Program
}
