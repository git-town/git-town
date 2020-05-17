package steps

import (
	"github.com/cucumber/godog"
)

// CoworkerSteps defines Gherkin step implementations around a coworker.
func CoworkerSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I am collaborating with a coworker$`, func() error {
		return fs.gitEnv.AddCoworkerRepo()
	})

	suite.Step(`^my coworker fetches updates$`, func() error {
		return fs.gitEnv.CoworkerRepo.Fetch()
	})

	suite.Step(`^my coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = fs.gitEnv.CoworkerRepo.Configuration(false).SetParentBranch(childBranch, parentBranch)
		return nil
	})

	suite.Step(`^my coworker is on the "([^"]*)" branch$`, func(branchName string) error {
		return fs.gitEnv.CoworkerRepo.CheckoutBranch(branchName)
	})

	suite.Step(`^my coworker runs "([^"]+)"$`, func(command string) error {
		fs.runRes, fs.runErr = fs.gitEnv.CoworkerRepo.Shell.RunString(command)
		return nil
	})
}
