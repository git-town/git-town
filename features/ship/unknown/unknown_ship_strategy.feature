Feature: unknown ship strategy

  Background:
    Given a Git repo with origin
    And the origin is "git@github.com:git-town/git-town.git"
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And Git setting "git-town.ship-strategy" is "zonk"
    And the current branch is "feature"
    When I run "git-town ship -m done"

  Scenario: result
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      unknown ship strategy in git-town.ship-strategy: "zonk"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints the error:
      """
      unknown ship strategy in git-town.ship-strategy: "zonk"
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
