Feature: unsupported hosting platform

  Background:
    Given the current branch is a feature branch "feature"
    When I run "git-town propose"

  Scenario: result
    Then it prints the error:
      """
      unsupported hosting platform

      This command requires hosting on one of these services:
      * Bitbucket
      * GitHub
      * GitLab
      * Gitea
      """
