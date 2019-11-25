package steps

import (
	"strconv"

	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
func ConfigurationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I haven\'t configured Git Town yet$`, func() error {
		// NOTE: nothing to do here yet since we don't configure Git Town in Go specs at this point.
		// In the future:
		// - delete_main_branch_configuration
		// - delete_perennial_branches_configuration
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is set to "(true|false)"$`, func(value string) error {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return errors.Wrapf(err, "cannot parse %q into bool", value)
		}
		outcome, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Configuration().SetNewBranchPush(b, false)
		if err != nil {
			return errors.Wrapf(err, "cannot set new-branch-push-flag configuration to %q: %s", value, err)
		}
		return err
	})
}
