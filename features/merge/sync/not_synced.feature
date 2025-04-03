Feature: merging branches using the "rebase" sync-strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
      | beta  | feature | alpha  | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
      | beta   | local, origin | beta commit  | beta-file  | beta content  |
    And the current branch is "beta"
    And Git setting "git-town.sync-feature-strategy" is "rebase"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branches "beta" and "alpha" are not in sync, please run "git town sync" and try again
      """


  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
