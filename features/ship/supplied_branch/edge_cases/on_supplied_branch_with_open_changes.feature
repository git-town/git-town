Feature: does not ship a branch that has open changes

  Background:
    Given the current branch is a feature branch "feature"
    And an uncommitted file
    When I run "git-town ship feature"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """
    And the current branch is still "feature"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "feature"
