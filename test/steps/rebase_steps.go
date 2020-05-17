package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// RebaseSteps defines Gherkin step implementations around rebases.
func RebaseSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^my repo (?:now|still) has a rebase in progress$`, func() error {
		hasRebase, err := state.gitEnv.DevRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if !hasRebase {
			return fmt.Errorf("expected rebase in progress")
		}
		return nil
	})

	suite.Step(`^there is no rebase in progress$`, func() error {
		hasRebase, err := state.gitEnv.DevRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if hasRebase {
			return fmt.Errorf("expected no rebase in progress")
		}
		return nil
	})
}
