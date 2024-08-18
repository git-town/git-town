package light

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/colors"
	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/config"
	"github.com/git-town/git-town/v15/internal/git"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/vm/program"
	"github.com/git-town/git-town/v15/internal/vm/shared"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

func Execute(args ExecuteArgs) {
	for _, opcode := range args.Prog {
		err := opcode.Run(shared.RunArgs{
			Backend:                         args.Backend,
			Config:                          args.Config,
			Connector:                       None[hostingdomain.Connector](),
			DialogTestInputs:                components.NewTestInputs(),
			FinalMessages:                   args.FinalMessages,
			Frontend:                        args.Frontend,
			Git:                             args.Git,
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
	Backend       gitdomain.RunnerQuerier
	Config        config.ValidatedConfig
	FinalMessages stringslice.Collector
	Frontend      gitdomain.Runner
	Git           git.Commands
	Prog          program.Program
}
