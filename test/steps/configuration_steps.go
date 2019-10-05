package steps

import (
	"github.com/DATA-DOG/godog"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
func ConfigurationSteps(s *godog.Suite, state *FeatureState) {
	s.Step(`^I haven\'t configured Git Town yet$`, state.iHaventConfiguredGitTownYet)
}

func (state *FeatureState) iHaventConfiguredGitTownYet() error {
	// delete_main_branch_configuration
	// delete_perennial_branches_configuration
	return nil
}
