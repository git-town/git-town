package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// RebaseSteps defines Gherkin step implementations around rebases.
func RebaseSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repo has a rebase in progress$`, func() error {
		hasRebase, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if !hasRebase {
			return fmt.Errorf("expected rebase in progress")
		}
		return nil
	})

	suite.Step(`^there is no rebase in progress$`, func() error {
		hasRebase, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if hasRebase {
			return fmt.Errorf("expected no rebase in progress")
		}
		return nil
	})
}
