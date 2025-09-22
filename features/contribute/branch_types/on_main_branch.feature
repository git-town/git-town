Feature: cannot make the main branch a contribution branch

  Background:
    Given a Git repo with origin
    When I run "git-town contribute"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot make the main branch a contribution branch
      """
    And the main branch is still "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the main branch is still "main"
