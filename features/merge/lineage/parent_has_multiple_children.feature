Feature: dry-run merging branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | alpha | feature | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE      | FILE NAME  | FILE CONTENT  |
      | alpha  | local, origin | alpha commit | alpha-file | alpha content |
    And the branches
      | NAME   | TYPE    | PARENT | LOCATIONS     |
      | beta_1 | feature | alpha  | local, origin |
      | beta_2 | feature | alpha  | local, origin |
    And the current branch is "beta_1"
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta_1 | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "alpha" has more than one child
      """

  Scenario: undo
    When I run "git-town undo"
    And Git Town runs no commands
    And the initial commits exist now
    And the initial branches exist now
