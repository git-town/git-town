package steps

import (
	"github.com/DATA-DOG/godog"
)

// WorkspaceSteps defines Cucumber step implementations around Git workspace management.
func WorkspaceSteps(s *godog.Suite, state *FeatureState) {
	s.Step("^my workspace is currently not a Git repository$", state.myWorkspaceIsCurrentlyNotAGitRepository)
}

func (state *FeatureState) myWorkspaceIsCurrentlyNotAGitRepository() error {
	// FileUtils.rm_rf '.git'
	return nil
}
