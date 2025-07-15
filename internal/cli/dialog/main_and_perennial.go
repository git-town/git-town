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
		Inputs:                args.DialogInputs.Next(),
		LocalBranches:         args.LocalBranches,
		LocalGitMainBranch:    args.UnvalidatedConfig.GitLocal.MainBranch,
		UnscopedGitMainBranch: args.UnvalidatedConfig.NormalConfig.Git.MainBranch,
	})
	if err != nil || exit {
		return "", gitdomain.LocalBranchNames{}, exit, err
	}
	mainBranch = mainBranchOpt.GetOrElse(args.UnvalidatedMain.GetOrPanic())
	perennials, exit, err = PerennialBranches(PerennialBranchesArgs{
		Inputs:                args.DialogInputs.Next(),
		LocalBranches:         args.LocalBranches,
		LocalGitPerennials:    args.UnvalidatedConfig.GitLocal.PerennialBranches,
		MainBranch:            mainBranch,
		UnscopedGitPerennials: args.UnvalidatedConfig.NormalConfig.Git.PerennialBranches,
	})
	return mainBranch, perennials, exit, err
}

type MainAndPerennialsArgs struct {
	Backend           subshelldomain.RunnerQuerier
	DialogInputs      dialogcomponents.TestInputs
	GetDefaultBranch  func(subshelldomain.Querier) Option[gitdomain.LocalBranchName]
	LocalBranches     gitdomain.LocalBranchNames
	UnvalidatedConfig config.UnvalidatedConfig
	UnvalidatedMain   Option[gitdomain.LocalBranchName]
}
