Feature: swapping a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE         | PARENT | LOCATIONS     |
      | parent  | feature      | main   | local, origin |
      | current | contribution | parent | local, origin |
    And the current branch is "current"
    When I run "git-town swap"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot swap: branch current is a contribution branch
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial commits exist now
