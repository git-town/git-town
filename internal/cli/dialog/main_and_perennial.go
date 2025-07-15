package dialog

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

func MainAndPerennials(args MainAndPerennialsArgs) (mainBranch gitdomain.LocalBranchName, perennials gitdomain.LocalBranchNames, exit dialogdomain.Exit, err error) {
	unvalidatedMain, hasMain := args.UnvalidatedMain.Get()
	if hasMain {
		return unvalidatedMain, args.UnvalidatedPerennials, false, nil
	}
	if args.HasConfigFile {
		return unvalidatedMain, args.UnvalidatedPerennials, false, errors.New(messages.ConfigMainbranchInConfigFile)
	}
	fmt.Print(messages.ConfigNeeded)
	mainBranch, exit, err = MainBranch(args.LocalBranches, args.GetDefaultBranch(args.Backend), args.DialogInputs)
	if err != nil || exit {
		return mainBranch, args.UnvalidatedPerennials, exit, err
	}
	perennials, exit, err = PerennialBranches(args.LocalBranches, args.UnvalidatedPerennials, mainBranch, args.DialogInputs)
	return mainBranch, perennials, exit, err
}

type MainAndPerennialsArgs struct {
	Backend               subshelldomain.RunnerQuerier
	DialogInputs          dialogcomponents.TestInputs
	GetDefaultBranch      func(subshelldomain.Querier) Option[gitdomain.LocalBranchName]
	HasConfigFile         bool
	LocalBranches         gitdomain.LocalBranchNames
	UnvalidatedMain       Option[gitdomain.LocalBranchName]
	UnvalidatedPerennials gitdomain.LocalBranchNames
}
