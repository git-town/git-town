Feature: does not ship perennial branches

  Background:
    Given a perennial branch "production"
    And the current branch is "main"
    And an uncommitted file
    When I run "git-town ship production"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | main   | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can be shipped
      """
    And the current branch is still "main"
    And the uncommitted file still exists
    And no branch hierarchy exists now

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And the current branch is still "main"
    And no branch hierarchy exists now
