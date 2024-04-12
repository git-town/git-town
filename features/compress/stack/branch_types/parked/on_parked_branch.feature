Feature: does not compress non-active parked branches

  Background:
    Given parked branch "parked" with these commits
      | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | local, origin | parked 1 | parked_1  | parked 1     |
      |               | parked 2 | parked_2  | parked 2     |
      |               | parked 3 | parked_3  | parked 3     |
    And feature branch "child" as a child of "parked" has these commits
      | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | local, origin | child 1 | child_1   | child 1      |
      |               | child 2 | child_2   | child 2      |
      |               | child 3 | child_3   | child 3      |
    And the current branch is "parked"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --soft main                           |
      |        | git commit -m "parked 1"                        |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout child                              |
      | child  | git reset --soft parked                         |
      |        | git commit -m "child 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout parked                             |
      | parked | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "parked"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | child  | local, origin | parked 1 |
      |        |               | child 1  |
      | parked | local, origin | parked 1 |
    And file "parked_1" still has content "parked 1"
    And file "parked_2" still has content "parked 2"
    And file "parked_3" still has content "parked 3"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git add -A                                      |
      |        | git stash                                       |
      |        | git checkout child                              |
      | child  | git reset --hard {{ sha 'child 3' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout parked                             |
      | parked | git reset --hard {{ sha 'parked 3' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And the current branch is still "parked"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
