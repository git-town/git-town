@this
Feature: cannot compress branches that are out of sync

  Scenario: local branch is behind
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | origin        | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then it prints the error:
      """
      Please sync branch "feature" before compressing it
      """
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: local branch is ahead
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | local         | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then it prints the error:
      """
      Please sync branch "feature" before compressing it
      """
    And the initial commits exist
    And the initial branches and lineage exist

  Scenario: branch is ahead and behind
    Given the current branch is a feature branch "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | local         | commit 2 | file_2    | content 2    |
      |         | origin        | commit 3 | file_3    | content 3    |
    When I run "git-town compress"
    Then it prints the error:
      """
      Please sync branch "feature" before compressing it
      """
    And the initial commits exist
    And the initial branches and lineage exist
