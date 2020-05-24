package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// OfflineSteps defines Cucumber step implementations around offline functionality.
func OfflineSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^offline mode is disabled$`, func() error {
		offline := state.gitEnv.DevRepo.IsOffline()
		if offline {
			return fmt.Errorf("expected to not be offline but am")
		}
		return nil
	})
	suite.Step(`^offline mode is enabled$`, func() error {
		offline := state.gitEnv.DevRepo.IsOffline()
		if !offline {
			return fmt.Errorf("expected to be offline but am not")
		}
		return nil
	})
}
