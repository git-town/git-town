package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// CoworkerSteps defines Gherkin step implementations around a coworker.
func CoworkerSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my coworker fetches updates$`, func() error {
		_, err := fs.activeScenarioState.gitEnvironment.CoworkerRepo.Shell.Run("git", "fetch")
		return err
	})

	suite.Step(`^my coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_, err := fs.activeScenarioState.gitEnvironment.CoworkerRepo.Shell.Run("git", "config", "git-town-branch."+childBranch+".parent", parentBranch)
		return err
	})

	suite.Step(`^my coworker is on the "([^"]*)" branch$`, func(branchName string) error {
		err := fs.activeScenarioState.gitEnvironment.CoworkerRepo.CheckoutBranch(branchName)
		if err != nil {
			return fmt.Errorf("cannot change to branch %q: %w", branchName, err)
		}
		return nil
	})

	suite.Step(`^my coworker runs "([^"]+)"$`, func(command string) error {
		fs.activeScenarioState.lastRunResult, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.CoworkerRepo.Shell.RunString(command)
		return nil
	})
}
