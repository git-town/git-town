Feature: does not ship an empty branch

  Background:
    Given the current branch is a feature branch "empty"
    And the commits
      | BRANCH | LOCATION | MESSAGE      | FILE NAME   | FILE CONTENT   |
      | main   | local    | main commit  | common_file | common content |
      | empty  | local    | empty commit | common_file | common content |
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
