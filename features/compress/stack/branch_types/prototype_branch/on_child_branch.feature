Feature: does not compress non-active prototype branches in the stack

  Background:
    Given a Git repo clone
    And the branch
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | prototype | local, origin | prototype 1 | parked_1  | prototype 1  |
      |           |               | prototype 2 | parked_2  | prototype 2  |
    And the branch
      | NAME  | TYPE    | PARENT    | LOCATIONS     |
      | child | feature | prototype | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | child  | local, origin | child 1 | child_1   | child 1      |
      |        |               | child 2 | child_2   | child 2      |
    And the current branch is "child"
    And an uncommitted file
    When I run "git-town compress --stack"

  Scenario: result
    Then it runs the commands
      | BRANCH    | COMMAND                                         |
      | child     | git fetch --prune --tags                        |
      |           | git add -A                                      |
      |           | git stash                                       |
      |           | git checkout prototype                          |
      | prototype | git reset --soft main                           |
      |           | git commit -m "prototype 1"                     |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout child                              |
      | child     | git reset --soft prototype                      |
      |           | git commit -m "child 1"                         |
      |           | git push --force-with-lease --force-if-includes |
      |           | git stash pop                                   |
    And all branches are now synchronized
    And the current branch is still "child"
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE     |
      | child     | local, origin | prototype 1 |
      |           |               | child 1     |
      | prototype | local, origin | prototype 1 |
    And file "parked_1" still has content "prototype 1"
    And file "parked_2" still has content "prototype 2"
    And the uncommitted file still exists

  Scenario: undo
    When I run "git-town undo"
    Then it runs the commands
      | BRANCH    | COMMAND                                         |
      | child     | git add -A                                      |
      |           | git stash                                       |
      |           | git reset --hard {{ sha 'child 2' }}            |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout prototype                          |
      | prototype | git reset --hard {{ sha 'prototype 2' }}        |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout child                              |
      | child     | git stash pop                                   |
    And the current branch is still "child"
    And the initial commits exist
    And the initial branches and lineage exist
    And the uncommitted file still exists
