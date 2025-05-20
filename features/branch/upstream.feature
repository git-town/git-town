Feature: display the local branch hierarchy when there is an upstream

  Background:
    Given a Git repo with origin
    And an upstream repo
    And I ran "git fetch upstream"
    When I run "git-town branch"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints:
      """
      * main
      """
