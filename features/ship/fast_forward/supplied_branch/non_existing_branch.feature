Feature: does not ship a non-existing branch

  Background:
    Given a Git repo with origin
    And the current branch is "main"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship non-existing-branch"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is no branch "non-existing-branch"
      """
    And the current branch is now "main"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is still "main"
