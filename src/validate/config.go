package validate

import (
	"errors"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
)

func Config(args ConfigArgs) (validatedResult config.Config, aborted bool, err error) {
	// check Git user data
	if args.Unvalidated.Config.GitUserEmail == "" {
		return validatedResult, false, errors.New(messages.GitUserEmailMissing)
	}
	if args.Unvalidated.Config.GitUserName == "" {
		return validatedResult, false, errors.New(messages.GitUserNameMissing)
	}

	// enter and save main and perennials
	var validatedMain gitdomain.LocalBranchName
	var validatedPerennials gitdomain.LocalBranchNames
	if args.Unvalidated.Config.MainBranch.IsEmpty() {
		validatedMain, validatedPerennials, aborted, err = dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
			DialogInputs:          args.TestInputs,
			GetDefaultBranch:      args.Backend.DefaultBranch,
			HasConfigFile:         args.Unvalidated.ConfigFile.IsSome(),
			LocalBranches:         args.LocalBranches,
			UnvalidatedMain:       None[gitdomain.LocalBranchName](),
			UnvalidatedPerennials: args.Unvalidated.Config.PerennialBranches,
		})
		if err != nil || aborted {
			return validatedResult, aborted, err
		}
		if err = args.Unvalidated.SetMainBranch(validatedMain); err != nil {
			return validatedResult, false, err
		}
		if err = args.Unvalidated.SetPerennialBranches(validatedPerennials); err != nil {
			return validatedResult, false, err
		}
	} else {
		validatedMain = args.Unvalidated.Config.MainBranch
		validatedPerennials = args.Unvalidated.Config.PerennialBranches
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, exit, err := dialog.Lineage(dialog.LineageArgs{
		BranchesToVerify: args.BranchesToValidate,
		Config:           args.Unvalidated.Config,
		DefaultChoice:    validatedMain,
		DialogTestInputs: args.TestInputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       validatedMain,
	})
	if err != nil || exit {
		return validatedResult, exit, err
	}
	for branch, parent := range additionalLineage {
		if err = args.Unvalidated.SetParent(branch, parent); err != nil {
			return validatedResult, false, err
		}
	}
	if len(additionalPerennials) > 0 {
		validatedPerennials = append(validatedPerennials, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(validatedPerennials); err != nil {
			return validatedResult, false, err
		}
	}

	return args.Unvalidated, false, err
}

type ConfigArgs struct {
	Backend            git.BackendCommands
	BranchesToValidate gitdomain.LocalBranchNames
	FinalMessages      stringslice.Collector
	LocalBranches      gitdomain.LocalBranchNames
	TestInputs         components.TestInputs
	Unvalidated        config.Config
}
