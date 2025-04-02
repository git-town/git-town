Feature: swapping a branch with its observed parent

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE     | PARENT | LOCATIONS     |
      | parent  | observed | main   | local, origin |
      | current | feature  | parent | local, origin |
    And the current branch is "current"
    When I run "git-town swap"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap: branch "parent" is a observed branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "current"
    And the initial commits exist now
    And the initial lineage exists now
