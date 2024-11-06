Feature: cannot ship observed branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    And the current branch is "observed"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | observed | local, origin | observed commit |
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot ship observed branches
      """
    And the current branch is still "observed"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "observed"
    And the initial commits exist now
    And the initial branches and lineage exist now
