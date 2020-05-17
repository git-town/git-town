package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// FolderSteps defines Cucumber step implementations around folders.
func FolderSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I am in the project root folder$`, func() error {
		actual, err := fs.state.gitEnvironment.DeveloperRepo.LastActiveDir()
		if err != nil {
			return fmt.Errorf("cannot determine the current working directory: %w", err)
		}
		expected := fs.state.gitEnvironment.DeveloperRepo.Dir
		if actual != expected {
			return fmt.Errorf("expected to be in %q but am in %q", expected, actual)
		}
		return nil
	})
}
