Feature: does not ship empty feature branches

  Background:
    Given the feature branches "empty" and "other"
    And the commits
      | BRANCH | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main   | local    | main commit    | common_file | common content |
      | empty  | local    | feature commit | common_file | common content |
    And the current branch is "other"
    And an uncommitted file
    When I run "git-town ship empty"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git stash pop            |
    And it prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And the current branch is still "other"
    And the uncommitted file still exists
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "other"
    And the initial commits exist
    And the initial branches and lineage exist
