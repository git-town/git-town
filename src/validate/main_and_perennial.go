package validate

import (
	"errors"
	"fmt"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

func MainAndPerennials(args MainAndPerennialsArgs) (MainAndPerennialsResult, error) {
	unvalidatedMain, hasMain := args.UnvalidatedMain.Get()
	if hasMain {
		return MainAndPerennialsResult{
			ValidatedMain:       unvalidatedMain,
			ValidatedPerennials: args.UnvalidatedPerennials,
		}, nil
	}
	if args.HasConfigFile {
		return emptyMainAndPerennialsResult(), errors.New(messages.ConfigMainbranchInConfigFile)
	}
	fmt.Print(messages.ConfigNeeded)
	var err error
	var aborted bool
	validatedMain, aborted, err := dialog.MainBranch(args.LocalBranches, args.Backend.DefaultBranch(), args.DialogInputs.Next())
	if err != nil || aborted {
		return emptyMainAndPerennialsResult(), err
	}
	if validatedMain != unvalidatedMain {
		err := args.LocalGitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, validatedMain.String())
		if err != nil {
			return emptyMainAndPerennialsResult(), err
		}
	}
	validatedPerennials, aborted, err := dialog.PerennialBranches(args.LocalBranches, args.UnvalidatedPerennials, validatedMain, args.DialogInputs.Next())
	if err != nil || aborted {
		return emptyMainAndPerennialsResult(), err
	}
	if slices.Compare(validatedPerennials, args.UnvalidatedPerennials) != 0 {
		err := args.LocalGitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, validatedPerennials.Join(" "))
		if err != nil {
			return emptyMainAndPerennialsResult(), err
		}
	}
	return MainAndPerennialsResult{
		ValidatedMain:       validatedMain,
		ValidatedPerennials: validatedPerennials,
	}, nil
}

type MainAndPerennialsArgs struct {
	Backend               *git.BackendCommands
	DialogInputs          *components.TestInputs
	HasConfigFile         bool
	LocalBranches         gitdomain.LocalBranchNames
	LocalGitConfig        *gitconfig.Access
	UnvalidatedMain       Option[gitdomain.LocalBranchName]
	UnvalidatedPerennials gitdomain.LocalBranchNames
}

type MainAndPerennialsResult struct {
	ValidatedMain       gitdomain.LocalBranchName
	ValidatedPerennials gitdomain.LocalBranchNames
	Err                 error
}

func emptyMainAndPerennialsResult() MainAndPerennialsResult {
	return MainAndPerennialsResult{
		ValidatedMain:       gitdomain.EmptyLocalBranchName(),
		ValidatedPerennials: gitdomain.LocalBranchNames{},
	}
}
