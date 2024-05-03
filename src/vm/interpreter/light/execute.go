package light

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/vm/program"
	"github.com/git-town/git-town/v14/src/vm/shared"
)

func Execute(prog program.Program, lineage configdomain.Lineage) {
	for _, opcode := range prog {
		err := opcode.Run(shared.RunArgs{
			Connector:                       nil,
			DialogTestInputs:                nil,
			Lineage:                         lineage,
			PrependOpcodes:                  nil,
			RegisterUndoablePerennialCommit: nil,
			UpdateInitialBranchLocalSHA:     nil,
		})
		if err != nil {
			fmt.Println(colors.Red().Styled("NOTICE: " + err.Error()))
		}
	}
}
