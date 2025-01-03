Feature: does not ship uncommitted changes using the always-merge strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And Git setting "git-town.ship-strategy" is "always-merge"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      you have uncommitted changes. Did you mean to commit them before shipping?
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
