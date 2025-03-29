Feature: swapping an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT   | LOCATIONS     |
      | branch-1 | feature  | main     | local, origin |
      | observed | observed | branch-1 | local, origin |
    And the current branch is "observed"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap: branch "observed" is a observed branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "observed"
    And the initial commits exist now
    And the initial lineage exists now
