Feature: compress the commits on a local feature branch

  Background:
    Given a Git repo with origin
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS |
      | feature | feature | main   | local     |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local    | commit 1 | file_1    | content 1    |
      |         |          | commit 2 | file_2    | content 2    |
      |         |          | commit 3 | file_3    | content 3    |
    And an uncommitted file
    When I run "git-town compress"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                  |
      | feature | git fetch --prune --tags |
      |         | git add -A               |
      |         | git stash                |
      |         | git reset --soft main    |
      |         | git commit -m "commit 1" |
      |         | git stash pop            |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION | MESSAGE  |
      | feature | local    | commit 1 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                               |
      | feature | git add -A                            |
      |         | git stash                             |
      |         | git reset --hard {{ sha 'commit 3' }} |
      |         | git stash pop                         |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
