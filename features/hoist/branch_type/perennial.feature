Feature: hoisting an empty branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE      | PARENT | LOCATIONS     |
      | staging | perennial | main   | local, origin |
    And the current branch is "staging"
    When I run "git-town hoist"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | staging | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot hoist perennial branches
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "main"
    And the initial commits exist now
    And the initial lineage exists now
