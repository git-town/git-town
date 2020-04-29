package steps

import "github.com/cucumber/godog"

// ConflictSteps defines Gherkin step implementations around merge conflicts.
func ConflictSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I resolve the conflict in "([^"]*)"$`, func(filename string) error {
		err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFile(filename, "resolved content")
		if err != nil {
			return err
		}
		err = fs.activeScenarioState.gitEnvironment.DeveloperRepo.StageFiles(filename)
		if err != nil {
			return err
		}
		return nil
	})
}
