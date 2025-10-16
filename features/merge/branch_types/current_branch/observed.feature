Feature: does not merge observed branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE     | PARENT | LOCATIONS     |
      | parent  | feature  | main   | local         |
      | current | observed | parent | local, origin |
    And the current branch is "current"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot merge observed branches
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial commits exist now
