Feature: unsupported hosting platform

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    When I run "git-town propose"

  Scenario: result
    Then Git Town prints the error:
      """
      unsupported hosting platform

      This command requires hosting on one of these services:
      * Bitbucket
      * GitHub
      * GitLab
      * Gitea
      """
