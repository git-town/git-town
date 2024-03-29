Feature: cannot compress branches that are out of sync

  @this
  Scenario: local branch behind
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | origin        | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then it prints the error:
      """
      Please sync this branch before compressing it.
      """
    And the initial commits exist
    And the initial branches and lineage exist
