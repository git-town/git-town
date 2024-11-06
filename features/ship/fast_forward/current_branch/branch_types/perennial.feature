Feature: cannot ship perennial branches using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | perennial | local, origin | perennial commit |
    And the current branch is "perennial"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot ship perennial branches
      """
    And the current branch is still "perennial"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "perennial"
    And the initial commits exist now
    And the initial branches and lineage exist now
