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
    And Git setting "git-town.ship-strategy" is "fast-forward"
    And I run "git-town ship feature"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
    And Git Town prints the error:
      """
      branch "feature" is not in sync
      """

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And no merge is in progress
    And the initial commits exist now
    And the initial lineage exists now
