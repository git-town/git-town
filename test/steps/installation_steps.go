package steps

import "github.com/DATA-DOG/godog"

// InstallationSteps defines Cucumber step implementations around installation of Git Town.
func InstallationSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^I have Git "([^"]*)" installed$`, fs.iHaveGitInstalled)
}

func (fs *FeatureState) iHaveGitInstalled(arg1 string) error {
	err := fs.activeScenarioState.gitEnvironment.DeveloperRepo.AddTempShellOverride(
		"git",
		`#!/usr/bin/env bash
		echo "git version 2.6.2"`)
	return err
}
