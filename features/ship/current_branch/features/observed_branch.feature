Feature: cannot ship observed branches

  Background:
    Given a Git repo clone
    And the branches
      | NAME     | TYPE     | LOCATIONS     |
      | observed | observed | local, origin |
    Given the current branch is "observed"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE         |
      | observed | local, origin | observed commit |
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
    And it prints the error:
      """
      cannot ship observed branches
      """
    And the current branch is still "observed"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "observed"
    And the initial commits exist
    And the initial branches and lineage exist
