Feature: cannot ship perennial branches

  Background:
    Given a Git repo with origin
    And the branch
      | NAME      | TYPE      | LOCATIONS     |
      | perennial | perennial | local, origin |
    And the current branch is "perennial"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | perennial | local, origin | perennial commit |
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And it prints the error:
      """
      cannot ship perennial branches
      """
    And the current branch is still "perennial"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the initial commits exist
    And the initial branches and lineage exist
