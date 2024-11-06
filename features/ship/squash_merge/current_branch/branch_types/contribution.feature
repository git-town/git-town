Feature: cannot ship contribution branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    And the current branch is "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE             |
      | contribution | local, origin | contribution commit |
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                  |
      | contribution | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot ship contribution branches
      """
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now
