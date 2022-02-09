Feature: does not ship a non-existing branch

  Background:
    Given the current branch is "main"
    And my workspace has an uncommitted file
    When I run "git-town ship non-existing-branch"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      there is no branch named "non-existing-branch"
      """
    And the current branch is now "main"
    And my workspace still contains my uncommitted file

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And the current branch is still "main"
