Feature: does not compress observed branches

  Background:
    Given the current branch is an observed branch "observed"
    And the commits
      | BRANCH   | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | observed | local, origin | commit 1 | file_1    | content 1    |
      |          |               | commit 2 | file_2    | content 2    |
    When I run "git-town compress"

  Scenario: result
    Then it runs the commands
      | BRANCH   | COMMAND                  |
      | observed | git fetch --prune --tags |
    And it prints the error:
      """
      you are merely observing branch "observed" and should leave compressing it to the branch owner
      """
    And the current branch is still "observed"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "observed"
    And the initial commits exist
    And the initial branches and lineage exist
