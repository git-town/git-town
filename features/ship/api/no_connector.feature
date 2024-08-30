Feature: cannot ship a branch without connector

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE        |
      | feature | local, origin | feature commit |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "api"
    And the origin is "git@unknown.com:something/whatever.git"
    And a proposal for this branch does not exist
    When I run "git-town ship -m done"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      shipping via the API requires a connector
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the initial branches and lineage exist
    And the initial commits exist
