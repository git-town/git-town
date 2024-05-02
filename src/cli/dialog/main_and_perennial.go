package dialog

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

func MainAndPerennials(args MainAndPerennialsArgs) (mainBranch gitdomain.LocalBranchName, additionalPerennials gitdomain.LocalBranchNames, aborted bool, err error) {
	unvalidatedMain, hasMain := args.UnvalidatedMain.Get()
	if hasMain {
		return unvalidatedMain, args.UnvalidatedPerennials, false, nil
	}
	if args.HasConfigFile {
		return unvalidatedMain, args.UnvalidatedPerennials, false, errors.New(messages.ConfigMainbranchInConfigFile)
	}
	fmt.Print(messages.ConfigNeeded)
	mainBranch, aborted, err = MainBranch(args.LocalBranches, args.GetDefaultBranch(), args.DialogInputs.Next())
	if err != nil || aborted {
		return mainBranch, args.UnvalidatedPerennials, aborted, err
	}
	additionalPerennials, aborted, err = PerennialBranches(args.LocalBranches, args.UnvalidatedPerennials, mainBranch, args.DialogInputs.Next())
	return mainBranch, additionalPerennials, aborted, err
}

type MainAndPerennialsArgs struct {
	Backend               *git.BackendCommands
	DialogInputs          *components.TestInputs
	GetDefaultBranch      func() Option[gitdomain.LocalBranchName]
	HasConfigFile         bool
	LocalBranches         gitdomain.LocalBranchNames
	LocalGitConfig        *gitconfig.Access
	UnvalidatedMain       Option[gitdomain.LocalBranchName]
	UnvalidatedPerennials gitdomain.LocalBranchNames
}
