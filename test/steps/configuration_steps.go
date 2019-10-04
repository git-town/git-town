package steps

import (
	"github.com/DATA-DOG/godog"
)

// ConfigurationSteps provides Cucumber step implementations around configuration.
func ConfigurationSteps(s *godog.Suite, gtf *GitTownFeature) {
	s.Step(`^I haven\'t configured Git Town yet$`, gtf.iHaventConfiguredGitTownYet)
}

func (gtf *GitTownFeature) iHaventConfiguredGitTownYet() error {
	// delete_main_branch_configuration
	// delete_perennial_branches_configuration
	return nil
}
