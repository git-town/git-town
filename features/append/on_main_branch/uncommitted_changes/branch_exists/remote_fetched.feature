Feature: already existing known remote branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | origin    |
    And an uncommitted file
    And I run "git fetch"
    When I run "git-town append existing"

  Scenario: result
    Then Git Town runs no commands
    And it prints the error:
      """
      there is already a branch "existing"
      """
    And the current branch is still "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is now "main"
    And the initial commits exist now
    And the initial branches and lineage exist now
    And the uncommitted file still exists
