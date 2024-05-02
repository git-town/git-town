package validate

import (
	"errors"
	"slices"

	"github.com/git-town/git-town/v14/src/cli/dialog"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

func Config(args ConfigArgs) (validatedResult config.ValidatedConfig, aborted bool, err error) {
	validatedMain, validatedPerennials, aborted, err := dialog.MainAndPerennials(dialog.MainAndPerennialsArgs{
		UnvalidatedMain:       args.unvalidated.Config.MainBranch,
		UnvalidatedPerennials: args.unvalidated.Config.PerennialBranches,
	})
	if err != nil || aborted {
		return validatedResult, aborted, err
	}
	err = args.LocalGitConfig.SetLocalConfigValue(gitconfig.KeyMainBranch, validatedMain.String())
	if err != nil {
		return emptyMainAndPerennialsResult(), err
	}
	if unvalidated.Config.GitUserEmail.IsNone() {
		return nil, errors.New(messages.GitUserEmailMissing)
	}
	if unvalidated.Config.GitUserName.IsNone() {
		return nil, errors.New(messages.GitUserNameMissing)
	}
	additionalLineage, additionalPerennials, aborted, err := Lineage(LineageArgs{
		BranchesToVerify: branchesToValidate,
		Config:           unvalidated.Config,
		DefaultChoice:    validateResult.ValidatedMain,
		DialogTestInputs: testInputs,
		LocalBranches:    localBranches,
		MainBranch:       validateResult.ValidatedMain,
	})
	if err != nil {
		return nil, err
	}
	if slices.Compare(additionalPerennials, args.UnvalidatedPerennials) != 0 {
		err := args.LocalGitConfig.SetLocalConfigValue(gitconfig.KeyPerennialBranches, additionalPerennials.Join(" "))
		if err != nil {
			return emptyMainAndPerennialsResult(), err
		}
	}
	validatedConfig := configdomain.ValidatedConfig{
		UnvalidatedConfig: unvalidated.Config,
		MainBranch:        validateResult.ValidatedMain,
	}
	vConfig := config.ValidatedConfig{
		Config: validatedConfig,
	}
	return &vConfig, nil
}

type ConfigArgs struct {
	unvalidated        config.UnvalidatedConfig
	branchesToValidate gitdomain.LocalBranchNames
	localBranches      gitdomain.LocalBranchNames
	backend            *git.BackendCommands
	testInputs         *components.TestInputs
}
