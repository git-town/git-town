package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func MainAndPerennials(args MainAndPerennialsArgs) (mainBranch gitdomain.LocalBranchName, perennials gitdomain.LocalBranchNames, exit dialogdomain.Exit, err error) {
	if unvalidatedMain, hasMain := args.UnvalidatedConfig.UnvalidatedConfig.MainBranch.Get(); hasMain {
		return unvalidatedMain, args.UnvalidatedConfig.NormalConfig.PerennialBranches, false, nil
	}
	fmt.Print(messages.ConfigNeeded)
	mainBranch, exit, err = MainBranch(args.LocalBranches, args.GetDefaultBranch(args.Backend), args.DialogInputs)
	if err != nil || exit {
		return mainBranch, args.UnvalidatedConfig.NormalConfig.PerennialBranches, exit, err
	}
	perennials, exit, err = PerennialBranches(args.LocalBranches, args.UnvalidatedConfig.NormalConfig.PerennialBranches, mainBranch, args.DialogInputs)
	return mainBranch, perennials, exit, err
}

type MainAndPerennialsArgs struct {
	Backend           subshelldomain.RunnerQuerier
	DialogInputs      dialogcomponents.TestInputs
	GetDefaultBranch  func(subshelldomain.Querier) Option[gitdomain.LocalBranchName]
	LocalBranches     gitdomain.LocalBranchNames
	UnvalidatedConfig config.UnvalidatedConfig
}
