package steps

import (
	"github.com/DATA-DOG/godog"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
func ConfigurationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I haven\'t configured Git Town yet$`, fs.iHaventConfiguredGitTownYet)
}

func (fs *FeatureState) iHaventConfiguredGitTownYet() error {
	// delete_main_branch_configuration
	// delete_perennial_branches_configuration
	return nil
}
