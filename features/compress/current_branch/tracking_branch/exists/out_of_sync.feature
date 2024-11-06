Feature: cannot compress branches that are out of sync

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"

  Scenario: local branch is behind
    Given the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | origin        | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then Git Town runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
    And Git Town prints the error:
      """
      please sync branch "feature" before compressing it
      """
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: local branch is ahead
    Given the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | local         | commit 2 | file_2    | content 2    |
    When I run "git-town compress"
    Then Git Town prints the error:
      """
      please sync branch "feature" before compressing it
      """
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: branch is ahead and behind
    Given the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         | local         | commit 2 | file_2    | content 2    |
      |         | origin        | commit 3 | file_3    | content 3    |
    When I run "git-town compress"
    Then Git Town prints the error:
      """
      please sync branch "feature" before compressing it
      """
    And the initial commits exist now
    And the initial branches and lineage exist now

  Scenario: branch is deleted at remote
    Given the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
    And origin deletes the "feature" branch
    When I run "git-town compress"
    Then Git Town prints the error:
      """
      please sync branch "feature" before compressing it
      """
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local    | commit 1 | file_1    | content 1    |
      |         |          | commit 2 | file_2    | content 2    |
    And the initial branches and lineage exist now
