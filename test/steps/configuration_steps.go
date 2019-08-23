package steps

import (
	"github.com/DATA-DOG/godog"
)

func ConfigurationSteps(s *godog.Suite) {

	s.Step(`^I haven\'t configured Git Town yet$`, func() error {
		// delete_main_branch_configuration
		// delete_perennial_branches_configuration
		return nil
	})

}
