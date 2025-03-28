Feature: detaching an observed branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME     | TYPE     | PARENT | LOCATIONS     |
      | observed | observed | main   | local, origin |
    And the current branch is "observed"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot detach observed branches
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "observed"
    And the initial commits exist now
    And the initial lineage exists now
