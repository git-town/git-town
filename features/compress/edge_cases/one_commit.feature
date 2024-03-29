Feature: does not compress branches containing only one commit

  Scenario: branch has 1 commit
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
    When I run "git-town compress"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      this branch has already just one commit
      """
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: branch has no commits
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH | LOCATION | MESSAGE | FILE NAME | FILE CONTENT |
    When I run "git-town compress"
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And it prints the error:
      """
      this branch has no commits
      """
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
