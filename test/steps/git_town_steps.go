package steps

import (
	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// GitTownSteps defines Cucumber step implementations around the Git Town setup.
func GitTownSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is in offline mode$`, func() error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.SetOffline(true)
		if err != nil {
			return errors.Wrap(err, "cannot enable offline mode")
		}
		return nil
	})
}
