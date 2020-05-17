package steps

import (
	"github.com/cucumber/godog"
)

// GitTownSteps defines Cucumber step implementations around the Git Town setup.
func GitTownSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^Git Town is in offline mode$`, func() error {
		return fs.state.gitEnv.DeveloperRepo.SetOffline(true)
	})
}
