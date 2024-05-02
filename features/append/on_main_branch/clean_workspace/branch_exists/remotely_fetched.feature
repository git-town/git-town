Feature: already existing known remote branch

  Background:
    Given a remote branch "existing"
    And I run "git fetch"
    When I run "git-town append existing"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is already a branch "existing" at the "origin" remote
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
