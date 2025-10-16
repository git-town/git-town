Feature: detaching the main branch

  Background:
    Given a Git repo with origin
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot detach branches without parent
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the initial lineage exists now
    And the initial commits exist now
