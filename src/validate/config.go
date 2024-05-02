package validate

import (
	"errors"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

func Config(args ConfigArgs) (validatedResult config.Config, aborted bool, err error) {
	// enter and save main and perennials
	var mainOpt Option[gitdomain.LocalBranchName]
	if !args.Unvalidated.Config.MainBranch.IsEmpty() {
		mainOpt = Some(args.Unvalidated.Config.MainBranch)
	}
	validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
		UnvalidatedMain:       mainOpt,
		UnvalidatedPerennials: args.Unvalidated.Config.PerennialBranches,
	})
	if err != nil || aborted {
		return validatedResult, aborted, err
	}
	if err = args.Unvalidated.SetMainBranch(validatedMain); err != nil {
		return validatedResult, false, err
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return validatedResult, false, err
		}
	}

	// check Git user data
	if args.Unvalidated.Config.GitUserEmail.IsNone() {
		return validatedResult, false, errors.New(messages.GitUserEmailMissing)
	}
	if args.Unvalidated.Config.GitUserName.IsNone() {
		return validatedResult, false, errors.New(messages.GitUserNameMissing)
	}

	// enter and save missing parent branches
	additionalLineage, additionalPerennials, abort, err := dialog.Lineage(dialog.LineageArgs{
		BranchesToVerify: args.BranchesToValidate,
		Config:           args.Unvalidated.Config,
		DefaultChoice:    validatedMain,
		DialogTestInputs: args.TestInputs,
		LocalBranches:    args.LocalBranches,
		MainBranch:       validatedMain,
	})
	if err != nil || abort {
		return validatedResult, abort, err
	}
	for branch, parent := range additionalLineage {
		args.Unvalidated.SetParent(branch, parent)
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return validatedResult, false, err
		}
	}
	fullConfig := configdomain.FullConfig{
		UnvalidatedConfig: args.Unvalidated.Config,
		MainBranch:        validatedMain,
	}
	vConfig := config.Config{
		Config: fullConfig,
	}
	return vConfig, false, nil
}

type ConfigArgs struct {
	Backend            *git.BackendCommands
	BranchesToValidate gitdomain.LocalBranchNames
	LocalBranches      gitdomain.LocalBranchNames
	TestInputs         *components.TestInputs
	Unvalidated        config.Config
}
