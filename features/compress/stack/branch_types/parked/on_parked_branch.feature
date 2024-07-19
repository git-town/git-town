Feature: compresses active parked branches

  Background:
    Given a Git repo clone
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
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH | COMMAND                                         |
      | parked | git add -A                                      |
      |        | git stash                                       |
      |        | git checkout child                              |
      | child  | git reset --hard {{ sha 'child 2' }}            |
      |        | git push --force-with-lease --force-if-includes |
      |        | git checkout parked                             |
      | parked | git reset --hard {{ sha 'parked 2' }}           |
      |        | git push --force-with-lease --force-if-includes |
      |        | git stash pop                                   |
    And the current branch is still "parked"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
