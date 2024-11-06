Feature: does not ship a branch that has open changes

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """
    And the current branch is still "feature"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is still "feature"
