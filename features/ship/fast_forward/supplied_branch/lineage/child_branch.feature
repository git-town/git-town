Feature: does not ship a child branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
      | gamma | feature | beta   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      |
      | alpha  | local, origin | alpha commit |
      | beta   | local, origin | beta commit  |
      | gamma  | local, origin | gamma commit |
    And the current branch is "alpha"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship gamma"

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
