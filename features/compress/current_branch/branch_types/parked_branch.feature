Feature: compress the commits on a parked branch

  Background:
    Given a Git repo clone
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked |        | local, origin |
    Given the current branch is "parked"
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | parked | local, origin | commit 1 | file_1    | content 1    |
      |        |               | commit 2 | file_2    | content 2    |
      |        |               | commit 3 | file_3    | content 3    |
    And an uncommitted file
    When I run "git-town compress"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --soft main                           |
      |        | git commit -m "commit 1"                        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "parked"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | parked | local, origin | commit 1 |
    And file "file_1" still has content "content 1"
    And file "file_2" still has content "content 2"
    And file "file_3" still has content "content 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --hard {{ sha 'commit 3' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And the current branch is still "parked"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
