Feature: cannot ship contribution branches

  Background:
    Given the current branch is a contribution branch "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | local, origin | contribution commit |
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH       | COMMAND                  |
      | contribution | git fetch --prune --tags |
    And it prints the error:
      """
      cannot ship contribution branches
      """
    And the current branch is still "contribution"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "contribution"
    And the initial commits exist
    And the initial branches and lineage exist
