Feature: does not compress non-active prototype branches in the stack

  Background:
    Given a Git repo with origin
    And the branches
      | NAME      | TYPE      | PARENT | LOCATIONS     |
      | prototype | prototype | main   | local, origin |
    And the commits
      | BRANCH    | LOCATION      | MESSAGE     | FILE NAME | FILE CONTENT |
      | prototype | local, origin | prototype 1 | parked_1  | prototype 1  |
      |           |               | prototype 2 | parked_2  | prototype 2  |
    And the branches
      | NAME  | TYPE    | PARENT    | LOCATIONS     |
      | child | feature | prototype | local, origin |
    And the commits
      | BRANCH | LOCATION      | MESSAGE | FILE NAME | FILE CONTENT |
      | child  | local, origin | child 1 | child_1   | child 1      |
      |        |               | child 2 | child_2   | child 2      |
    And the current branch is "child"
    When I run "git-town compress --stack"

  Scenario: result
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | child     | git fetch --prune --tags                        |
      |           | git checkout prototype                          |
      | prototype | git reset --soft main --                        |
      |           | git commit -m "prototype 1"                     |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout child                              |
      | child     | git reset --soft prototype --                   |
      |           | git commit -m "child 1"                         |
      |           | git push --force-with-lease --force-if-includes |
    And all branches are now synchronized
    And these commits exist now
      | BRANCH    | LOCATION      | MESSAGE     |
      | prototype | local, origin | prototype 1 |
      | child     | local, origin | child 1     |
    And file "parked_1" still has content "prototype 1"
    And file "parked_2" still has content "prototype 2"

  Scenario: undo
    When I run "git-town undo"
    Then Git Town runs the commands
      | BRANCH    | COMMAND                                         |
      | child     | git reset --hard {{ sha 'child 2' }}            |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout prototype                          |
      | prototype | git reset --hard {{ sha 'prototype 2' }}        |
      |           | git push --force-with-lease --force-if-includes |
      |           | git checkout child                              |
    And the initial branches and lineage exist now
    And the initial commits exist now
