package steps

import (
	"github.com/cucumber/godog"
)

// CoworkerSteps defines Gherkin step implementations around a coworker.
func CoworkerSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^I am collaborating with a coworker$`, func() error {
		return state.gitEnv.AddCoworkerRepo()
	})

	suite.Step(`^my coworker fetches updates$`, func() error {
		return state.gitEnv.CoworkerRepo.Fetch()
	})

	suite.Step(`^my coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = state.gitEnv.CoworkerRepo.Configuration(false).SetParentBranch(childBranch, parentBranch)
		return nil
	})

	suite.Step(`^my coworker is on the "([^"]*)" branch$`, func(branchName string) error {
		return state.gitEnv.CoworkerRepo.CheckoutBranch(branchName)
	})

	suite.Step(`^my coworker runs "([^"]+)"$`, func(command string) error {
		state.runRes, state.runErr = state.gitEnv.CoworkerRepo.Shell.RunString(command)
		return nil
	})
}
