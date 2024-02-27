Feature: cannot park non-existing branches

  Background:
    Given the current branch is a feature branch "feature"
    And an uncommitted file
    When I run "git-town park feature non-existing"

  Scenario: result
    Then it runs no commands
    And the current branch is still "feature"
    And the uncommitted file still exists
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And there are still no parked branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are still no parked branches
    And the current branch is still "feature"
