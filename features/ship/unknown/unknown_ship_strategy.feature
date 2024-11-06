Feature: unknown ship strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "zonk"
    And the origin is "git@github.com:git-town/git-town.git"
    And a proposal for this branch does not exist
    When I run "git-town ship -m done"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND |
    And it prints the error:
      """
      unknown ship strategy: "zonk"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And it prints the error:
      """
      unknown ship strategy: "zonk"
      """
    And the initial branches and lineage exist now
    And the initial commits exist now
