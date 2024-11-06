Feature: does not ship empty feature branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | empty | feature | main   | local     |
      | other | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE        | FILE NAME   | FILE CONTENT   |
      | main   | local    | main commit    | common_file | common content |
      | empty  | local    | feature commit | common_file | common content |
    And the current branch is "other"
    And an uncommitted file
    And Git Town setting "ship-strategy" is "squash-merge"
    When I run "git-town ship empty"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                  |
      | other  | git fetch --prune --tags |
      |        | git add -A               |
      |        | git stash                |
      |        | git stash pop            |
    And Git Town prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And the current branch is still "other"
    And the uncommitted file still exists
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And Git Town prints:
      """
      nothing to undo
      """
    And the current branch is still "other"
    And the initial commits exist now
    And the initial branches and lineage exist now
