Feature: detaching an empty branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE      | PARENT | LOCATIONS     |
      | staging | perennial | main   | local, origin |
    And the current branch is "staging"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | staging | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot detach perennial branches since you don't own them
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial commits exist now
    And the initial lineage exists now
