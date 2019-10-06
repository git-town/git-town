package steps

import (
	"fmt"
	"os"
	"path"

	"github.com/DATA-DOG/godog"
	"github.com/pkg/errors"
)

// WorkspaceSteps defines Cucumber step implementations around Git workspace management.
func WorkspaceSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my workspace has an uncommitted file$`, func() error {
		fs.activeScenarioState.gitEnvironment.DeveloperRepo.CreateFile("uncommitted file", "uncommitted content")
		return nil
	})

	suite.Step(`^my workspace is currently not a Git repository$`, func() error {
		os.RemoveAll(path.Join(fs.activeScenarioState.gitEnvironment.DeveloperRepo.Dir, ".git"))
		return nil
	})

	suite.Step(`^my workspace still contains my uncommitted file$`, func() error {
		hasFile, err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.HasFile("uncommitted file", "uncommitted content")
		if err != nil {
			return errors.Wrap(err, "cannot determine if workspace contains uncommitted file")
		}
		if !hasFile {
			return fmt.Errorf("expected the uncommitted file but didn't find one")
		}
		return nil
	})
}
