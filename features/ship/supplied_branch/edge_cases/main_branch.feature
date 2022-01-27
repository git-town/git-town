Feature: does not ship the main branch

  Background:
    Given my repo has a feature branch named "feature"
    And I am on the "feature" branch
    And my workspace has an uncommitted file
    When I run "git-town ship main"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can be shipped
      """
    And I am still on the "feature" branch
    And my workspace still contains my uncommitted file

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "feature" branch
