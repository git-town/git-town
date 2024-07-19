Feature: already existing local branch

  Background:
    Given a Git repo clone
    And the branch
      | NAME     | TYPE    | PARENT | LOCATIONS |
      | existing | feature | main   | local     |
    And an uncommitted file
    When I run "git-town append existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is already a branch "existing"
      """
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
