Feature: already existing remote branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE    | PARENT | LOCATIONS     |
      | old      | feature | main   | local, origin |
      | existing | feature | main   | origin        |
    Given the current branch is "old"
    And an uncommitted file
    And I run "git fetch"
    When I run "git-town prepend existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is already a branch "existing"
      """
    And the current branch is still "old"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is now "old"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial lineage exists
