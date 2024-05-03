package validate

import (
	"errors"
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/git-town/git-town/v14/src/messages"
)

func Config(args ConfigArgs) (validatedResult config.ValidatedConfig, aborted bool, err error) {
	// check Git user data
	if args.Unvalidated.Config.GitUserEmail.IsNone() {
		return validatedResult, false, errors.New(messages.GitUserEmailMissing)
	}
	if args.Unvalidated.Config.GitUserName.IsNone() {
		return validatedResult, false, errors.New(messages.GitUserNameMissing)
	}

	// enter and save main and perennials
	validatedMain, additionalPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
		DialogInputs:          args.TestInputs,
		GetDefaultBranch:      args.Backend.DefaultBranch,
		HasConfigFile:         args.Unvalidated.ConfigFile.IsSome(),
		LocalBranches:         args.LocalBranches,
		UnvalidatedMain:       args.Unvalidated.Config.MainBranch,
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
		if err = args.Unvalidated.SetParent(branch, parent); err != nil {
			return validatedResult, abort, err
		}
	}
	if len(additionalPerennials) > 0 {
		newPerennials := append(args.Unvalidated.Config.PerennialBranches, additionalPerennials...)
		if err = args.Unvalidated.SetPerennialBranches(newPerennials); err != nil {
			return validatedResult, false, err
		}
	}

	// remove outdated lineage
	err = args.Unvalidated.RemoveOutdatedConfiguration(args.LocalBranches)
	if err != nil {
		return validatedResult, abort, err
	}
	err = cleanupPerennialParentEntries(args.Unvalidated.Config.Lineage, args.Unvalidated.Config.PerennialBranches, args.Unvalidated.GitConfig, args.FinalMessages)

	validatedConfig := configdomain.ValidatedConfig{
		UnvalidatedConfig: args.Unvalidated.Config,
		MainBranch:        validatedMain,
	}
	vConfig := config.ValidatedConfig{
		Config: validatedConfig,
	}
	return vConfig, false, nil
}

type ConfigArgs struct {
	Backend            *git.BackendCommands
	BranchesToValidate gitdomain.LocalBranchNames
	FinalMessages      *stringslice.Collector
	LocalBranches      gitdomain.LocalBranchNames
	TestInputs         *components.TestInputs
	Unvalidated        config.UnvalidatedConfig
}

// cleanupPerennialParentEntries removes outdated entries from the configuration.
func cleanupPerennialParentEntries(lineage configdomain.Lineage, perennialBranches gitdomain.LocalBranchNames, access gitconfig.Access, finalMessages *stringslice.Collector) error {
	for _, perennialBranch := range perennialBranches {
		if lineage.Parent(perennialBranch).IsSome() {
			if err := access.RemoveLocalConfigValue(gitconfig.NewParentKey(perennialBranch)); err != nil {
				return err
			}
			lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
	return nil
}
