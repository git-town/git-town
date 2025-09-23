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
    And branch "main" still has type "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And branch "main" still has type "main"
