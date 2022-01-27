Feature: cannot ship a non-existing branch

  Background:
    Given I am on the "main" branch
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
    And I am now on the "main" branch
    And my workspace still contains my uncommitted file

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "main" branch
