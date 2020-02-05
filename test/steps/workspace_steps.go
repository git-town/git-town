package steps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DATA-DOG/godog"
)

// WorkspaceSteps defines Cucumber step implementations around Git workspace management.
func WorkspaceSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my workspace has an uncommitted file$`, func() error {
		fs.activeScenarioState.uncommittedFileName = "uncommitted file"
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFile(fs.activeScenarioState.uncommittedFileName, "uncommitted content")
	})

	suite.Step(`^my workspace has an uncommitted file with name: "([^"]+)"$`, func(filename string) error {
		fs.activeScenarioState.uncommittedFileName = filename
		return fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFile(fs.activeScenarioState.uncommittedFileName, "uncommitted content")
	})

	suite.Step(`^my workspace is currently not a Git repository$`, func() error {
		os.RemoveAll(filepath.Join(fs.activeScenarioState.gitEnvironment.DeveloperRepo.Dir, ".git"))
		return nil
	})

	suite.Step(`^my workspace still contains my uncommitted file$`, func() error {
		hasFile, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasFile(fs.activeScenarioState.uncommittedFileName, "uncommitted content")
		if err != nil {
			return fmt.Errorf("cannot determine if workspace contains uncommitted file: %w", err)
		}
		if !hasFile {
			return fmt.Errorf("expected the uncommitted file but didn't find one")
		}
		return nil
	})
}
