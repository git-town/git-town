Feature: does not ship perennial branches

  Background:
    Given my repo has the perennial branch "production"
    And I am on the "main" branch
    And my workspace has an uncommitted file
    When I run "git-town ship production"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can be shipped
      """
    And I am still on the "main" branch
    And my workspace still contains my uncommitted file
    And Git Town now has no branch hierarchy information

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "main" branch
