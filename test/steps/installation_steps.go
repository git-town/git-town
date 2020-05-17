package steps

import (
	"github.com/cucumber/godog"
)

// InstallationSteps defines Cucumber step implementations around installation of Git Town.
func InstallationSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^my computer has a broken "([^"]*)" tool installed$`, func(name string) error {
		return state.gitEnv.DevShell.MockBrokenCommand(name)
	})

	suite.Step(`^my computer has no tool to open browsers installed$`, func() error {
		return state.gitEnv.DevShell.MockNoCommandsInstalled()
	})

	suite.Step(`^my computer has the "([^"]*)" tool installed$`, func(tool string) error {
		return state.gitEnv.DevShell.MockCommand(tool)
	})

	suite.Step(`^I have Git "([^"]*)" installed$`, func(version string) error {
		err := state.gitEnv.DevShell.MockGit(version)
		return err
	})
}
