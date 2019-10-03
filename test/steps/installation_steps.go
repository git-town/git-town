package steps

import "github.com/DATA-DOG/godog"

// InstallationSteps provides Cucumber step implementations around installation of Git Town.
func InstallationSteps(s *godog.Suite, gtf *GitTownFeature) {
	s.Step(`^I have Git "([^"]*)" installed$`, gtf.iHaveGitInstalled)
}

func (gtf *GitTownFeature) iHaveGitInstalled(arg1 string) error {
	err := gtf.gitEnvironment.DeveloperRepo.AddTempShellOverride(
		"git",
		`#!/usr/bin/env bash
		echo "git version 2.6.2"`)
	return err
}
