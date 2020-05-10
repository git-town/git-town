package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// MergeSteps defines Cucumber step implementations around Git merges
// nolint:funlen
func MergeSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^there is no merge in progress$`, func() error {
		has, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasMergeInProgress()
		if err != nil {
			return err
		}
		if has {
			return fmt.Errorf("expected no merge in progress, but has one")
		}
		return nil
	})
}
