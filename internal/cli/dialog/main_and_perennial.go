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
	unvalidatedMain, hasMain := args.UnvalidatedMain.Get()
	if hasMain {
		return unvalidatedMain, gitdomain.LocalBranchNames{}, false, nil
	}
	fmt.Print(messages.ConfigNeeded)
	mainBranchOpt, exit, err := MainBranch(MainBranchArgs{
		GitStandardBranch:     args.GetDefaultBranch(args.Backend),
		UnscopedGitMainBranch: args.UnvalidatedConfig.NormalConfig.Git.MainBranch,
		LocalGitMainBranch:    args.UnvalidatedConfig.GitLocal.MainBranch,
		LocalBranches:         args.LocalBranches,
		Inputs:                args.DialogInputs.Next(),
	})
	if err != nil || exit {
		return "", gitdomain.LocalBranchNames{}, exit, err
	}
	mainBranch = mainBranchOpt.GetOrElse(args.UnvalidatedMain.GetOrPanic())
	perennials, exit, err = PerennialBranches(PerennialBranchesArgs{
		LocalBranches:         args.LocalBranches,
		MainBranch:            mainBranch,
		UnscopedGitPerennials: args.UnvalidatedConfig.NormalConfig.Git.PerennialBranches,
		LocalGitPerennials:    args.UnvalidatedConfig.GitLocal.PerennialBranches,
		Inputs:                args.DialogInputs.Next(),
	})
	return mainBranch, perennials, exit, err
}

type MainAndPerennialsArgs struct {
	Backend           subshelldomain.RunnerQuerier
	DialogInputs      dialogcomponents.TestInputs
	GetDefaultBranch  func(subshelldomain.Querier) Option[gitdomain.LocalBranchName]
	LocalBranches     gitdomain.LocalBranchNames
	UnvalidatedMain   Option[gitdomain.LocalBranchName]
	UnvalidatedConfig config.UnvalidatedConfig
}
