Feature: does not ship perennial branches

  Scenario: try to ship the main branch
    Given the current branch is "main"
    When I run "git-town ship -m done"
    Then it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can be shipped
      """
    And the current branch is still "main"

  Scenario: try to ship a perennial branch
    Given the perennial branches "qa" and "production"
    And the current branch is "production"
    When I run "git-town ship"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can be shipped
      """
    And the current branch is still "production"
