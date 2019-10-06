package steps

import (
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

	suite.Step(`^the "([^"]+)" configuration is set to "([^"]+)"$`, func(name, value string) error {
		output, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Run("git", "config", "git-town."+name, value)
		if err != nil {
			return errors.Wrapf(err, "cannot set Git configuration %q to %q: %s", name, value, output)
		}
		return err
	})
}
