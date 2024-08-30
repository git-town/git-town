Feature: does not ship the given out-of-sync branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
      | other   | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE       |
      | feature | local    | local commit  |
      |         | origin   | origin commit |
    And the current branch is "other"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "fast-forward"
    And I run "git-town ship feature"

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
