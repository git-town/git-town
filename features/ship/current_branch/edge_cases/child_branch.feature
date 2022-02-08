Feature: does not ship a child branch

  Background:
    Given my repo has a feature branch "alpha"
    And my repo has a feature branch "beta" as a child of "alpha"
    And my repo has a feature branch "gamma" as a child of "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, remote | alpha commit |
      | beta   | local, remote | beta commit  |
      | gamma  | local, remote | gamma commit |
    And I am on the "gamma" branch
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | gamma  | git fetch --prune --tags |
    And it prints the error:
      """
      shipping this branch would ship "alpha, beta" as well,
      please ship "alpha" first
      """
    And I am still on the "gamma" branch
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And I am still on the "gamma" branch
    And my repo is left with my initial commits
    And my repo now has its initial branches and branch hierarchy
