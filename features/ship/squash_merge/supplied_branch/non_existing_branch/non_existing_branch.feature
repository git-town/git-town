Feature: does not ship a non-existing branch

  Background:
    Given a Git repo with origin
    And Git setting "git-town.ship-strategy" is "squash-merge"
    And the current branch is "main"
    When I run "git-town ship non-existing-branch"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      there is no branch "non-existing-branch"
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
