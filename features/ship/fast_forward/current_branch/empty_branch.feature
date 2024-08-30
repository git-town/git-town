Feature: does not ship an empty branch using the fast-forward strategy

  Background:
    Given a Git repo with origin
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS |
      | empty | feature | main   | local     |
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME | FILE CONTENT |
      | main   | local    | main commit  | same_file | same content |
      | empty  | local    | empty commit | same_file | same content |
    And the current branch is "empty"
    And Git Town setting "ship-strategy" is "fast-forward"
    When I run "git-town ship"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                  |
      | empty  | git fetch --prune --tags |
    And it prints the error:
      """
      the branch "empty" has no shippable changes
      """
    And the current branch is still "empty"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And it prints:
      """
      nothing to undo
      """
    And the current branch is still "empty"
    And the initial commits exist
    And the initial branches and lineage exist
