Feature: does not ship a child branch

  Background:
    Given the current branch is a feature branch "alpha"
    And a feature branch "beta" as a child of "alpha"
    And a feature branch "gamma" as a child of "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    When I run "git-town ship gamma -m 'gamma done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And it prints the error:
      """
      shipping this branch would ship "alpha" and "beta" as well,
      please ship "alpha" first
      """
    And the current branch is now "alpha"
    And the initial commits exist
    And the initial lineage exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "alpha"
    And the initial commits exist
    And the initial lineage exists
