package steps

import (
	"github.com/DATA-DOG/godog"
)

// WorkspaceSteps provides Cucumber step implementations around Git workspace management.
func WorkspaceSteps(s *godog.Suite, gtf *GitTownFeature) {
	s.Step("^my workspace is currently not a Git repository$", gtf.myWorkspaceIsCurrentlyNotAGitRepository)
}

func (gtf *GitTownFeature) myWorkspaceIsCurrentlyNotAGitRepository() error {
	// FileUtils.rm_rf '.git'
	return nil
}
