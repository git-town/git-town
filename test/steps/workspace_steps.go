package steps

import (
	"github.com/DATA-DOG/godog"
)

// WorkspaceSteps provides Cucumber step implementations around Git workspace management.
func WorkspaceSteps(s *godog.Suite) {
	s.Step("^my workspace is currently not a Git repository$",
		func() error {
			// FileUtils.rm_rf '.git'
			return nil
		})
}
