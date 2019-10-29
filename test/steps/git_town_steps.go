package steps

import (
	"github.com/DATA-DOG/godog"
)

// GitTownSteps defines Cucumber step implementations around the Git Town setup.
func GitTownSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is in offline mode$`, func() error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.SetOffline(true)
	})
}
