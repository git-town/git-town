Feature: hack an existing contribution branch

  Background:
    Given a Git repo with origin
    And the current branch is "main"
    When I run "git-town hack main"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      you are trying to convert the main branch to a feature branch
      """
    And branch "main" still has type "main"
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial branches and lineage exist now
    And branch "main" now has type "main"
