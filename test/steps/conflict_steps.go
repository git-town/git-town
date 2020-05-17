package steps

import "github.com/cucumber/godog"

// ConflictSteps defines Gherkin step implementations around merge conflicts.
func ConflictSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(filename, content string) error {
		if content == "" {
			content = "resolved content"
		}
		err := fs.state.gitEnvironment.DeveloperRepo.CreateFile(filename, content)
		if err != nil {
			return err
		}
		err = fs.state.gitEnvironment.DeveloperRepo.StageFiles(filename)
		if err != nil {
			return err
		}
		return nil
	})
}
