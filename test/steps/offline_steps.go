package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// OfflineSteps defines Cucumber step implementations around offline functionality.
func OfflineSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^offline mode is disabled$`, func() error {
		offline, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.IsOffline()
		if err != nil {
			return err
		}
		if offline {
			return fmt.Errorf("expected to not be offline but am")
		}
		return nil
	})
	suite.Step(`^offline mode is enabled$`, func() error {
		offline, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.IsOffline()
		if err != nil {
			return err
		}
		if !offline {
			return fmt.Errorf("expected to be offline but am not")
		}
		return nil
	})
}
