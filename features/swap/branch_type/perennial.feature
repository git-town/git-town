Feature: swapping a perennial branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE      | LOCATIONS     |
      | current | perennial | local, origin |
    And the current branch is "current"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap a branch without parent
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "current"
    And the initial commits exist now
    And the initial lineage exists now
