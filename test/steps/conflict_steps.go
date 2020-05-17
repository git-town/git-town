package steps

import "github.com/cucumber/godog"

// ConflictSteps defines Gherkin step implementations around merge conflicts.
func ConflictSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(filename, content string) error {
		if content == "" {
			content = "resolved content"
		}
		err := state.gitEnv.DevRepo.CreateFile(filename, content)
		if err != nil {
			return err
		}
		err = state.gitEnv.DevRepo.StageFiles(filename)
		if err != nil {
			return err
		}
		return nil
	})
}
