package validate

import (
	"errors"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

func Config(unvalidated config.UnvalidatedConfig, branchesToValidate gitdomain.LocalBranchNames, localBranches gitdomain.LocalBranchNames, backend *git.BackendCommands, testInputs *components.TestInputs) (*config.ValidatedConfig, error) {
	validateResult, err := MainAndPerennials(MainAndPerennialsArgs{
		UnvalidatedMain:       unvalidated.Config.MainBranch,
		UnvalidatedPerennials: unvalidated.Config.PerennialBranches,
	})
	if err != nil {
		return nil, err
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
	validatedConfig := configdomain.ValidatedConfig{
		UnvalidatedConfig: unvalidated.Config,
		MainBranch:        validateResult.ValidatedMain,
	}
	vConfig := config.ValidatedConfig{
		Config: validatedConfig,
	}
	return &vConfig, nil
}
