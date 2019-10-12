package steps

import (
	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// GitTownSteps defines Cucumber step implementations around the Git Town setup.
func GitTownSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is in offline mode$`, func() error {
		output, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.Run("git", "config", "git-town.offline", "true")
		if err != nil {
			return errors.Wrapf(err, "cannot enable offline mode: %s", output)
		}
		return nil
	})
}
