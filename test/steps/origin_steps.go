package steps

import "github.com/cucumber/godog"

// OriginSteps defines Cucumber step implementations around Git origins.
func OriginSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repo does not have a remote origin$`, func() error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.RemoveRemote("origin")
		if err != nil {
			return err
		}
		fs.activeScenarioState.gitEnvironment.OriginRepo = nil
		return nil
	})
}
