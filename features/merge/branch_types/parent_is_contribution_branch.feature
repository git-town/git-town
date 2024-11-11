Feature: does not merge branches whose parent is contribution

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE         | PARENT | LOCATIONS |
      | parent  | contribution |        | local     |
      | current | feature      | parent | local     |
    And the current branch is "current"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | current | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot merge branch "current" because its parent branch (parent) has no parent
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "current"
    And the initial commits exist now
    And the initial lineage exists now
