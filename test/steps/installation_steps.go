package steps

import (
	"github.com/cucumber/godog"
)

// InstallationSteps defines Cucumber step implementations around installation of Git Town.
func InstallationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my computer has no tool to open browsers installed$`, func() error {
		return fs.activeScenarioState.gitEnvironment.DeveloperShell.MockNoCommandsInstalled()
	})

	suite.Step(`^my computer has the "([^"]*)" tool installed$`, func(tool string) error {
		return fs.activeScenarioState.gitEnvironment.DeveloperShell.MockCommand(tool)
	})

	suite.Step(`^I have Git "([^"]*)" installed$`, func(version string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperShell.MockGit(version)
		return err
	})
}
