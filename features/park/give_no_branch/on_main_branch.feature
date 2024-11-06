Feature: park the current main branch

  Background:
    Given a Git repo with origin
    When I run "git-town park"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot park the main branch
      """
    And the current branch is still "main"
    And the main branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "main"
    And the main branch is still "main"
    And there are now no parked branches
