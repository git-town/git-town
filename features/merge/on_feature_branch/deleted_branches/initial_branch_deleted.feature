Feature: merging a branch that was deleted at the remote

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
    And origin deletes the "beta" branch
    When I run "git-town merge"

  Scenario: result
    Then it prints the error:
      """
      branch "beta" was deleted at the remote
      """
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then it prints:
      """
      nothing to undo
      """
    Then it runs no commands
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
