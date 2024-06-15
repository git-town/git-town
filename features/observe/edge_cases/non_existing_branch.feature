Feature: cannot observe non-existing branches

  Background:
    When I run "git-town observe non-existing"

  @this
  Scenario: result
    Then it runs no commands
    And it prints the error:
      """
      there is no branch "non-existing"
      """
    And the current branch is still "feature"
    And the uncommitted file still exists
    And there are still no observed branches

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And there are still no observed branches
    And the current branch is still "feature"
