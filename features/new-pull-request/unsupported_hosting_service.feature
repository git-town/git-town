Feature: git-new-pull-request: when origin is unsupported

  As a developer trying to create a pull request in a repository on an unsupported hosting service
  I should get an error that my hosting service is not supported
  So that I know why the command doesn't work.

  Background:
    Given my repo has a feature branch named "feature"
    And I am on the "feature" branch
    When I run "git-town new-pull-request"

  Scenario: result
    Then it prints the error:
      """
      unsupported hosting service

      This command requires hosting on one of these services:
      * Bitbucket
      * GitHub
      * GitLab
      * Gitea
      """
