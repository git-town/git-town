package lightinterpreter

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v23/internal/config"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/git"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks"
	"github.com/git-town/git-town/v23/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v23/internal/messages"
	"github.com/git-town/git-town/v23/internal/subshell/subshelldomain"
	"github.com/git-town/git-town/v23/internal/vm/program"
	"github.com/git-town/git-town/v23/internal/vm/shared"
	"github.com/git-town/git-town/v23/pkg/colors"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

func Execute(args ExecuteArgs) {
	for {
		nextStep, hasNextStep := args.Prog.Pop().Get()
		if !hasNextStep {
			return
		}
		runnable, isRunnable := nextStep.(shared.Runnable)
		if !isRunnable {
			panic(fmt.Errorf(messages.OpcodeNotRunnable, gohacks.TypeName(nextStep)))
		}
		err := runnable.Run(shared.RunArgs{
			Backend:                         args.Backend,
			BranchInfos:                     args.BranchInfos,
			Config:                          NewMutable(&args.Config),
			Connector:                       args.Connector,
			DetectedForgeType:               args.DetectedForgeType,
			FinalMessages:                   args.FinalMessages,
			Frontend:                        args.Frontend,
			Git:                             args.Git,
			Inputs:                          dialogcomponents.NewInputs(),
			PrependOpcodes:                  args.Prog.Prepend,
			RegisterUndoablePerennialCommit: nil,
			UpdateInitialSnapshotLocalSHA:   nil,
		})
		if err != nil {
			fmt.Println(colors.Red().Styled("NOTICE: " + err.Error()))
		}
	}
}

type ExecuteArgs struct {
	Backend           subshelldomain.RunnerQuerier
	BranchInfos       gitdomain.BranchInfos
	Config            config.ValidatedConfig
	Connector         Option[forgedomain.Connector]
	DetectedForgeType Option[forgedomain.DetectedForgeType]
	FinalMessages     stringslice.Collector
	Frontend          subshelldomain.Runner
	Git               git.Commands
	Prog              program.Program
}
