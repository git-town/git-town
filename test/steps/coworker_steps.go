package steps

import (
	"github.com/cucumber/godog"
)

// CoworkerSteps defines Gherkin step implementations around a coworker.
func CoworkerSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I am collaborating with a coworker$`, func() error {
		return fs.activeScenarioState.gitEnvironment.AddCoworkerRepo()
	})

	suite.Step(`^my coworker fetches updates$`, func() error {
		return fs.activeScenarioState.gitEnvironment.CoworkerRepo.Fetch()
	})

	suite.Step(`^my coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = fs.activeScenarioState.gitEnvironment.CoworkerRepo.Configuration(false).SetParentBranch(childBranch, parentBranch)
		return nil
	})

	suite.Step(`^my coworker is on the "([^"]*)" branch$`, func(branchName string) error {
		return fs.activeScenarioState.gitEnvironment.CoworkerRepo.CheckoutBranch(branchName)
	})

	suite.Step(`^my coworker runs "([^"]+)"$`, func(command string) error {
		fs.activeScenarioState.lastRunResult, fs.activeScenarioState.lastRunErr = fs.activeScenarioState.gitEnvironment.CoworkerRepo.Shell.RunString(command)
		return nil
	})
}
