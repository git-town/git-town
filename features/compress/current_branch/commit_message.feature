Feature: compress the commits on a feature branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME    | TYPE    | PARENT | LOCATIONS     |
      | feature | feature | main   | local, origin |
    And the commits
      | BRANCH  | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | feature | local, origin | commit 1 | file_1    | content 1    |
      |         |               | commit 2 | file_2    | content 2    |
      |         |               | commit 3 | file_3    | content 3    |
    And the current branch is "feature"
    And an uncommitted file
    When I run "git-town compress -m compressed"

  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git fetch --prune --tags                        |
      |         | git add -A                                      |
      |         | git stash                                       |
      |         | git reset --soft main                           |
      |         | git commit -m compressed                        |
      |         | git push --force-with-lease --force-if-includes |
      |         | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "feature"
    And these commits exist now
      | BRANCH  | LOCATION      | MESSAGE    |
      | feature | local, origin | compressed |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH  | COMMAND                                         |
      | feature | git add -A                                      |
      |         | git stash                                       |
      |         | git reset --hard {{ sha 'commit 3' }}           |
      |         | git push --force-with-lease --force-if-includes |
      |         | git stash pop                                   |
    And the current branch is still "feature"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
