Feature: does not ship out-of-sync branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION | MESSAGE                   | FILE NAME        | FILE CONTENT   |
      | feature | local    | conflicting local commit  | conflicting_file | local content  |
      |         | origin   | conflicting origin commit | conflicting_file | origin content |
    And the current branch is "feature"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      branch "feature" is not in sync
      """

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And no merge is in progress
    And the initial commits exist
    And the initial lineage exists
