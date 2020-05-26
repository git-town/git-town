package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// OfflineSteps defines Cucumber step implementations around offline functionality.
func OfflineSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^offline mode is disabled$`, func() error {
		state.gitEnv.DevRepo.Configuration.Reload()
		if state.gitEnv.DevRepo.IsOffline() {
			return fmt.Errorf("expected to not be offline but am")
		}
		return nil
	})
	suite.Step(`^offline mode is enabled$`, func() error {
		state.gitEnv.DevRepo.Configuration.Reload()
		if !state.gitEnv.DevRepo.IsOffline() {
			return fmt.Errorf("expected to be offline but am not")
		}
		return nil
	})
}
