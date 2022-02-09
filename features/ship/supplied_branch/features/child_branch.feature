Feature: does not ship a child branch

  Background:
    Given a feature branch "alpha"
    And a feature branch "beta" as a child of "alpha"
    And a feature branch "gamma" as a child of "beta"
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "alpha"
    When I run "git-town ship gamma -m 'gamma done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | alpha  | git fetch --prune --tags |
    And it prints the error:
      """
      shipping this branch would ship "alpha, beta" as well,
      please ship "alpha" first
      """
    And the current branch is now "alpha"
    And now the initial commits exist
    And Git Town is now aware of the initial branch hierarchy

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints the error:
      """
      nothing to undo
      """
    And the current branch is still "alpha"
    And now the initial commits exist
    And Git Town is now aware of the initial branch hierarchy
