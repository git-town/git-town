Feature: does not compress perennial branches

  Background:
    Given the current branch is a perennial branch "perennial"
    And the commits
      | BRANCH    | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | perennial | local, origin | commit 1 | file_1    | content 1    |
      |           |               | commit 2 | file_2    | content 2    |
    When I run "git-town compress"

  Scenario: result
    Then it prints the error:
      """
      better not compress perennial branches
      """
    And it runs the commands
      | BRANCH    | COMMAND                  |
      | perennial | git fetch --prune --tags |
    And all branches are still synchronized
    And the current branch is still "perennial"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "perennial"
    And the initial commits exist
    And the initial branches and lineage exist
