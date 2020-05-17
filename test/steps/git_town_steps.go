package steps

import (
	"github.com/cucumber/godog"
)

// GitTownSteps defines Cucumber step implementations around the Git Town setup.
func GitTownSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^Git Town is in offline mode$`, func() error {
		return state.gitEnv.DevRepo.SetOffline(true)
	})
}
