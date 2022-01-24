Feature: git town-ship: errors when trying to ship a perennial branch


  Background:
    Given my repo has the perennial branches "qa" and "production"
    And I am on the "production" branch
    When I run "git-town ship"


  Scenario: result
    Then it runs the commands
      | BRANCH     | COMMAND                  |
      | production | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "production" is not a feature branch. Only feature branches can be shipped
      """
    And I am still on the "production" branch
    And my repo now has the following commits
      | BRANCH | LOCATION |
