Feature: make the main branch a feature brancha

  Background:
    Given a Git repo with origin
    And the current branch is "main"
    When I run "git-town feature main"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      cannot make the main branch a feature branch
      """
    And the initial branches and lineage exist now
    And branch "main" still has type "main"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial branches and lineage exist now
    And branch "main" now has type "main"
    And the initial commits exist now
