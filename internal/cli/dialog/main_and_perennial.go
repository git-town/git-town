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
	fmt.Print(messages.ConfigNeeded)
	_, mainBranch, exit, err = MainBranch(MainBranchArgs{
		GitStandardBranch:     args.GetDefaultBranch(args.Backend),
		Inputs:                args.DialogInputs,
		LocalBranches:         args.LocalBranches,
		LocalGitMainBranch:    args.UnvalidatedConfig.GitLocal.MainBranch,
		UnscopedGitMainBranch: args.UnvalidatedConfig.GitUnscoped.MainBranch,
	})
	if err != nil || exit {
		return "", gitdomain.LocalBranchNames{}, exit, err
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
