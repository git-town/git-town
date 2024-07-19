Feature: already existing remote branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | origin    |
    And an uncommitted file
    And I run "git fetch"
    When I run "git-town hack existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is already a branch "existing" at the "origin" remote
      """
    And the current branch is still "main"
    And no commits exist now
    And the initial branches and lineage exist
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
