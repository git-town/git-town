Feature: cannot observe the main branch

  Background:
    Given a Git repo with origin
    When I run "git-town observe"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot observe the main branch
      """
    And the main branch is still "main"
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the main branch is still "main"
    And there are still no observed branches
