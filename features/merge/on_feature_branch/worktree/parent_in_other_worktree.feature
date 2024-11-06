Feature: merging a branch in a stack with its parent

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
    And branch "alpha" is active in another worktree
    And the current branch is "beta"
    When I run "git-town merge"

  Scenario: result
    Then it prints the error:
      """
      branch "alpha"
      """
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | beta   | git fetch --prune --tags |
    And the current branch is still "beta"
    And the initial lineage exists now
    And the initial commits exist now

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "beta"
    And the initial commits exist now
    And the initial lineage exists now
