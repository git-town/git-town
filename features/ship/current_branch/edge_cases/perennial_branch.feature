Feature: does not ship perennial branches

  Scenario: try to ship the main branch
    Given I am on the "main" branch
    When I run "git-town ship -m done"
    Then it prints the error:
      """
      the branch "main" is not a feature branch. Only feature branches can be shipped
      """
    And I am still on the "main" branch

  Scenario: try to ship a perennial branch
    Given my repo has the perennial branches "qa" and "production"
    And I am on the "production" branch
    When I run "git-town ship"
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can be shipped
      """
    And I am still on the "production" branch
