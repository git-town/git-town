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
    When I run "git-town compress --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git fetch --prune --tags                        |
      |        | git reset --soft parked --                      |
      |        | git commit -m "child 1"                         |
      |        | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH | LOCATION      | MESSAGE  |
      | parked | local, origin | parked 1 |
      |        |               | parked 2 |
      | child  | local, origin | child 1  |
    And file "parked_1" still has content "parked 1"
    And file "parked_2" still has content "parked 2"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH | COMMAND                                         |
      | child  | git reset --hard {{ sha 'child 2' }}            |
      |        | git push --force-with-lease --force-if-includes |
    And the initial branches and lineage exist now
    And the initial commits exist now
