Feature: cannot ship perennial branches

  Background:
    Given the current branch is a perennial branch "perennial"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE          |
      | perennial | local, origin | perennial commit |
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
