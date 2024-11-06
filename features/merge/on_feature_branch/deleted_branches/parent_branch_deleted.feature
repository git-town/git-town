Feature: merging a branch whose parent was deleted remotely

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
    And origin deletes the "alpha" branch
    When I run "git-town merge"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "alpha" was deleted at the remote
      """
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
