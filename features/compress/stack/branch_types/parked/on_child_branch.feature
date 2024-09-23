Feature: does not compress non-active parked branches in the stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME   | TYPE   | PARENT | LOCATIONS     |
      | parked | parked | main   | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE  | FILE NAME | FILE CONTENT |
      | parked | local, origin | parked 1 | parked_1  | parked 1     |
      |        |               | parked 2 | parked_2  | parked 2     |
    And the branches
      | NAME  | TYPE    | PARENT | LOCATIONS     |
      | child | feature | parked | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | child  | local, origin | child 1 | child_1   | child 1      |
      |        |               | child 2 | child_2   | child 2      |
    And the current branch is "child"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --soft parked                         |
      |        | git commit -m "child 1"                         |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "child"
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | child  | local, origin | parked 1 |
      |        |               | parked 2 |
      |        |               | child 1  |
      | parked | local, origin | parked 1 |
      |        |               | parked 2 |
    And file "parked_1" still has content "parked 1"
    And file "parked_2" still has content "parked 2"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git add -A                                      |
      |        | git stash                                       |
      |        | git reset --hard {{ sha 'child 2' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And the current branch is still "child"
    And the initial commits exist now
    And the initial branches and lineage exist now
    And the uncommitted file still exists
