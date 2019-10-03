package steps

import (
	"github.com/DATA-DOG/godog"
)

// InstallationSteps provides Cucumber step implementations around installation of Git Town.
func InstallationSteps(s *godog.Suite) {
	s.Step(`^I have Git "([^"]*)" installed$`,
		func(arg1 string) error {
			err := gitEnvironment.DeveloperRepo.AddTempShellOverride(
				"git",
				`#!/usr/bin/env bash
		echo "git version 2.6.2"`)
			return err
		})
}
