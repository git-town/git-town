package steps

import "github.com/DATA-DOG/godog"

// InstallationSteps defines Cucumber step implementations around installation of Git Town.
func InstallationSteps(s *godog.Suite, state *FeatureState) {
	s.Step(`^I have Git "([^"]*)" installed$`, state.iHaveGitInstalled)
}

func (state *FeatureState) iHaveGitInstalled(arg1 string) error {
	err := state.gitEnvironment.DeveloperRepo.AddTempShellOverride(
		"git",
		`#!/usr/bin/env bash
		echo "git version 2.6.2"`)
	return err
}
