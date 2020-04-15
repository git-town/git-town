package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// FileSteps defines Cucumber step implementations around files.
func FileSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my uncommitted file is stashed$`, func() error {
		uncommittedFiles, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.UncommittedFiles()
		if err != nil {
			return err
		}
		for ucf := range uncommittedFiles {
			if uncommittedFiles[ucf] == fs.activeScenarioState.uncommittedFileName {
				return fmt.Errorf("expected file %q to be stashed but it is still uncommitted", fs.activeScenarioState.uncommittedFileName)
			}
		}
		stashSize, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.StashSize()
		if err != nil {
			return err
		}
		if stashSize != 1 {
			return fmt.Errorf("expected 1 stash but found %d", stashSize)
		}
		return nil
	})
}
