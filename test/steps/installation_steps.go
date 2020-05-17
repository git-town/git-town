package steps

import (
	"github.com/cucumber/godog"
)

// InstallationSteps defines Cucumber step implementations around installation of Git Town.
func InstallationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my computer has a broken "([^"]*)" tool installed$`, func(name string) error {
		return fs.state.gitEnv.DeveloperShell.MockBrokenCommand(name)
	})

	suite.Step(`^my computer has no tool to open browsers installed$`, func() error {
		return fs.state.gitEnv.DeveloperShell.MockNoCommandsInstalled()
	})

	suite.Step(`^my computer has the "([^"]*)" tool installed$`, func(tool string) error {
		return fs.state.gitEnv.DeveloperShell.MockCommand(tool)
	})

	suite.Step(`^I have Git "([^"]*)" installed$`, func(version string) error {
		err := fs.state.gitEnv.DeveloperShell.MockGit(version)
		return err
	})
}
