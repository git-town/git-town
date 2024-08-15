Feature: dry-run compressing the commits on a feature branch

  Background:
    Given a Git repo with origin
    And the branch
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the current branch is "feature"
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
      |         |               | commit 3 | file_3    | content 3    |
    And an uncommitted file
    When I run "git-town compress --dry-run"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git add -A                                      |
      |         | git stash                                       |
      |         | git reset --soft main                           |
      |         | git commit -m "commit 1"                        |
      |         | git push --force-with-lease --force-if-includes |
      |         | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "feature"
    And the initial commits exist
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs no commands
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
