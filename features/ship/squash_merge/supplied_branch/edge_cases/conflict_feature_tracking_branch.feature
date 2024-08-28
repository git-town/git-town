Feature: refuses shipping a branch with conflicts between the supplied feature branch and its tracking branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "other"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "squash-merge"
    And I run "git-town ship feature -m 'feature done'"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" is not in sync
      """
    And the current branch is still "other"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "other"
    And the uncommitted file still exists
    And no merge is in progress
    And the initial commits exist
    And the initial lineage exists
