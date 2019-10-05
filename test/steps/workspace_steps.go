package steps

import (
	"os"
	"path"

	"github.com/DATA-DOG/godog"
)

// WorkspaceSteps defines Cucumber step implementations around Git workspace management.
func WorkspaceSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my workspace is currently not a Git repository$`, fs.myWorkspaceIsCurrentlyNotAGitRepository)
}

func (fs *FeatureState) myWorkspaceIsCurrentlyNotAGitRepository() error {
	os.RemoveAll(path.Join(fs.activeScenarioState.gitEnvironment.DeveloperRepo.Dir, ".git"))
	return nil
}
