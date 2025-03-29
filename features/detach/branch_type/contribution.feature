Feature: detaching a contribution branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution | main   | local, origin |
    And the current branch is "contribution"
    When I run "git-town detach"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                  |
      | contribution | git fetch --prune --tags |
    And Git Town prints the error:
      """
      cannot detach contribution branches since you don't own them
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial lineage exists now
