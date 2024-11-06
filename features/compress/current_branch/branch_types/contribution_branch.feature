Feature: does not compress contribution branches

  Background:
    Given a Git repo with origin
    And the branches
      | NAME         | TYPE         | PARENT | LOCATIONS     |
      | contribution | contribution |        | local, origin |
    And the current branch is "contribution"
    And the commits
      | BRANCH       | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | contribution | local, origin | commit 1 | file_1    | content 1    |
      |              |               | commit 2 | file_2    | content 2    |
    When I run "git-town compress"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH       | COMMAND                  |
      | contribution | git fetch --prune --tags |
    And it prints the error:
      """
      you are merely contributing to branch "contribution" and should leave compressing it to the branch owner
      """
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs no commands
    And the current branch is still "contribution"
    And the initial commits exist now
    And the initial branches and lineage exist now
