Feature: already existing local branch

  Background:
    Given a local feature branch "existing"
    And an uncommitted file
    When I run "git-town append existing"

  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is already a branch "existing"
      """
    And the uncommitted file still exists

  @debug @this
  Scenario: undo
    When I run "git-town undo -v"
    Then it runs no commands
    And the current branch is now "main"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
