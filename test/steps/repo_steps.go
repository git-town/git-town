package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// RepoSteps defines Gherkin step implementations around running things in subshells.
func RepoSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repo has a merge in progress$`, func() error {
		has, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasMergeInProgress()
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("expected merge in progress but no merge detected")
		}
		return nil
	})

	suite.Step(`^my repo has an upstream repo$`, func() error {
		return fs.activeScenarioState.gitEnvironment.AddUpstream()
	})

	suite.Step(`^my repository knows about the remote branch$`, func() error {
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.Fetch()
	})
}
